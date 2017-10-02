package sysstat

import (
	"syscall"
	"unsafe"
)

// FileSystemStat is a statistics for a filesystem.
type FileSystemStat struct {
	BlockSize       uint64
	TotalBlocks     uint64
	FreeBlocks      uint64
	AvailableBlocks uint64
	TotalINodes     uint64
	FreeINodes      uint64
}

// FileSystemStatReader is used for reading uptime.
// FileSystemStatReader is not safe for concurrent accesses.
type FileSystemStatReader struct {
	path []byte
}

// NewFileSystemStatReader creates a FileSystemStatReader
func NewFileSystemStatReader(path string) *FileSystemStatReader {
	return &FileSystemStatReader{path: []byte(path)}
}

// Read reads the uptime
func (r *FileSystemStatReader) Read(s *FileSystemStat) error {
	var buf syscall.Statfs_t
	err := statfs(r.path, &buf)
	if err != nil {
		return err
	}

	s.BlockSize = uint64(buf.Bsize)
	s.TotalBlocks = buf.Blocks
	s.FreeBlocks = buf.Bfree
	s.AvailableBlocks = buf.Bavail
	s.TotalINodes = buf.Files
	s.FreeINodes = buf.Ffree
	return nil
}

func statfs(path []byte, buf *syscall.Statfs_t) (err error) {
	var _p0 unsafe.Pointer
	if len(path) > 0 {
		_p0 = unsafe.Pointer(&path[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	_, _, e1 := syscall.Syscall(syscall.SYS_STATFS, uintptr(_p0), uintptr(unsafe.Pointer(buf)), 0)
	if e1 != 0 {
		err = e1
	}
	return
}
