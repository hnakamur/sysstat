package sysstat

import (
	"os"
	"syscall"
)

// Uptime represents time elapsed from boot.
type Uptime struct {
	Uptime float64
}

// UptimeReader is used for reading uptime.
// UptimeReader is not safe for concurrent accesses.
type UptimeReader struct {
	buf [80]byte
}

// NewUptimeReader creates a UptimeReader
func NewUptimeReader() *UptimeReader {
	return new(UptimeReader)
}

// Read reads the uptime
func (r *UptimeReader) Read(u *Uptime) error {
	fd, err := open([]byte("/proc/uptime"), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	n, err := syscall.Read(fd, r.buf[:])
	if err != nil {
		return err
	}
	return r.parse(r.buf[:n], u)
}

func (r *UptimeReader) parse(buf []byte, u *Uptime) error {
	var err error
	u.Uptime, err = readFloat64Field(&buf)
	if err != nil {
		return err
	}
	return nil
}
