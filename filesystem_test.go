package sysstat

import (
	"testing"
)

func TestFileSystemStatReader_Read(t *testing.T) {
	paths := []string{"/"}
	r := NewFileSystemStatReader(paths)
	stats := make([]FileSystemStat, len(paths))
	err := r.Read(stats)
	if err != nil {
		t.Fatal(err)
	}
}

func BenchmarkFileSystemStatReader_Read(b *testing.B) {
	paths := []string{"/"}
	r := NewFileSystemStatReader(paths)
	stats := make([]FileSystemStat, len(paths))
	for i := 0; i < b.N; i++ {
		err := r.Read(stats)
		if err != nil {
			b.Fatal(err)
		}
	}
}
