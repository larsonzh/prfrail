package snapshot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

type CaptureOptions struct {
	SourceDir          string
	Kind               string // "baseline" or "candidate"
	SnapshotID         string
	RunID              string
	ParentSnapshotHash *string
	Task               *TaskBinding
	Exclusions         []Exclusion
	Store              *Store
	Environment        *Environment
	MaxBytesQuota      int64
}

type fileIdentity struct {
	path string
	info os.FileInfo
}

func Capture(ctx context.Context, opts CaptureOptions) (SnapshotManifest, error) {
	if opts.Store == nil {
		return SnapshotManifest{}, errors.New("snapshot store is required")
	}
	srcInfo, err := os.Stat(opts.SourceDir)
	if err != nil {
		return SnapshotManifest{}, fmt.Errorf("source dir stat: %w", err)
	}
	if !srcInfo.IsDir() {
		return SnapshotManifest{}, fmt.Errorf("%w: source path is not a directory", ErrInvalidPath)
	}

	// Prepare exclusions
	exclusions := make([]Exclusion, 0, len(opts.Exclusions)+3)
	// Default exclusions
	defaultExclusions := []Exclusion{
		{Pattern: ".git/**", Source: "default", Reason: "version-control-metadata"},
		{Pattern: ".prfrail/**", Source: "default", Reason: "runtime-store"},
		{Pattern: "tmp/**", Source: "default", Reason: "temporary-scratch"},
	}
	exclusions = append(exclusions, defaultExclusions...)
	exclusions = append(exclusions, opts.Exclusions...)

	matchers := make([]func(string) bool, len(exclusions))
	for i, ex := range exclusions {
		m, err := compileGlobMatcher(ex.Pattern)
		if err != nil {
			return SnapshotManifest{}, fmt.Errorf("compile exclusion pattern %q: %w", ex.Pattern, err)
		}
		matchers[i] = m
	}

	isExcluded := func(relPath string) bool {
		for _, m := range matchers {
			if m(relPath) {
				return true
			}
		}
		return false
	}

	caseMap := make(map[string]string)
	entriesMap := make(map[string]Entry)
	var capturedBytes int64
	var seenFiles []fileIdentity
	hardlinkCounter := 0
	hardlinkGroups := make(map[string]string) // path -> groupID

	err = filepath.Walk(opts.SourceDir, func(fullPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(opts.SourceDir, fullPath)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		relPath := filepath.ToSlash(rel)

		if isExcluded(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if err := ValidatePath(relPath); err != nil {
			return err
		}

		lower := strings.ToLower(relPath)
		if prev, exists := caseMap[lower]; exists && prev != relPath {
			return fmt.Errorf("%w: %q and %q", ErrCaseCollision, relPath, prev)
		}
		caseMap[lower] = relPath

		// Check for unportable reparse points / junctions
		if isReparsePointOrJunction(info, fullPath) {
			return fmt.Errorf("%w: path %q is an unportable reparse point or junction", ErrUnsupportedLink, relPath)
		}

		mode := info.Mode()

		if mode&os.ModeSymlink != 0 {
			target, err := os.Readlink(fullPath)
			if err != nil {
				return fmt.Errorf("readlink %q: %w", relPath, err)
			}
			normTarget := filepath.ToSlash(target)
			if err := validateLinkTargetEscape(relPath, normTarget); err != nil {
				return err
			}
			targetHash := evidence.Digest("", []byte(normTarget))
			entriesMap[relPath] = Entry{
				Path:       relPath,
				Type:       EntrySymbolicLink,
				Target:     normTarget,
				TargetHash: targetHash,
			}
			return nil
		}

		if mode.IsDir() {
			entriesMap[relPath] = Entry{
				Path:     relPath,
				Type:     EntryDirectory,
				ReadOnly: mode.Perm()&0222 == 0,
			}
			return nil
		}

		if mode.IsRegular() {
			// Concurrent modification check: stat before, read, stat after
			statBefore, err := os.Lstat(fullPath)
			if err != nil {
				return fmt.Errorf("lstat before read %q: %w", relPath, err)
			}

			data, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("read file %q: %w", relPath, err)
			}

			statAfter, err := os.Lstat(fullPath)
			if err != nil {
				return fmt.Errorf("lstat after read %q: %w", relPath, err)
			}

			if statBefore.Size() != statAfter.Size() || !statBefore.ModTime().Equal(statAfter.ModTime()) || int64(len(data)) != statBefore.Size() {
				return fmt.Errorf("%w: file %q changed during capture", ErrConcurrentChange, relPath)
			}

			// Quota check
			if opts.MaxBytesQuota > 0 && capturedBytes+int64(len(data)) > opts.MaxBytesQuota {
				return fmt.Errorf("%w: captured %d + new %d > quota %d", ErrQuotaExceeded, capturedBytes, len(data), opts.MaxBytesQuota)
			}
			capturedBytes += int64(len(data))

			hash, err := opts.Store.PutObject(ctx, data)
			if err != nil {
				return fmt.Errorf("store object %q: %w", relPath, err)
			}

			// Hardlink detection
			var groupID *string
			for _, prev := range seenFiles {
				if os.SameFile(info, prev.info) {
					existingGroup, ok := hardlinkGroups[prev.path]
					if !ok {
						hardlinkCounter++
						existingGroup = fmt.Sprintf("hlg-%03d", hardlinkCounter)
						hardlinkGroups[prev.path] = existingGroup
						// Update previous entry
						if prevEntry, exists := entriesMap[prev.path]; exists {
							prevEntry.HardlinkGroup = &existingGroup
							entriesMap[prev.path] = prevEntry
						}
					}
					hardlinkGroups[relPath] = existingGroup
					groupID = &existingGroup
					break
				}
			}
			seenFiles = append(seenFiles, fileIdentity{path: relPath, info: info})

			// Role detection
			role := RoleRegular
			base := filepath.Base(relPath)
			switch base {
			case "package.json", "go.mod", "pom.xml", "Cargo.toml", "requirements.txt", "pyproject.toml":
				role = RolePackageManifest
			case "package-lock.json", "go.sum", "Cargo.lock", "yarn.lock", "pnpm-lock.yaml":
				role = RoleLockfile
			}

			entriesMap[relPath] = Entry{
				Path:          relPath,
				Type:          EntryRegularFile,
				ContentHash:   hash,
				SizeBytes:     int64(len(data)),
				Executable:    mode.Perm()&0111 != 0,
				ReadOnly:      mode.Perm()&0222 == 0,
				Role:          role,
				HardlinkGroup: groupID,
			}
			return nil
		}

		return fmt.Errorf("%w: unsupported file mode %s for %q", ErrInvalidPath, mode.String(), relPath)
	})

	if err != nil {
		return SnapshotManifest{}, err
	}

	// Ensure directory closure: parent directories must exist as EntryDirectory
	for p := range entriesMap {
		segments := strings.Split(p, "/")
		for j := 1; j < len(segments); j++ {
			parent := strings.Join(segments[:j], "/")
			if _, exists := entriesMap[parent]; !exists {
				entriesMap[parent] = Entry{
					Path:     parent,
					Type:     EntryDirectory,
					ReadOnly: false,
				}
			}
		}
	}

	// Collect and sort entries
	entries := make([]Entry, 0, len(entriesMap))
	for _, e := range entriesMap {
		entries = append(entries, e)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Path < entries[j].Path
	})

	// Environment
	env := opts.Environment
	if env == nil {
		detected := detectEnvironment()
		env = &detected
	}

	now := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	body := ManifestBody{
		SnapshotID:         opts.SnapshotID,
		Kind:               opts.Kind,
		CapturedAt:         now,
		RunID:              opts.RunID,
		ParentSnapshotHash: opts.ParentSnapshotHash,
		Task:               opts.Task,
		Entries:            entries,
		Exclusions:         exclusions,
		Environment:        *env,
	}

	manifest, err := NewSnapshotManifest(body)
	if err != nil {
		return SnapshotManifest{}, fmt.Errorf("create snapshot manifest: %w", err)
	}
	return manifest, nil
}

func compileGlobMatcher(pattern string) (func(string) bool, error) {
	// Special cases:
	// "dir/**" matches "dir" and "dir/..."
	isDirPrefix := strings.HasSuffix(pattern, "/**")
	baseDir := ""
	if isDirPrefix {
		baseDir = pattern[:len(pattern)-3]
	}

	// Convert glob pattern to regular expression
	var b strings.Builder
	b.WriteString("^")
	i := 0
	for i < len(pattern) {
		if i+1 < len(pattern) && pattern[i:i+2] == "**" {
			if i+2 < len(pattern) && pattern[i+2] == '/' {
				b.WriteString("(?:.*/)?")
				i += 3
			} else {
				b.WriteString(".*")
				i += 2
			}
		} else if pattern[i] == '*' {
			b.WriteString("[^/]*")
			i++
		} else if pattern[i] == '?' {
			b.WriteString("[^/]")
			i++
		} else {
			ch := pattern[i]
			if strings.ContainsRune(`.+()|[]{}^$\`, rune(ch)) {
				b.WriteByte('\\')
			}
			b.WriteByte(ch)
			i++
		}
	}
	b.WriteString("$")

	re, err := regexp.Compile(b.String())
	if err != nil {
		return nil, err
	}

	return func(path string) bool {
		if isDirPrefix && (path == baseDir || strings.HasPrefix(path, baseDir+"/")) {
			return true
		}
		return re.MatchString(path)
	}, nil
}

func detectEnvironment() Environment {
	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "darwin"
	}
	archName := runtime.GOARCH

	caseSensitive := true
	if osName == "windows" || osName == "darwin" {
		caseSensitive = false
	}

	// Extract env var names
	rawEnv := os.Environ()
	varNames := make([]string, 0, len(rawEnv))
	seen := make(map[string]struct{}, len(rawEnv))
	for _, env := range rawEnv {
		idx := strings.Index(env, "=")
		if idx > 0 {
			name := env[:idx]
			if envVarRegex.MatchString(name) && len(name) <= 256 {
				if _, ok := seen[name]; !ok {
					seen[name] = struct{}{}
					varNames = append(varNames, name)
				}
			}
		}
	}
	sort.Strings(varNames)

	return Environment{
		OS:                       osName,
		Arch:                     archName,
		CaseSensitive:            caseSensitive,
		EnvironmentVariableNames: varNames,
		ToolchainEvidence:        []string{},
	}
}
