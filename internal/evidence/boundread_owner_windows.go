//go:build windows

package evidence

// Owner has no Windows answer: ownership is expressed through the DACL, which the
// Windows install-source observation reads directly (ai_install_proof_windows.go).
// ok=false makes a caller that expects unix ownership fail closed instead of
// inventing uid/gid facts.
func (bound *BoundFile) Owner() (uid, gid int, ok bool) {
	return 0, 0, false
}
