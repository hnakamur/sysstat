package sysstat

import (
	"os"
	"syscall"
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

func parseLoadAvg(buf []byte, a *LoadAvg) error {
	var err error
	a.Load1, err = readFloat64Field(&buf)
	if err != nil {
		return err
	}
	a.Load5, err = readFloat64Field(&buf)
	if err != nil {
		return err
	}
	a.Load15, err = readFloat64Field(&buf)
	if err != nil {
		return err
	}
	return nil
}
