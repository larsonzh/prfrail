package guard

import "errors"

// ErrExecutableBinding reports a spawn refused because the approved executable
// could not be proven to be the object about to be launched (P83 T3). It is
// always a refusal, never a fallback to an unverified path.
//
// The sentinel lives in an untagged file because the unsupported-platform start
// refuses with it before any binding helper is reachable.
var ErrExecutableBinding = errors.New("managed executable binding refused")
