package sysstat

import "testing"

func TestCPUStatReader_parse(t *testing.T) {
	buf := []byte(`cpu  70688924 148688 17620091 2036888025 2351467 0 1426771 0 18331249 0
cpu0 35370168 74531 8761044 1019034230 1255954 0 443050 0 9172844 0
cpu1 35318756 74156 8859047 1017853794 1095513 0 983720 0 9158404 0
intr 4879759801 44 5693 0 0 0 0 0 0 1 0 0 0 4 0 0 0 62 0 0 118489261 0 0 0 0 0 151888559 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
ctxt 11431573079
btime 1494729939
processes 23742342
procs_running 3
procs_blocked 0
softirq 4630023313 4 1347808053 2165753 637294967 113231921 0 964818 1200133498 0 1328424299
`)
	var s rawCPUStat
	reader, err := NewCPUStatReader()
	if err != nil {
		t.Fatal(err)
	}
	err = reader.parse(buf, &s)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		name string
		ptr  *uint64
		want uint64
	}{
		{"User", &s.User, 70688924},
		{"Nice", &s.Nice, 148688},
		{"Sys", &s.Sys, 17620091},
		{"Idle", &s.Idle, 2036888025},
		{"IOWait", &s.IOWait, 2351467},
		{"HardIRQ", &s.HardIRQ, 0},
		{"SoftIRQ", &s.SoftIRQ, 1426771},
		{"Steal", &s.Steal, 0},
		{"Guest", &s.Guest, 18331249},
		{"GuestNice", &s.GuestNice, 0},
	}
	for _, c := range testCases {
		if *c.ptr != c.want {
			t.Errorf("cpu stat d %s unmatch, got %d, want %d", c.name, *c.ptr, c.want)
		}
	}
}
