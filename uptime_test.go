package sysstat

import (
	"testing"
)

func TestUptimeReader_parse(t *testing.T) {
	buf := []byte("10654673.98 20455002.81\n")
	var u Uptime
	r := NewUptimeReader()
	err := r.parse(buf, &u)
	if err != nil {
		t.Fatal(err)
	}
	want := 10654673.98
	if u.Uptime != want {
		t.Errorf("uptime unmatch, got %g, want %g", u.Uptime, want)
	}
}

func BenchmarkUptimeReader_parse(b *testing.B) {
	buf := []byte("10654673.98 20455002.81\n")
	var u Uptime
	r := NewUptimeReader()
	for i := 0; i < b.N; i++ {
		err := r.parse(buf, &u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUptimeReader_Read(b *testing.B) {
	var u Uptime
	r := NewUptimeReader()
	for i := 0; i < b.N; i++ {
		err := r.Read(&u)
		if err != nil {
			b.Fatal(err)
		}
	}
}
