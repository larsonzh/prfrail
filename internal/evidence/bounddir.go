package evidence

// P42-T027B-EVIDENCE-BIND: the Q2 ancestor-directory binding primitive.
//
// P83 SECURITY-FINDINGS-REMEDIATION (X-1) moved this primitive out of
// internal/adapters and into this package. The gates runner binds a workspace
// directory handle before it collects artifacts (F1), and internal/adapters
// already imports internal/gates, so the container has to live below both: the
// primitive stays next to BoundFile, which it returns.
//
// docs/validation/evidence/P42-T027B-EVIDENCE-BIND/01-arch-design/review.md
// section 2 "Q2 ancestor directory replacement" requires a directory handle that
// outlives the path: once the per-request evidence directory is bound, an
// ancestor renamed away and re-created by an attacker must not be able to
// redirect a read. Both platform implementations resolve a relative artifact
// name against the bound directory object, never against a re-resolved path:
//
//   - Windows: CreateFile with FILE_FLAG_BACKUP_SEMANTICS opens the directory
//     object; a relative open derives the object's current real name with
//     GetFinalPathNameByHandle and opens under it.
//   - Unix: openat2(2) with RESOLVE_BENEATH resolves relative to the directory
//     descriptor, so a replaced ancestor cannot redirect it and a symlink in any
//     component is refused.
//
// BoundDir, OpenBoundDir and OpenBoundDirFile are the shared surface; the
// platform files own the mechanics and the post-open verification. The relative
// name guard lives in bounddir_clean.go, which the darwin build excludes because
// its fail-closed half resolves no names.
//
// internal/adapters keeps thin aliases (see bounddir_alias.go there) so every
// existing caller and every existing BoundDir test keeps working unchanged.
//
// P83 fix round 7 (item T4b) added BoundDirNames: the enumeration an artifact glob
// needs is read from a descriptor resolved beneath the bound one, never from the
// caller's path, so a directory replaced after the bind cannot contribute a name
// that the reads would then never see.

import "errors"

// ErrBoundDirUnsupported reports that directory binding is unavailable on this
// platform. It is defined here, in a file built on every platform, because the
// adapters alias refers to it on all of them.
var ErrBoundDirUnsupported = errors.New("evidence: directory binding is supported on windows and linux only")
