package sysstat

import (
	"os"
	"strconv"
	"strings"
)

// LoadAvg represents load averages for 1 minute, 5 minutes, and 15 minutes.
type LoadAvg struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

// ReadLoadAvg read the load average values.
func ReadLoadAvg(a *LoadAvg) error {
	file, err := os.Open("/proc/loadavg")
	if err != nil {
		return err
	}
	defer file.Close()

	var buf [80]byte
	n, err := file.Read(buf[:])
	if err != nil {
		return err
	}
	return parseLoadAvg(string(buf[:n]), a)
}

func parseLoadAvg(b string, a *LoadAvg) error {
	fields := strings.Fields(b)
	load1, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return err
	}

	load5, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return err
	}

	load15, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return err
	}

	a.Load1 = load1
	a.Load5 = load5
	a.Load15 = load15
	return nil
}
