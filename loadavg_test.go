package sysstat

import "testing"

func TestParseLoadAvg(t *testing.T) {
	var a LoadAvg
	line := []byte("1.31 1.39 1.43 2/1081 24188\n")
	err := parseLoadAvg(line, &a)
	if err != nil {
		t.Fatal(err)
	}
	if a.Load1 != 1.31 {
		t.Errorf("Load1 unmatch, got=%g, want=%g", a.Load1, 1.31)
	}
	if a.Load5 != 1.39 {
		t.Errorf("Load5 unmatch, got=%g, want=%g", a.Load5, 1.39)
	}
	if a.Load15 != 1.43 {
		t.Errorf("Load15 unmatch, got=%g, want=%g", a.Load15, 1.43)
	}
}

func BenchmarkReadLoadAvg(b *testing.B) {
	var a LoadAvg
	for i := 0; i < b.N; i++ {
		err := ReadLoadAvg(&a)
		if err != nil {
			b.Fatal(err)
		}
	}
}
