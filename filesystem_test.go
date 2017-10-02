package sysstat

import (
	"testing"
)

func BenchmarkFileSystemStatReader_Read(b *testing.B) {
	r := NewFileSystemStatReader("/")
	var s FileSystemStat
	for i := 0; i < b.N; i++ {
		err := r.Read(&s)
		if err != nil {
			b.Fatal(err)
		}
	}
}
