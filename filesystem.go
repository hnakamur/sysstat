package sysstat

import (
	"os"
	"syscall"
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
	fd, err := open(r.path, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	var buf syscall.Statfs_t
	err = syscall.Fstatfs(fd, &buf)
	if err != nil {
		return err
	}
	syscall.Close(fd)

	s.BlockSize = uint64(buf.Bsize)
	s.TotalBlocks = buf.Blocks
	s.FreeBlocks = buf.Bfree
	s.AvailableBlocks = buf.Bavail
	s.TotalINodes = buf.Files
	s.FreeINodes = buf.Ffree
	return nil
}
