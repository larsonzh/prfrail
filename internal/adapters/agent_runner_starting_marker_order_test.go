package adapters

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/larsonzh/prfrail/internal/guard"
)

// DR13-FIX-D3 pins one ordering contract of the starting marker: the published
// name must disappear *before* the host lock is dropped, so a concurrent scan can
// never observe "free lock + marker still at path" (the W1 state). The delivered
// order is os.Remove(claim.path) -> unlockStartingMarker(claim.file) ->
// claim.file.Close().
//
// The release path has no injection seam and unlockStartingMarker is a plain
// package function, so the ordering itself is pinned structurally on the AST while
// the runtime consequences (remove-while-open, idempotence, fail-closed on a
// refused delete) are pinned behaviourally.

const dr13StartingMarkerSourceFile = "agent_runner_starting_marker.go"

// dr13SourcePath locates the production source next to this test file, so the
// structural assertion follows the package instead of a hard-coded absolute path.
func dr13SourcePath(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate this test file to derive the package directory")
	}
	path := filepath.Join(filepath.Dir(file), dr13StartingMarkerSourceFile)
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("the release source %s must exist next to this test: %v", path, err)
	}
	return path
}

// dr13ReleaseFuncBody returns the body of (*startingMarkerClaim).release.
func dr13ReleaseFuncBody(t *testing.T, parsed *ast.File) *ast.BlockStmt {
	t.Helper()
	var body *ast.BlockStmt
	ast.Inspect(parsed, func(n ast.Node) bool {
		decl, ok := n.(*ast.FuncDecl)
		if !ok || decl.Body == nil || decl.Recv == nil || decl.Name.Name != "release" {
			return true
		}
		if len(decl.Recv.List) != 1 {
			return true
		}
		star, ok := decl.Recv.List[0].Type.(*ast.StarExpr)
		if !ok {
			return true
		}
		receiver, ok := star.X.(*ast.Ident)
		if !ok || receiver.Name != "startingMarkerClaim" {
			return true
		}
		body = decl.Body
		return false
	})
	if body == nil {
		t.Fatalf("cannot find (*startingMarkerClaim).release in %s", dr13StartingMarkerSourceFile)
	}
	return body
}

// dr13IsClaimPathExpr reports whether expr is exactly claim.path.
func dr13IsClaimPathExpr(expr ast.Expr) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok || selector.Sel == nil || selector.Sel.Name != "path" {
		return false
	}
	target, ok := selector.X.(*ast.Ident)
	return ok && target.Name == "claim"
}

// dr13ClaimPathAliases collects local names initialised from claim.path, so a
// rewrite such as "path := claim.path" keeps the structural test honest instead
// of forcing one exact spelling.
func dr13ClaimPathAliases(body *ast.BlockStmt) map[string]bool {
	aliases := map[string]bool{}
	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok || len(assign.Lhs) != len(assign.Rhs) {
			return true
		}
		for i, lhs := range assign.Lhs {
			name, ok := lhs.(*ast.Ident)
			if !ok || !dr13IsClaimPathExpr(assign.Rhs[i]) {
				continue
			}
			aliases[name.Name] = true
		}
		return true
	})
	return aliases
}

// dr13IsRemoveOfClaimPath reports whether call is os.Remove(<the claim path>).
func dr13IsRemoveOfClaimPath(call *ast.CallExpr, aliases map[string]bool) bool {
	selector, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || selector.Sel == nil || selector.Sel.Name != "Remove" {
		return false
	}
	pkg, ok := selector.X.(*ast.Ident)
	if !ok || pkg.Name != "os" {
		return false
	}
	if len(call.Args) != 1 {
		return false
	}
	if dr13IsClaimPathExpr(call.Args[0]) {
		return true
	}
	arg, ok := call.Args[0].(*ast.Ident)
	return ok && aliases[arg.Name]
}

// dr13StatementListing renders the top-level statements of a function body in
// source order with their line numbers, so a failure names which statement ran
// first instead of only reporting an offset.
func dr13StatementListing(fset *token.FileSet, body *ast.BlockStmt, lines []string) string {
	listing := make([]string, 0, len(body.List))
	for _, stmt := range body.List {
		position := fset.Position(stmt.Pos())
		text := ""
		if position.Line-1 >= 0 && position.Line-1 < len(lines) {
			text = strings.TrimSpace(lines[position.Line-1])
		}
		listing = append(listing, fmt.Sprintf("  line %d: %s", position.Line, text))
	}
	return strings.Join(listing, "\n")
}

// TestStartingMarkerReleaseRemovesBeforeUnlock is the deterministic red test of
// DR13-FIX-D3: in (*startingMarkerClaim).release the remove of the published path
// must appear before the unlock call in source order. The pre-fix order
// (unlock -> close -> remove) lets a scan take the free lock while the marker is
// still published, which is exactly the W1 state the slice closes.
//
// Honest scope: this is a source-order gate, not a runtime proof. The runtime
// consequences are covered by the behavioural tests below.
func TestStartingMarkerReleaseRemovesBeforeUnlock(t *testing.T) {
	path := dr13SourcePath(t)
	source, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	lines := strings.Split(string(source), "\n")
	fset := token.NewFileSet()
	parsed, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}
	body := dr13ReleaseFuncBody(t, parsed)
	aliases := dr13ClaimPathAliases(body)
	listing := dr13StatementListing(fset, body, lines)

	var removePos, unlockPos token.Pos
	var removeDeferred, unlockDeferred bool
	var stack []ast.Node
	ast.Inspect(body, func(n ast.Node) bool {
		if n == nil {
			if len(stack) > 0 {
				stack = stack[:len(stack)-1]
			}
			return true
		}
		stack = append(stack, n)
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		deferred := false
		for _, ancestor := range stack {
			switch ancestor.(type) {
			case *ast.DeferStmt, *ast.GoStmt:
				deferred = true
			}
		}
		if dr13IsRemoveOfClaimPath(call, aliases) && (removePos == token.NoPos || call.Pos() < removePos) {
			removePos = call.Pos()
			removeDeferred = deferred
		}
		if ident, ok := call.Fun.(*ast.Ident); ok && ident.Name == "unlockStartingMarker" && (unlockPos == token.NoPos || call.Pos() < unlockPos) {
			unlockPos = call.Pos()
			unlockDeferred = deferred
		}
		return true
	})

	if removePos == token.NoPos {
		t.Fatalf("release() must remove the published marker with os.Remove(claim.path); no such call was found.\nstatements:\n%s", listing)
	}
	if unlockPos == token.NoPos {
		t.Fatalf("release() must call unlockStartingMarker(claim.file); no such call was found.\nstatements:\n%s", listing)
	}
	if removeDeferred || unlockDeferred {
		t.Fatalf("release() must remove and unlock on the normal path, not deferred (remove deferred=%v, unlock deferred=%v).\nstatements:\n%s",
			removeDeferred, unlockDeferred, listing)
	}

	removeAt := fset.Position(removePos)
	unlockAt := fset.Position(unlockPos)
	if removeAt.Offset > unlockAt.Offset {
		removeLine := ""
		if removeAt.Line-1 >= 0 && removeAt.Line-1 < len(lines) {
			removeLine = strings.TrimSpace(lines[removeAt.Line-1])
		}
		unlockLine := ""
		if unlockAt.Line-1 >= 0 && unlockAt.Line-1 < len(lines) {
			unlockLine = strings.TrimSpace(lines[unlockAt.Line-1])
		}
		t.Fatalf("release() drops the host lock before the published name is gone: %q (line %d) runs before %q (line %d), so a scan can lock a marker that is still published (W1). Required order: remove -> unlock -> close.\nstatements:\n%s",
			unlockLine, unlockAt.Line, removeLine, removeAt.Line, listing)
	}
}

// TestStartingMarkerReleaseRemovesThePublishedPath pins the runtime share-mode
// premise of the new order: the claim handle is created with FILE_SHARE_DELETE,
// so the published marker can be removed while the claim still holds it open.
// A regression that narrows the share mode would turn every release into a
// silent "remove failed" residual.
func TestStartingMarkerReleaseRemovesThePublishedPath(t *testing.T) {
	runRoot := t.TempDir()
	claim, err := writeAgentRunnerStartingMarker(runRoot, newAgentRunnerStartingMarker(
		"launch-release-order",
		"request-release-order",
		guard.ProcessIdentity{PID: os.Getpid(), StartToken: "release-order-token"},
	))
	if err != nil {
		t.Fatalf("write the starting marker: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(claim.path) })
	if _, statErr := os.Stat(claim.path); statErr != nil {
		t.Fatalf("the marker must be published before release: %v", statErr)
	}
	if removed := claim.release(); !removed {
		t.Fatal("release must remove the published marker while its own handle is still open (the claim handle must share delete)")
	}
	if _, statErr := os.Stat(claim.path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the published marker must be gone after release (stat err=%v)", statErr)
	}
}

// TestStartingMarkerReleaseIsIdempotent pins the second half of the release
// contract: a claim released twice reports false the second time, never panics,
// and cannot resurrect or re-delete the name. The ordering change must not turn
// this into a nil-handle dereference or a second unlock on a closed handle.
func TestStartingMarkerReleaseIsIdempotent(t *testing.T) {
	runRoot := t.TempDir()
	claim, err := writeAgentRunnerStartingMarker(runRoot, newAgentRunnerStartingMarker(
		"launch-release-idempotent",
		"request-release-idempotent",
		guard.ProcessIdentity{PID: os.Getpid(), StartToken: "release-idempotent-token"},
	))
	if err != nil {
		t.Fatalf("write the starting marker: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(claim.path) })
	if removed := claim.release(); !removed {
		t.Fatal("the first release must remove the published marker")
	}
	if removed := claim.release(); removed {
		t.Fatal("a second release of the same claim must report false, never a second delete")
	}
	if _, statErr := os.Stat(claim.path); !errors.Is(statErr, os.ErrNotExist) {
		t.Fatalf("the marker must stay gone after the second release (stat err=%v)", statErr)
	}
}

// TestStartingMarkerReleaseNeverExposesReleasedContent is the probabilistic
// regression of DR13-FIX-D3 (registered at the R3 bar, not a deterministic
// proof). A writer repeatedly publishes and releases a claim while scanners run
// the recovery recipe - open, tryLock, Stat, decode - and no scan may observe
// "lock acquired + name still present at path + the record of a claim whose
// release already deleted the name". Such an observation is the W1 state.
//
// The scanners use plain os.OpenFile like RecoverStarting does. Two residuals are
// deliberately not counted, both registered as design R5: a delete refused because
// a scanner still holds a no-delete-share handle (release then reports false and
// the marker stays for the next scan), and a scan that locks an orphaned file
// after the name was replaced by a newer claim (os.SameFile tells the two apart;
// D2's Stat re-check covers it).
func TestStartingMarkerReleaseNeverExposesReleasedContent(t *testing.T) {
	const (
		dr13Releasing = iota + 1
		dr13Removed
		dr13Residual
	)
	const dr13Window = 3 * time.Second

	runRoot := t.TempDir()
	path := startingMarkerPath(runRoot)

	var mu sync.Mutex
	states := map[string]int{}
	observed := map[string]int{}
	forcedResiduals := 0

	deadline := time.Now().Add(dr13Window)
	stop := make(chan struct{})

	var scanners sync.WaitGroup
	for i := 0; i < 3; i++ {
		scanners.Add(1)
		go func() {
			defer scanners.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				file, err := os.OpenFile(path, os.O_RDWR, 0)
				if err != nil {
					if !errors.Is(err, os.ErrNotExist) {
						time.Sleep(time.Millisecond)
					}
					continue
				}
				locked, lockErr := tryLockStartingMarker(file)
				if lockErr != nil || !locked {
					_ = file.Close()
					continue
				}
				fdInfo, fdErr := file.Stat()
				pathInfo, statErr := os.Stat(path)
				if fdErr == nil && statErr == nil && os.SameFile(fdInfo, pathInfo) {
					if marker, decodeErr := decodeAgentRunnerStartingMarker(file); decodeErr == nil {
						mu.Lock()
						if states[marker.LaunchID] == dr13Releasing {
							observed[marker.LaunchID]++
						}
						mu.Unlock()
					}
				}
				unlockStartingMarker(file)
				_ = file.Close()
			}
		}()
	}

	iterations := 0
	for i := 0; time.Now().Before(deadline); i++ {
		marker := newAgentRunnerStartingMarker(
			fmt.Sprintf("launch-release-stress-%d", i),
			"request-release-stress",
			guard.ProcessIdentity{PID: os.Getpid(), StartToken: "release-stress-token"},
		)
		claim, err := writeAgentRunnerStartingMarker(runRoot, marker)
		if err != nil {
			// A scan holding the name without delete sharing refuses the rename;
			// retry instead of failing the run.
			time.Sleep(200 * time.Microsecond)
			continue
		}
		iterations++
		mu.Lock()
		states[marker.LaunchID] = dr13Releasing
		mu.Unlock()
		removed := claim.release()
		mu.Lock()
		if removed {
			states[marker.LaunchID] = dr13Removed
		} else {
			states[marker.LaunchID] = dr13Residual
			forcedResiduals++
		}
		mu.Unlock()
	}
	close(stop)
	scanners.Wait()
	_ = os.Remove(path)

	mu.Lock()
	defer mu.Unlock()
	var violations []string
	observations := 0
	for launchID, count := range observed {
		observations += count
		if states[launchID] == dr13Removed {
			violations = append(violations, launchID)
		}
	}
	sort.Strings(violations)
	if len(violations) > 0 {
		shown := violations
		if len(shown) > 5 {
			shown = shown[:5]
		}
		t.Fatalf("%d of %d release/scan interleavings let a scan lock a published marker whose own release then deleted the name (W1 observed between %d iterations, %d residual deletes); launch ids: %s",
			len(violations), observations, iterations, forcedResiduals, strings.Join(shown, ", "))
	}
	t.Logf("stress window clean: %d releases, %d free-lock observations during release, %d residual deletes", iterations, observations, forcedResiduals)
}
