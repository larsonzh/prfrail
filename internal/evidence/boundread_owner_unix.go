//go:build unix

package evidence

import "syscall"

// Owner reports the owning user and group of the object this handle addresses.
//
// P85 F06: the Linux install-source rule used to take the leaf's ownership from
// an os.Lstat of the path, i.e. from a second resolution of the name, while the
// digest and identity came from the bound handle. The two could describe
// different objects. The owner is now read from the stat the BoundFile captured
// on its own open, so path, mode, owner, digest and identity are one observation.
func (bound *BoundFile) Owner() (uid, gid int, ok bool) {
	if bound == nil || bound.info == nil {
		return 0, 0, false
	}
	stat, ok := bound.info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, false
	}
	return int(stat.Uid), int(stat.Gid), true
}
