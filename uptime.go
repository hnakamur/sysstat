package sysstat

import (
	"os"
	"syscall"

	"github.com/hnakamur/ascii"
	"github.com/hnakamur/bytesconv"
)

// Uptime represents time elapsed from boot.
type Uptime struct {
	Uptime float64
}

type uptimeReader struct {
	buf [80]byte
}

var gUptimeReader uptimeReader

// ReadUptime read the uptime.
// Note ReadUptime is not goroutine safe.
func ReadUptime(u *Uptime) error {
	return gUptimeReader.readUptime(u)
}

func (r *uptimeReader) readUptime(u *Uptime) error {
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

func (r *uptimeReader) parse(buf []byte, u *Uptime) error {
	var err error
	u.Uptime, err = r.readFloat64(&buf)
	if err != nil {
		return err
	}
	return nil
}

func (r *uptimeReader) readFloat64(buf *[]byte) (float64, error) {
	start, end := ascii.NextField(*buf)
	field := (*buf)[start:end]
	val, err := bytesconv.ParseFloat(field, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}
