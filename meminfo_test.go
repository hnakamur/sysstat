package sysstat

import (
	"testing"
)

func BenchmarkReadMemInfo(b *testing.B) {
	var m MemInfo
	for i := 0; i < b.N; i++ {
		err := ReadMemInfo(&m)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestReadMemInfo(t *testing.T) {
	var m MemInfo
	err := ReadMemInfo(&m)
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseMemInfo(t *testing.T) {
	buf := []byte(`
MemTotal:       16260508 kB
MemFree:          543220 kB
MemAvailable:    6990124 kB
Buffers:         1931976 kB
Cached:          5597668 kB
SwapCached:         7168 kB
Active:          6838880 kB
Inactive:        2903424 kB
Active(anon):    1811632 kB
Inactive(anon):  1788208 kB
Active(file):    5027248 kB
Inactive(file):  1115216 kB
Unevictable:        7120 kB
Mlocked:            7120 kB
SwapTotal:      16600572 kB
SwapFree:       16582724 kB
Dirty:               692 kB
Writeback:             0 kB
AnonPages:       2212660 kB
Mapped:           461952 kB
Shmem:           1383156 kB
Slab:             823872 kB
SReclaimable:     642376 kB
SUnreclaim:       181496 kB
KernelStack:       16384 kB
PageTables:        42036 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:    24730824 kB
Committed_AS:   13751760 kB
VmallocTotal:   34359738367 kB
VmallocUsed:           0 kB
VmallocChunk:          0 kB
HardwareCorrupted:     0 kB
AnonHugePages:    569344 kB
CmaTotal:              0 kB
CmaFree:               0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:     1521344 kB
DirectMap2M:    15083520 kB
`)
	var m MemInfo
	err := gMemInfoReader.parseMemInfo(buf, &m)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name   string
		ptr    *uint64
		wantKB uint64
	}{
		{"MemTotal", &m.MemTotal, 16260508},
		{"MemFree", &m.MemFree, 543220},
		{"MemAvailable", &m.MemAvailable, 6990124},
		{"Buffers", &m.Buffers, 1931976},
		{"Cached", &m.Cached, 5597668},
		{"SwapCached", &m.SwapCached, 7168},
		{"SwapTotal", &m.SwapTotal, 16600572},
		{"SwapFree", &m.SwapFree, 16582724},
	}
	for _, c := range testCases {
		if *c.ptr != c.wantKB*1024 {
			t.Errorf("%s unmatch, got %d, want %d", c.name, *c.ptr, c.wantKB*1024)
		}
	}
}
