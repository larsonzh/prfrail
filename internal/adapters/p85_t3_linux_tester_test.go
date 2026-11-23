//go:build linux

package adapters

// P85 T3 (finding F6) - the Linux half of the install source counterexamples.
//
// The Linux rule is: every component of the executable's absolute path (the
// installation directory included) and of its shebang interpreter must be
// non-writable by the executing user AND not owned by that user. The superuser can
// rewrite any path, so euid 0 can never prove it - which is why the positive path
// skips there instead of pretending.
//
// These tests are skipped on Windows by the build tag and, on a Linux run without a
// provable fixture, by an explicit skip.
//
// P85 r1: the same-inode rewrite leg is REGISTERED as an exemption (see its skip
// message) instead of a bare skip, and the ELF-magic counterexample left open by F3
// is measured below (TestP85T3LinuxNonELFWithoutShebangIsNotProven).
//
// P85 final review (F17): the exemption covers EXACTLY ONE leg - the same-inode
// rewrite that needs a PRIVILEGED fixture (an install directory that is unwritable
// AND foreign-owned so the proof succeeds, yet still mutable by the test). The other
// native counterexamples - the single-right refusals, the same-digest-different-object
// case and the interpreter-chain drift - are NOT exempted with it: they must run on
// the native Ubuntu VM (10.0.0.199) as a non-root user, and only a genuinely
// privileged fixture may be reported as not executed.
//
// P85 final review (F18): ② single A changed the interpreter rule to prove the
// ORIGINAL name chain (openat2 with RESOLVE_NO_SYMLINKS). A linked layout - the
// merged-/usr shape where /bin/sh or the PT_INTERP loader is a symlink - is now
// REFUSED, which is the fail-closed branch the review allows. The assertions below
// therefore expect a refusal on a linked chain instead of a proven target.

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// p85T3LinuxGuard skips when this build is not Linux or when the process is the
// superuser, whose three-component rule cannot prove anything.
func p85T3LinuxGuard(t *testing.T) {
	t.Helper()
	if runtime.GOOS != "linux" {
		t.Skipf("the linux install source rule does not apply on %s", runtime.GOOS)
	}
	if os.Geteuid() == 0 {
		t.Skip("euid 0 cannot prove immutability by design; run this leg as a non-root user")
	}
}

// TestP85T3LinuxRefusesWritableInstallDirectory is the "directory became writable"
// counterexample: a directory this user can write cannot protect a probe.
func TestP85T3LinuxRefusesWritableInstallDirectory(t *testing.T) {
	p85T3LinuxGuard(t)
	dir := t.TempDir()
	executable := filepath.Join(dir, "copilot")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := VerifyAIInstallSource(executable); !errors.Is(err, ErrAIInstallSourceProof) {
		t.Fatalf("a writable installation directory must be refused, got %v", err)
	}
}

// TestP85T3LinuxRefusesADirectoryThisUserOwns is the ownership half: removing the
// write bits is not enough while the directory still belongs to the executing user,
// because that user can restore them.
func TestP85T3LinuxRefusesADirectoryThisUserOwns(t *testing.T) {
	p85T3LinuxGuard(t)
	dir := t.TempDir()
	executable := filepath.Join(dir, "copilot")
	if err := os.WriteFile(executable, []byte("#!/bin/sh\nexit 0\n"), 0o555); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })
	if _, err := VerifyAIInstallSource(executable); !errors.Is(err, ErrAIInstallSourceProof) {
		t.Fatalf("a directory owned by the executing user must be refused even without write bits, got %v", err)
	}
}

// TestP85T3LinuxShebangInterpreterInWritableDirectoryRefused isolates the
// interpreter rule: the kernel re-resolves the shebang name at exec time, so an
// interpreter that lives in a writable directory makes the script unprovable.
//
// Mutation check: deleting the verifyUnwritableInstallPath call for the interpreter
// returns an interpreter fact instead of an error and turns this red.
func TestP85T3LinuxShebangInterpreterInWritableDirectoryRefused(t *testing.T) {
	p85T3LinuxGuard(t)
	dir := t.TempDir()
	// A real interpreter binary copied into a directory this user owns and can write.
	realInterpreter, err := os.ReadFile("/bin/sh")
	if err != nil {
		t.Skipf("no /bin/sh available to copy: %v", err)
	}
	interpreter := filepath.Join(dir, "sh")
	if err := os.WriteFile(interpreter, realInterpreter, 0o755); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(dir, "probe-script")
	if err := os.WriteFile(script, []byte("#!"+interpreter+"\nexit 0\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := installSourceInterpreter(script); err == nil {
		t.Fatalf("an interpreter in a writable directory must make the script unprovable")
	}
}

// TestP85T3LinuxNonELFWithoutShebangIsNotProven is the F3 counterexample added by the
// r1 review: before the fix a file carrying neither the ELF magic nor a shebang was
// silently classified as a plain object ("not a script") and passed. The header rule
// is called directly so the ELF-magic assertion itself is what is measured; going
// through VerifyAIInstallSource on a temporary fixture would also trip the
// writable-directory guard and make the leg over-determined.
//
// Mutation check: restoring the pre-fix branch that treats every non-shebang file as
// an executable object turns the plain-text and short-file cases red.
func TestP85T3LinuxNonELFWithoutShebangIsNotProven(t *testing.T) {
	p85T3LinuxGuard(t)
	dir := t.TempDir()
	candidates := map[string][]byte{
		"plain-text": []byte("this is an executable candidate with neither magic nor interpreter line\n"),
		"short-file": []byte("hi"),
	}
	for name, body := range candidates {
		path := filepath.Join(dir, "p85-t3-"+name)
		if err := os.WriteFile(path, body, 0o755); err != nil {
			t.Fatal(err)
		}
		interpreter, err := installSourceInterpreter(path)
		if err == nil {
			t.Fatalf("%s: a non-ELF file without a shebang must not be proven, got interpreter=%+v", name, interpreter)
		}
		if !errors.Is(err, ErrAIInstallSourceProof) {
			t.Fatalf("%s: the refusal must carry ErrAIInstallSourceProof, got %v", name, err)
		}
		// Fail-closed at the observation layer too. The writable temporary directory
		// fires there as well, so only the classification is asserted here.
		if _, observeErr := VerifyAIInstallSource(path); !errors.Is(observeErr, ErrAIInstallSourceProof) {
			t.Fatalf("%s: the install source observation must refuse as well, got %v", name, observeErr)
		}
	}
}

// p85T3LinuxPathHasSymlinkComponent reports whether any component of an absolute
// path is a symbolic link, without following it. It is the test-side oracle for
// "this layout is linked", so the expectation below is derived from the host and not
// asserted blindly.
func p85T3LinuxPathHasSymlinkComponent(path string) bool {
	cleaned := filepath.Clean(path)
	current := string(filepath.Separator)
	for _, part := range strings.Split(strings.TrimPrefix(cleaned, string(filepath.Separator)), string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return true
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return true
		}
	}
	return false
}

// p85T3LinuxLinkRefusalSet is the closed set of attribution texts that describe a
// link refusal by openat2 with RESOLVE_NO_SYMLINKS, spelled the way Go renders them.
const p85T3LinuxLinkRefusalSet = `"symbolic link", "ELOOP", "not a directory", "ENOTDIR"`

// p85T3LinuxLinkRefusalNamesLinkChain reports whether a refusal attributes the
// refusal to a link in the name chain.
//
// openat2 with RESOLVE_NO_SYMLINKS has TWO failure shapes for a linked component,
// and which one comes back depends on the host layout:
//
//   - ELOOP, rendered as "too many levels of symbolic links";
//   - ENOTDIR, rendered as "not a directory".
//
// P85 T3 minimal fix (VM evidence): the native Ubuntu 24.04 VM (10.0.0.199) uses
// the merged-/usr layout, so /bin is a symlink to usr/bin and openat2 answers
// `openat /bin: not a directory` instead of ELOOP. The refusal itself is correct -
// only the text shape differs from the one this leg used to demand. The assertion
// below therefore pins the refusal to this closed set instead of one spelling; the
// "the refusal must happen" half stays a hard assertion.
func p85T3LinuxLinkRefusalNamesLinkChain(err error) bool {
	if !errors.Is(err, ErrAIInstallSourceProof) {
		return false
	}
	message := err.Error()
	for _, shape := range []string{"symbolic link", "ELOOP", "not a directory", "ENOTDIR"} {
		if strings.Contains(message, shape) {
			return true
		}
	}
	return false
}

// TestP85T3LinuxELFObjectIsThePositiveControl keeps the F3 refusal honest: a real
// ELF object carries no shebang and is accepted by the header rule. The test binary
// itself is that ELF object, so no privileged fixture is needed.
//
// P85 F07: an ELF object is not automatically interpreter-free. A DYNAMIC object is
// exec'd through its PT_INTERP loader, and the kernel re-resolves that loader at exec
// time, so the observation must record (and immutability-verify) it; only a STATIC
// object is exec'd directly and records none.
//
// P85 F18: the loader is verified from its ORIGINAL name chain (openat2 with
// RESOLVE_NO_SYMLINKS), so on a merged-/usr host - where the PT_INTERP loader is a
// symlink - the correct outcome is a fail-closed REFUSAL, not a proof through the
// link target. The old assertion demanded a successful proof and would go red there.
// The host's own layout decides which expectation applies; the ELF header, not the
// filesystem, decides whether a loader exists at all.
func TestP85T3LinuxELFObjectIsThePositiveControl(t *testing.T) {
	p85T3LinuxGuard(t)
	self, err := os.Executable()
	if err != nil {
		t.Skipf("the test binary cannot be located: %v", err)
	}
	body, err := os.ReadFile(self)
	if err != nil {
		t.Fatalf("read the test binary %q: %v", self, err)
	}
	// The header rule is pure: it describes the object, not the permissions of the
	// directory the test binary happens to sit in (which this user owns, so the whole
	// path walk would refuse it for the writable-directory reason).
	loader, err := installSourceELFInterpreter(body, self)
	if err != nil {
		t.Fatalf("the ELF header of the test binary must be well formed, got %v", err)
	}
	if loader == "" {
		// Static object: the kernel execs it directly, so no interpreter is recorded.
		return
	}
	verified, verifyErr := installSourceVerifyInterpreter(loader, self, 0)
	if p85T3LinuxPathHasSymlinkComponent(loader) {
		if verifyErr == nil {
			t.Fatalf("a linked PT_INTERP chain (%q) must be refused, got %+v", loader, verified)
		}
		if !errors.Is(verifyErr, ErrAIInstallSourceProof) {
			t.Fatalf("the refusal must carry ErrAIInstallSourceProof, got %v", verifyErr)
		}
		if !strings.Contains(verifyErr.Error(), "symbolic link") {
			t.Fatalf("the refusal must name the link chain, not some other rule: %q", verifyErr)
		}
		return
	}
	if verifyErr != nil {
		t.Fatalf("a link-free PT_INTERP loader %q must be provable, got %v", loader, verifyErr)
	}
	if verified == nil {
		t.Fatalf("a dynamic ELF must record its PT_INTERP loader %q", loader)
	}
	if filepath.Clean(verified.Path) != filepath.Clean(loader) {
		t.Fatalf("the recorded interpreter must be the ORIGINAL loader name %q, got %q", loader, verified.Path)
	}
	if verified.Hash == "" || verified.FileID == "" {
		t.Fatalf("the recorded interpreter must carry its digest and object identity, got %+v", verified)
	}
}

// TestP85T3LinuxLinkedInterpreterChainIsRefused is the F18 native counterexample for
// the single-A rule: the kernel walks the ORIGINAL interpreter name chain, so a link
// anywhere in it cannot be described by the rule and must be refused instead of
// proven through its target.
//
// The expectation is derived from the host: on a merged-/usr Ubuntu run /bin/sh is a
// symlink and the refusal is the correct outcome; on a link-free layout the same call
// must succeed, which keeps the refusal from being satisfied by refusing everything.
//
// P85 T3 minimal fix: the refusal must HAPPEN (hard assertion) and must be attributed
// to the link rule, but the attribution text is taken from a closed set - openat2 with
// RESOLVE_NO_SYMLINKS reports a linked component as ELOOP ("symbolic link") on some
// layouts and as ENOTDIR ("not a directory") on the merged-/usr Ubuntu 24.04 VM, where
// /bin is a symlink to usr/bin and the VM run observed `openat /bin: not a directory`.
func TestP85T3LinuxLinkedInterpreterChainIsRefused(t *testing.T) {
	p85T3LinuxGuard(t)
	const interpreter = "/bin/sh"
	linked := p85T3LinuxPathHasSymlinkComponent(interpreter)
	proven, err := installSourceInterpreter(interpreter)
	if linked {
		if err == nil {
			t.Fatalf("a linked interpreter chain (%q) must be refused, got %+v", interpreter, proven)
		}
		if !errors.Is(err, ErrAIInstallSourceProof) {
			t.Fatalf("the refusal must carry ErrAIInstallSourceProof, got %v", err)
		}
		if !p85T3LinuxLinkRefusalNamesLinkChain(err) {
			t.Fatalf("the refusal must name the link chain (one of %s), not some other rule: %q", p85T3LinuxLinkRefusalSet, err)
		}
		return
	}
	if err != nil {
		t.Fatalf("a link-free interpreter %q must be provable, got %v", interpreter, err)
	}
	if proven == nil {
		t.Fatalf("the interpreter %q must be reported, got nil", interpreter)
	}
}

// TestP85T3LinuxSameInodeRewriteRequiresAPrivilegedFixture registers the one Linux
// counterexample that cannot be built from this process: it needs an install
// directory that is BOTH unwritable and foreign-owned (so the proof succeeds) and
// still mutable by the test (so an in-place rewrite can be attempted).
//
// P85 r1 F1 / P85 final review F17: this skip is an EXEMPTION REGISTRATION, not a
// silent gap. Its four elements are in the skip message itself so the test log
// carries them verbatim, and the same text is the entry for this slice's report and
// the ⑦ input package.
//
// F17 scope (final review): this exemption covers the SAME-INODE PRIVILEGED FIXTURE
// LEG ONLY. The other native counterexamples are NOT covered and must not be
// reported as exempt with it: the single-right refusals
// (TestP85T3LinuxRefusesWritableInstallDirectory,
// TestP85T3LinuxRefusesADirectoryThisUserOwns), the refusal of a linked interpreter
// chain (TestP85T3LinuxLinkedInterpreterChainIsRefused) and the interpreter drift
// assertions of TestP85T3LinuxELFObjectIsThePositiveControl. They run as a non-root
// user on the native Ubuntu VM together with ⑨; only a genuinely privileged fixture
// may be registered as not executed.
func TestP85T3LinuxSameInodeRewriteRequiresAPrivilegedFixture(t *testing.T) {
	p85T3LinuxGuard(t)
	t.Skip("EXEMPT-REGISTERED P85/F1-F17 | scope: the same-inode rewrite leg ONLY | reason: an end-to-end same-inode rewrite needs a fixture directory that is unwritable AND foreign-owned (so the proof succeeds) yet still mutable by the test (so the rewrite can be attempted); building it needs root chown - the ownership rule is part of the ②b-2 verdict - and the VM's larson account has no passwordless sudo, so the privileged fixture is not constructible in this round | not exempted with it: the single-right refusals, the same-digest-different-object case and the interpreter-chain drift, which must run natively as a non-root user on the Ubuntu VM with ⑨ | carried-by: TestP85T3WritableInstallDirectoryIsNotProven (a writable installation directory is refused) and TestP85T3MatchesDetectsAChangedObservation (the same fileID with a changed digest stops matching, which is what an in-place rewrite produces) | destination: this slice's report and the ⑦ review input package")
}
