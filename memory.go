package sysstat

import (
	"bytes"
	"os"
	"syscall"

	"github.com/hnakamur/ascii"
	"github.com/hnakamur/bytesconv"
)

// MemoryStat represents memory statistics in bytes.
// Only interested fields are provided for performance reason.
type MemoryStat struct {
	MemTotal     uint64
	MemFree      uint64
	MemAvailable uint64
	Buffers      uint64
	Cached       uint64
	SwapCached   uint64
	SwapTotal    uint64
	SwapFree     uint64
}

// MemoryStatReader is used for reading memory statistics.
// MemoryStatReader is not safe for concurrent accesses from multiple goroutines.
type MemoryStatReader struct {
	buf [4096]byte
}

// NewMemoryStatReader crates a MemoryStatReader.
func NewMemoryStatReader() *MemoryStatReader {
	return new(MemoryStatReader)
}

// Read reads a statistics about memory.
func (r *MemoryStatReader) Read(m *MemoryStat) error {
	fd, err := open([]byte("/proc/meminfo"), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	n, err := syscall.Read(fd, r.buf[:])
	if err != nil {
		return err
	}
	return r.parse(r.buf[:n], m)
}

func (r *MemoryStatReader) parse(buf []byte, m *MemoryStat) error {
	var err error
	m.MemTotal, err = r.readValue(&buf, []byte("MemTotal:"))
	if err != nil {
		return err
	}
	m.MemFree, err = r.readValue(&buf, []byte("MemFree:"))
	if err != nil {
		return err
	}
	m.MemAvailable, err = r.readValue(&buf, []byte("MemAvailable:"))
	if err != nil {
		return err
	}
	m.Buffers, err = r.readValue(&buf, []byte("Buffers:"))
	if err != nil {
		return err
	}
	m.Cached, err = r.readValue(&buf, []byte("Cached:"))
	if err != nil {
		return err
	}
	m.SwapCached, err = r.readValue(&buf, []byte("SwapCached:"))
	if err != nil {
		return err
	}
	m.SwapTotal, err = r.readValue(&buf, []byte("SwapTotal:"))
	if err != nil {
		return err
	}
	m.SwapFree, err = r.readValue(&buf, []byte("SwapFree:"))
	return err
}

func (r *MemoryStatReader) readValue(buf *[]byte, prefix []byte) (uint64, error) {
	line := r.findLineByPrefix(*buf, prefix)
	if line == nil {
		return 0, ErrUnexpectedFormat
	}
	start, end := ascii.NthField(line, 1)
	val, err := bytesconv.ParseUint(line[start:end], 10, 64)
	if err != nil {
		return 0, err
	}
	*buf = (*buf)[len(line):]
	return val * 1024, nil
}

func (r *MemoryStatReader) findLineByPrefix(buf, prefix []byte) []byte {
	for len(buf) > 0 {
		line := ascii.GetLine(buf)
		if bytes.HasPrefix(line, prefix) {
			return line
		}
		buf = buf[len(line):]
	}
	return nil
}
