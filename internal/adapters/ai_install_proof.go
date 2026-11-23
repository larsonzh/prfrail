package adapters

// T3 (finding F6): a durable proof that the probe executable's INSTALLATION
// SOURCE is immutable.
//
// The P83 identity-plus-digest binding proves that the object that is executed is
// the object that was approved, but it only spans the window between the bound
// open and the spawn. It says nothing about who may REPLACE the object at that
// path between two checks (Windows) or rewrite its bytes in place while keeping
// the inode (Linux). This file owns the store-local carrier
// `ai-install-source-proof.json` and the platform-neutral observation that fills
// it:
//
//   - the carrier is written by `ai check` after the receiver is approved and
//     before the probe starts, and it is re-read - and re-verified LIVE - by the
//     check command and by the probe executor before every start;
//   - the carrier is a MEMO of a fact that was true when it was recorded, never a
//     verdict: every judgement is made from a fresh observation, so a stale
//     carrier cannot admit a probe;
//   - the observation is platform-specific (Windows DACL vs Unix mode/owner) and
//     every platform that cannot prove it fails closed.
//
// The body is store-local protocol: no schema and no wire record depends on it.
// Its bytes enter the proofrail:ai-install-source-proof:1\n domain digest, which
// is the only digest this file publishes.

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/larsonzh/prfrail/internal/evidence"
)

const (
	// AIInstallSourceProofFileName is the carrier name inside the trusted root
	// derived from the chain-config directory. A caller never chooses the path.
	AIInstallSourceProofFileName = "ai-install-source-proof.json"
	// AIInstallSourceProofByteBound bounds one carrier read. A carrier larger than
	// this is already far past anything this file produces, so an oversized one is
	// a refusal, never a silently truncated read.
	AIInstallSourceProofByteBound int64 = 1 << 20
	// AIInstallSourceProofDomain is the digest domain of the carrier bytes. The
	// platform-specific evidence digests use their own domains (declared in the
	// platform files), so a value observed on one platform can never be replayed
	// as another platform's fact.
	AIInstallSourceProofDomain = "proofrail:ai-install-source-proof:1\n"

	aiInstallSourceProofSchemaVersion = 1
	aiInstallSourcePlatformDomain     = "proofrail:ai-install-source-platform:1\n"
)

// ErrAIInstallSourceProof reports an installation source whose immutability could
// not be proven. It is raised before any credential is resolved and again before
// the spawn; the check command classifies it as cli_unavailable.
var ErrAIInstallSourceProof = errors.New("AI install source is not proven immutable")

// AIInstallSourcePlatform names the evaluation environment the proof was taken
// in. A proof is only valid for the platform that produced it.
type AIInstallSourcePlatform struct {
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	PlatformHash string `json:"platformHash"`
}

// AIInstallSourceWindowsFacts is the Windows observation: the DACL of the install
// directory (and of the executable itself) grants this user's token no write,
// delete or permission-changing right.
type AIInstallSourceWindowsFacts struct {
	DirACLDeniesWriteDelete bool   `json:"dirAclDeniesWriteDelete"`
	ACLEvidenceHash         string `json:"aclEvidenceHash"`
}

// AIInstallSourceInterpreter is the interpreter a script executable (or a dynamic
// ELF loader) delegates to. The kernel re-resolves that name at exec time, so the
// proof records the interpreter's identity and digest and every start re-verifies
// them.
//
// P85 F07: the chain is recursive, and every level of it is executed by the kernel
// in turn. A script whose interpreter is itself a script (or a dynamic loader) used
// to be verified to depth but only the FIRST level was recorded - the inner
// identity was checked and then discarded, so a carrier could not describe the fact
// the check had established. Interpreter carries that level, and a nil value means
// "this object is exec'd directly" (a static ELF object).
type AIInstallSourceInterpreter struct {
	Path        string                      `json:"path"`
	Hash        string                      `json:"hash"`
	FileID      string                      `json:"fileID"`
	Interpreter *AIInstallSourceInterpreter `json:"interpreter,omitempty"`
}

// AIInstallSourceLinuxFacts is the Linux observation: every component of the
// executable's path (install directory included) and of its shebang interpreter
// is unwritable by the executing user and not owned by that user.
type AIInstallSourceLinuxFacts struct {
	UnwritablePath        bool                        `json:"unwritablePath"`
	OwnershipEvidenceHash string                      `json:"ownershipEvidenceHash"`
	Interpreter           *AIInstallSourceInterpreter `json:"interpreter,omitempty"`
}

// AIInstallSourceImmutability carries the platform observation. Exactly one field
// is set, matching the proof's platform: os, so a carrier cannot claim a Windows
// fact on Linux or the other way round.
type AIInstallSourceImmutability struct {
	Windows *AIInstallSourceWindowsFacts `json:"windows,omitempty"`
	Linux   *AIInstallSourceLinuxFacts   `json:"linux,omitempty"`
}

// AIInstallSourceFacts is one live observation of an executable's installation
// source: what it is (path, digest, object identity), where it was observed, and
// why it cannot be replaced.
type AIInstallSourceFacts struct {
	ExecutablePath string
	ExecutableHash string
	FileID         string
	Platform       AIInstallSourcePlatform
	Immutability   AIInstallSourceImmutability
}

// AIInstallSourceProof is the persisted carrier. It is store-local protocol: no
// schema and no wire record depends on it.
type AIInstallSourceProof struct {
	SchemaVersion  int                         `json:"schemaVersion"`
	ExecutablePath string                      `json:"executablePath"`
	ExecutableHash string                      `json:"executableHash"`
	FileID         string                      `json:"fileID"`
	Platform       AIInstallSourcePlatform     `json:"platform"`
	Immutability   AIInstallSourceImmutability `json:"immutability"`
	RecordedAt     string                      `json:"recordedAt"`
}

// AIInstallSourcePlatformHash digests the running platform identity. It is the
// probe's own measurement of the evaluation environment.
func AIInstallSourcePlatformHash() string {
	return evidence.Digest(aiInstallSourcePlatformDomain, []byte(runtime.GOOS+"\n"+runtime.GOARCH))
}

// AIInstallSourcePlatformFor returns the platform identity of one os/arch pair.
func AIInstallSourcePlatformFor(osName, arch string) AIInstallSourcePlatform {
	return AIInstallSourcePlatform{
		OS:           osName,
		Arch:         arch,
		PlatformHash: evidence.Digest(aiInstallSourcePlatformDomain, []byte(osName+"\n"+arch)),
	}
}

// NewAIInstallSourceProof turns one live observation into the carrier body. The
// caller decides when to record it; this function never observes anything.
func NewAIInstallSourceProof(facts AIInstallSourceFacts, recordedAt time.Time) AIInstallSourceProof {
	return AIInstallSourceProof{
		SchemaVersion:  aiInstallSourceProofSchemaVersion,
		ExecutablePath: facts.ExecutablePath,
		ExecutableHash: facts.ExecutableHash,
		FileID:         facts.FileID,
		Platform:       facts.Platform,
		Immutability:   facts.Immutability,
		RecordedAt:     recordedAt.UTC().Format(TimestampLayout),
	}
}

// Facts is the observation one carrier describes. P85 F05: the carrier - not a
// fresh re-derivation of the executable at spawn time - is the normative
// description of the object the operator approved, so the console hands the
// carrier's facts to the executor and the executor compares its live observation
// against them. The mapping is a pure field copy; it observes nothing.
func (proof AIInstallSourceProof) Facts() AIInstallSourceFacts {
	return AIInstallSourceFacts{
		ExecutablePath: proof.ExecutablePath,
		ExecutableHash: proof.ExecutableHash,
		FileID:         proof.FileID,
		Platform:       proof.Platform,
		Immutability:   proof.Immutability,
	}
}

// EncodeAIInstallSourceProof returns the exact carrier bytes: canonical JSON plus
// the document newline, never a BOM.
func EncodeAIInstallSourceProof(proof AIInstallSourceProof) ([]byte, error) {
	canonical, err := evidence.EncodeCanonical(proof)
	if err != nil {
		return nil, fmt.Errorf("%w: the install source proof cannot be canonicalized: %v", ErrAIInstallSourceProof, err)
	}
	return append(append([]byte(nil), canonical...), '\n'), nil
}

// DecodeAIInstallSourceProof decodes and validates one carrier body. A body that
// is malformed, carries unknown fields or fails validation is an error: a carrier
// that cannot be read is a refusal, never "recompute and continue".
func DecodeAIInstallSourceProof(wire []byte) (AIInstallSourceProof, error) {
	var proof AIInstallSourceProof
	if err := evidence.DecodeStrictJSON(wire, &proof); err != nil {
		return AIInstallSourceProof{}, fmt.Errorf("%w: the install source proof is malformed: %v", ErrAIInstallSourceProof, err)
	}
	if err := proof.Validate(); err != nil {
		return AIInstallSourceProof{}, err
	}
	return proof, nil
}

// AIInstallSourceProofDigest is the domain digest of one carrier body. It answers
// "are these two carrier bodies byte-identical?"; it is never a verdict.
func AIInstallSourceProofDigest(wire []byte) string {
	return evidence.Digest(AIInstallSourceProofDomain, wire)
}

// Validate refuses a carrier that cannot describe a proven install source. The
// platform fact has to match the recorded platform exactly, so the carrier cannot
// mix a Windows claim with a Linux one.
func (proof AIInstallSourceProof) Validate() error {
	if proof.SchemaVersion != aiInstallSourceProofSchemaVersion {
		return fmt.Errorf("%w: unsupported schema version %d", ErrAIInstallSourceProof, proof.SchemaVersion)
	}
	if !filepath.IsAbs(proof.ExecutablePath) {
		return fmt.Errorf("%w: the recorded executable %q is not an absolute path", ErrAIInstallSourceProof, proof.ExecutablePath)
	}
	if !evidence.ValidHash(proof.ExecutableHash) || proof.FileID == "" {
		return fmt.Errorf("%w: the recorded executable carries no valid digest or identity", ErrAIInstallSourceProof)
	}
	if !validAIInstallSourcePlatform(proof.Platform) {
		return fmt.Errorf("%w: the recorded platform %q/%q is not valid", ErrAIInstallSourceProof, proof.Platform.OS, proof.Platform.Arch)
	}
	if _, err := time.Parse(TimestampLayout, proof.RecordedAt); err != nil {
		return fmt.Errorf("%w: the recorded timestamp %q is not valid: %v", ErrAIInstallSourceProof, proof.RecordedAt, err)
	}
	switch proof.Platform.OS {
	case "windows":
		if proof.Immutability.Windows == nil || proof.Immutability.Linux != nil {
			return fmt.Errorf("%w: a windows proof must carry exactly the windows facts", ErrAIInstallSourceProof)
		}
		if !proof.Immutability.Windows.DirACLDeniesWriteDelete {
			return fmt.Errorf("%w: the recorded windows facts do not deny write and delete", ErrAIInstallSourceProof)
		}
		if !evidence.ValidHash(proof.Immutability.Windows.ACLEvidenceHash) {
			return fmt.Errorf("%w: the recorded windows ACL evidence carries no valid digest", ErrAIInstallSourceProof)
		}
	case "linux":
		if proof.Immutability.Linux == nil || proof.Immutability.Windows != nil {
			return fmt.Errorf("%w: a linux proof must carry exactly the linux facts", ErrAIInstallSourceProof)
		}
		if !proof.Immutability.Linux.UnwritablePath {
			return fmt.Errorf("%w: the recorded linux facts do not prove an unwritable path", ErrAIInstallSourceProof)
		}
		if !evidence.ValidHash(proof.Immutability.Linux.OwnershipEvidenceHash) {
			return fmt.Errorf("%w: the recorded linux ownership evidence carries no valid digest", ErrAIInstallSourceProof)
		}
		if interpreter := proof.Immutability.Linux.Interpreter; interpreter != nil {
			if err := validateAIInstallSourceInterpreter(interpreter, 1); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("%w: %q cannot carry an install source proof", ErrAIInstallSourceProof, proof.Platform.OS)
	}
	return nil
}

// aiInstallSourceInterpreterDepthLimit bounds how many nested interpreter levels a
// carrier may describe. It matches the observation limit of the Linux walk: a
// chain deeper than that is refused there, so a deeper carrier could only be
// forged, and a forged deeper chain is a refusal.
const aiInstallSourceInterpreterDepthLimit = 8

// validateAIInstallSourceInterpreter refuses an incomplete or over-deep
// interpreter level. Every recorded level is normative: the executor compares the
// live observation against it level by level, so a level that cannot be described
// must not be admitted (P85 F07).
func validateAIInstallSourceInterpreter(interpreter *AIInstallSourceInterpreter, depth int) error {
	if interpreter == nil {
		return fmt.Errorf("%w: the recorded interpreter chain carries an empty level", ErrAIInstallSourceProof)
	}
	if depth > aiInstallSourceInterpreterDepthLimit {
		return fmt.Errorf("%w: the recorded interpreter chain is deeper than %d levels", ErrAIInstallSourceProof, aiInstallSourceInterpreterDepthLimit)
	}
	if !filepath.IsAbs(interpreter.Path) || !evidence.ValidHash(interpreter.Hash) || interpreter.FileID == "" {
		return fmt.Errorf("%w: the recorded shebang interpreter is incomplete", ErrAIInstallSourceProof)
	}
	if interpreter.Interpreter != nil {
		return validateAIInstallSourceInterpreter(interpreter.Interpreter, depth+1)
	}
	return nil
}

// Matches reports whether a live observation still describes the carrier. The
// executable path is compared the way the host filesystem compares it; every
// other field is compared exactly, because a changed digest, identity, platform or
// permission fact is exactly what the proof exists to detect.
func (facts AIInstallSourceFacts) Matches(proof AIInstallSourceProof) bool {
	return facts.MatchesFacts(proof.Facts())
}

// MatchesFacts compares two observations item by item under the same rules Matches
// uses. P85 F05 needs it because the spawn boundary compares a LIVE observation
// against the CARRIER's facts (the approved object), not against a carrier body.
func (facts AIInstallSourceFacts) MatchesFacts(other AIInstallSourceFacts) bool {
	if !aiInstallSourcePathEqual(facts.ExecutablePath, other.ExecutablePath) {
		return false
	}
	if facts.ExecutableHash != other.ExecutableHash || facts.FileID != other.FileID {
		return false
	}
	if facts.Platform != other.Platform {
		return false
	}
	return aiInstallSourceImmutabilityEqual(facts.Immutability, other.Immutability)
}

// VerifyAIInstallSource observes the installation source of one executable: the
// name has to be a stable absolute path, the object identity and digest come from
// ONE bound open, and the platform observation has to prove that this user cannot
// write, delete or re-permission the object.
//
// P85 F06: the leaf is opened the way the frozen-policy artifacts are opened - the
// derived directory is bound to a handle (evidence.OpenBoundDir) and the leaf name
// is resolved INSIDE that handle (evidence.OpenBoundDirFile, O_NOFOLLOW on Linux
// and an explicit reparse classification on Windows). A symlink or reparse point at
// the leaf therefore refuses instead of being followed, and the identity, digest
// and (on Linux) the interpreter/ELF header are all taken from that one handle.
func VerifyAIInstallSource(executable string) (AIInstallSourceFacts, error) {
	path := strings.TrimSpace(executable)
	if path == "" {
		return AIInstallSourceFacts{}, fmt.Errorf("%w: no executable was named", ErrAIInstallSourceProof)
	}
	if !filepath.IsAbs(path) {
		return AIInstallSourceFacts{}, fmt.Errorf("%w: %q is not a stable absolute path", ErrAIInstallSourceProof, path)
	}
	file, err := openInstallSourceLeaf(path)
	if err != nil {
		return AIInstallSourceFacts{}, err
	}
	defer func() { _ = file.Close() }()
	fileID, err := file.FileID()
	if err != nil {
		return AIInstallSourceFacts{}, fmt.Errorf("%w: identify the executable %q: %v", ErrAIInstallSourceProof, path, err)
	}
	digest, err := file.Hash(copilotCLIExecutableDomain)
	if err != nil {
		return AIInstallSourceFacts{}, fmt.Errorf("%w: digest the executable %q: %v", ErrAIInstallSourceProof, path, err)
	}
	immutability, err := verifyInstallSourceImmutability(path, file)
	if err != nil {
		return AIInstallSourceFacts{}, err
	}
	return AIInstallSourceFacts{
		ExecutablePath: path,
		ExecutableHash: digest,
		FileID:         fileID,
		Platform:       AIInstallSourcePlatformFor(runtime.GOOS, runtime.GOARCH),
		Immutability:   immutability,
	}, nil
}

// installSourceRequireCanonicalSpelling refuses any path whose spelling is not
// already the canonical, cleaned form (P85 F07, design rev 3 section 3).
//
// The kernel does not clean the name it executes: an exec of "/a/../b" walks "/a"
// (following a link there if one exists) and then "..", so the object that runs is
// NOT the object "/b" names. Resolving such a spelling ourselves would prove one
// name chain and hand the kernel another. The only honest rule for a chain this
// layer cannot describe byte for byte is to refuse it, so every non-canonical
// shape - a "." or ".." component, a doubled separator, a trailing separator
// outside the root, a NUL - is rejected before anything is opened.
//
// On a canonical spelling filepath.Clean/Dir/Base are all IDENTITY operations, so
// the callers below may keep using them without changing a single component: that
// is exactly what makes "the string this rule proves" and "the string the kernel
// walks" the same string.
func installSourceRequireCanonicalSpelling(path string) error {
	if path == "" {
		return fmt.Errorf("%w: an empty path names no install source", ErrAIInstallSourceProof)
	}
	if strings.ContainsRune(path, 0) {
		return fmt.Errorf("%w: %q carries a NUL byte, so its spelling cannot be canonical", ErrAIInstallSourceProof, path)
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("%w: %q is not a stable absolute path", ErrAIInstallSourceProof, path)
	}
	if cleaned := filepath.Clean(path); cleaned != path {
		return fmt.Errorf("%w: %q is not the canonical spelling %q, so the name chain this rule would prove is not the chain the kernel walks",
			ErrAIInstallSourceProof, path, cleaned)
	}
	return nil
}

// openInstallSourceLeaf opens one executable relative to a BOUND directory handle.
// The caller's path is only used to open the directory; the leaf name is resolved
// inside that handle, so a directory replaced between the two steps cannot redirect
// the open and a symlink or reparse point at the leaf is refused by the primitive.
//
// P85 F02/F06/F07: the directory handle itself is opened WITHOUT following a link
// at any component. The kernel walks the ORIGINAL name chain at exec time, so a
// chain that carries a symbolic link is a chain this rule cannot describe: the link
// target would be the object that runs while a different name was verified.
// Refusing the link is the fail-closed half; a link-free chain is verified exactly
// as the kernel would walk it.
//
// P85 R6 (review r3 finding): the binding is evidence.OpenBoundDirVerified, whose
// Windows half IS the root-chain walk (openBoundDirChainHandle): every component
// below the volume root is opened as itself and classified, and a reparse point -
// junction, symbolic link or any other tag - refuses. The previous call,
// evidence.OpenBoundDirNoFollow, fell through to a plain OpenBoundDir on Windows,
// which follows a reparse point and then judged the target. The identity returned
// by the SAME resolution is kept here, so the classification, the directory the
// leaf is resolved in and the digest below are one handle chain - not a
// classification followed by a separate by-name open.
//
// P85 F07 (fix round 3): the ORIGINAL spelling is what is split - never a cleaned
// one. The spelling is first required to be canonical, and the Dir/Base split is
// then asserted to be an identity (Dir+Base reproduce the original bytes), so a
// resolver can no longer fold "..", ".", doubled or trailing separators away
// before the directory is bound while the kernel still walks the folded-away form.
func openInstallSourceLeaf(path string) (*evidence.BoundFile, error) {
	if err := installSourceRequireCanonicalSpelling(path); err != nil {
		return nil, err
	}
	directory, leaf := filepath.Dir(path), filepath.Base(path)
	// On a canonical spelling this is an identity; the assertion makes it a proven
	// property rather than an assumption of filepath, and it fires if a future
	// refactor reintroduces a cleaning step upstream.
	if rebuilt := filepath.Join(directory, leaf); rebuilt != path {
		return nil, fmt.Errorf("%w: the split of %q into %q and %q is not an identity, so the opened chain would differ from the named one",
			ErrAIInstallSourceProof, path, directory, leaf)
	}
	boundDirectory, directoryID, err := evidence.OpenBoundDirVerified(directory)
	if err != nil {
		return nil, fmt.Errorf("%w: bind the directory of %q: %v", ErrAIInstallSourceProof, path, err)
	}
	defer func() { _ = boundDirectory.Close() }()
	file, err := evidence.OpenBoundDirFile(boundDirectory, leaf, copilotCLIExecutableMaxBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: open the executable %q inside its bound directory (object %s): %v", ErrAIInstallSourceProof, path, directoryID, err)
	}
	return file, nil
}

func validAIInstallSourcePlatform(platform AIInstallSourcePlatform) bool {
	if platform.OS == "" || platform.Arch == "" || !evidence.ValidHash(platform.PlatformHash) {
		return false
	}
	derived := evidence.Digest(aiInstallSourcePlatformDomain, []byte(platform.OS+"\n"+platform.Arch))
	return derived == platform.PlatformHash
}

// aiInstallSourcePathEqual compares two resolved paths the way the host
// filesystem does: case-insensitively on Windows, byte-exactly elsewhere.
func aiInstallSourcePathEqual(left, right string) bool {
	left, right = filepath.Clean(left), filepath.Clean(right)
	if runtime.GOOS == "windows" {
		return strings.EqualFold(left, right)
	}
	return left == right
}

// aiInstallSourceImmutabilityEqual compares two platform observations exactly,
// including the optional shebang interpreter.
func aiInstallSourceImmutabilityEqual(left, right AIInstallSourceImmutability) bool {
	if (left.Windows == nil) != (right.Windows == nil) {
		return false
	}
	if left.Windows != nil && *left.Windows != *right.Windows {
		return false
	}
	if (left.Linux == nil) != (right.Linux == nil) {
		return false
	}
	if left.Linux == nil {
		return true
	}
	if left.Linux.UnwritablePath != right.Linux.UnwritablePath ||
		left.Linux.OwnershipEvidenceHash != right.Linux.OwnershipEvidenceHash {
		return false
	}
	if (left.Linux.Interpreter == nil) != (right.Linux.Interpreter == nil) {
		return false
	}
	if left.Linux.Interpreter != nil && !aiInstallSourceInterpreterEqual(left.Linux.Interpreter, right.Linux.Interpreter) {
		return false
	}
	return true
}

// aiInstallSourceInterpreterEqual compares two interpreter chains level by level.
// The struct is compared field by field rather than with == because the nested
// pointer makes it non-comparable; a chain of different depth is a difference, so
// an inner level can never be dropped from the comparison (P85 F07).
func aiInstallSourceInterpreterEqual(left, right *AIInstallSourceInterpreter) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	if left.Path != right.Path || left.Hash != right.Hash || left.FileID != right.FileID {
		return false
	}
	return aiInstallSourceInterpreterEqual(left.Interpreter, right.Interpreter)
}
