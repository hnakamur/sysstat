package sysstat

import (
	"os"
	"strconv"
	"syscall"
)

// LoadAvg represents load averages for 1 minute, 5 minutes, and 15 minutes.
type LoadAvg struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

// ReadLoadAvg read the load average values.
func ReadLoadAvg(a *LoadAvg) error {
	fd, err := syscall.Open("/proc/loadavg", os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	var buf [80]byte
	n, err := syscall.Read(fd, buf[:])
	if err != nil {
		return err
	}
	return parseLoadAvg(string(buf[:n]), a)
}

func parseLoadAvg(s string, a *LoadAvg) error {
	start, end := nextField(s)
	load1, err := strconv.ParseFloat(s[start:end], 64)
	if err != nil {
		return err
	}

	s = s[end+1:]
	start, end = nextField(s)
	load5, err := strconv.ParseFloat(s[start:end], 64)
	if err != nil {
		return err
	}

	s = s[end+1:]
	start, end = nextField(s)
	load15, err := strconv.ParseFloat(s[start:end], 64)
	if err != nil {
		return err
	}

	a.Load1 = load1
	a.Load5 = load5
	a.Load15 = load15
	return nil
}
