package sysstat

import (
	"testing"
	"time"
)

func BenchmarkReadUptime(b *testing.B) {
	var u Uptime
	for i := 0; i < b.N; i++ {
		err := ReadUptime(&u)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func TestReadUptime(t *testing.T) {
	var u Uptime
	err := ReadUptime(&u)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUptimeReader_parse(t *testing.T) {
	buf := []byte("10654673.98 20455002.81\n")
	var u Uptime
	err := gUptimeReader.parse(buf, &u)
	if err != nil {
		t.Fatal(err)
	}
	want := 10654673980 * time.Millisecond
	if u.Uptime != want {
		t.Errorf("uptime unmatch, got %d, want %d", u.Uptime, want)
	}
}
