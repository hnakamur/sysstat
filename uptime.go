package sysstat

import (
	"bytes"
	"os"
	"syscall"
	"time"

	"github.com/hnakamur/ascii"
	"github.com/hnakamur/bytesconv"
)

// Uptime represents time elapsed from boot.
type Uptime struct {
	Uptime time.Duration
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
	u.Uptime, err = r.readDuration(&buf)
	if err != nil {
		return err
	}
	return nil
}

func (r *uptimeReader) readDuration(buf *[]byte) (time.Duration, error) {
	start, end := ascii.NextField(*buf)
	field := (*buf)[start:end]
	dot := bytes.IndexByte(field, '.')
	if dot == -1 {
		return 0, ErrUnexpectedFormat
	}
	seconds, err := bytesconv.ParseUint(field[:dot], 10, 64)
	if err != nil {
		return 0, err
	}
	centiSeconds, err := bytesconv.ParseUint(field[dot+1:], 10, 64)
	if err != nil {
		return 0, err
	}
	*buf = (*buf)[end+1:]
	return time.Duration(seconds*1000+centiSeconds*10) * time.Millisecond, nil
}
