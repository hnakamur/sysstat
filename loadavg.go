package sysstat

import (
	"os"
	"syscall"

	"github.com/hnakamur/ascii"
	"github.com/hnakamur/bytesconv"
)

// LoadAvg represents load averages for 1 minute, 5 minutes, and 15 minutes.
type LoadAvg struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

var loadAvgBuf [80]byte

// ReadLoadAvg read the load average values.
// Note ReadLoadAvg is not goroutine safe.
func ReadLoadAvg(a *LoadAvg) error {
	fd, err := open([]byte("/proc/loadavg"), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	n, err := syscall.Read(fd, loadAvgBuf[:])
	if err != nil {
		return err
	}
	return parseLoadAvg(loadAvgBuf[:n], a)
}

func parseLoadAvg(s []byte, a *LoadAvg) error {
	start, end := ascii.NextField(s)
	load1, err := bytesconv.ParseFloat(s[start:end], 64)
	if err != nil {
		return err
	}

	s = s[end+1:]
	start, end = ascii.NextField(s)
	load5, err := bytesconv.ParseFloat(s[start:end], 64)
	if err != nil {
		return err
	}

	s = s[end+1:]
	start, end = ascii.NextField(s)
	load15, err := bytesconv.ParseFloat(s[start:end], 64)
	if err != nil {
		return err
	}

	a.Load1 = load1
	a.Load5 = load5
	a.Load15 = load15
	return nil
}
