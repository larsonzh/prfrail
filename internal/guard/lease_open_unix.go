//go:build !windows

package guard

import (
	"os"
	"syscall"
)

// openLeaseRecord opens one lease record for the bounded record read (P83 fix
// round 7, item T5#3).
//
// O_NOFOLLOW refuses a record that is a link: the read must describe the record the
// store wrote, never a target chosen by whoever can write the directory. O_NONBLOCK
// keeps an open on a FIFO from parking the caller while it holds the resource lock;
// it is a no-op for the regular file a record has to be, which the caller checks on
// the returned descriptor before reading anything.
func openLeaseRecord(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
}
