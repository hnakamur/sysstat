package sysstat

import (
	"bytes"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/hnakamur/ascii"
)

// CPUStat is a statistics about CPU.
type CPUStat struct {
	UserPercent   float64
	NicePercent   float64
	SysPercent    float64
	IOWaitPercent float64
}

// https://github.com/torvalds/linux/blob/486088bc4689f826b80aa317b45ac9e42e8b25ee/Documentation/filesystems/proc.txt#L1290-L1358
// https://github.com/torvalds/linux/blob/486088bc4689f826b80aa317b45ac9e42e8b25ee/Documentation/cpu-load.txt
type rawCPUStat struct {
	User      uint64
	Nice      uint64
	Sys       uint64
	Idle      uint64
	IOWait    uint64
	HardIRQ   uint64
	SoftIRQ   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
}

// CPUStatReader reads the CPU statistics.
type CPUStatReader struct {
	buf      [8192]byte
	curr     int
	stats    [2]rawCPUStat
	prevTime time.Time
	numCPU   int
}

// NewCPUStatReader creates a CPUStatReader.
func NewCPUStatReader() (*CPUStatReader, error) {
	r := &CPUStatReader{numCPU: runtime.NumCPU()}
	err := r.readCPUStat(nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *CPUStatReader) Read(s *CPUStat) error {
	return r.readCPUStat(s)
}

func (r *CPUStatReader) readCPUStat(s *CPUStat) error {
	fd, err := open([]byte("/proc/stat"), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	n, err := syscall.Read(fd, r.buf[:])
	if err != nil {
		return err
	}
	err = r.parse(r.buf[:n], &r.stats[r.curr])
	if err != nil {
		return err
	}

	now := time.Now()
	if s != nil {
		intervalSeconds := float64(now.Sub(r.prevTime)) / float64(time.Second)
		r.fillCPUStat(s, intervalSeconds)
	}
	r.prevTime = now
	r.switchCurr()
	return nil
}

func (r *CPUStatReader) switchCurr() {
	r.curr = 1 - r.curr
}

func (r *CPUStatReader) parse(buf []byte, s *rawCPUStat) error {
	for len(buf) > 0 {
		line := ascii.GetLine(buf)
		start, end := ascii.NextField(line)
		if bytes.Equal(line[start:end], []byte("cpu")) {
			err := r.parseLineAfterName(line[end+1:], s)
			if err != nil {
				return err
			}
			break
		}
		buf = buf[len(line):]
	}
	return nil
}

func (r *CPUStatReader) parseLineAfterName(buf []byte, s *rawCPUStat) error {
	var err error
	s.User, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.Nice, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.Sys, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.Idle, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.IOWait, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.HardIRQ, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.SoftIRQ, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.Steal, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.Guest, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.GuestNice, err = readUint64Field(&buf)
	return err
}

func (r *CPUStatReader) fillCPUStat(s *CPUStat, intervalSeconds float64) {
	curr := &r.stats[r.curr]
	prev := &r.stats[1-r.curr]
	s.UserPercent = r.calcUserPercent(curr, prev, intervalSeconds)
	s.NicePercent = r.calcNicePercent(curr, prev, intervalSeconds)
	s.SysPercent = r.calcSysPercent(curr, prev, intervalSeconds)
	s.IOWaitPercent = r.calcIOWaitPercent(curr, prev, intervalSeconds)
}

func (r *CPUStatReader) calcUserPercent(c, p *rawCPUStat, intervalSeconds float64) float64 {
	if c.User-c.Guest < p.User-p.Guest {
		return 0
	}
	return r.llSpValue(p.User-p.Guest, c.User-c.Guest, intervalSeconds)
}

func (r *CPUStatReader) calcNicePercent(c, p *rawCPUStat, intervalSeconds float64) float64 {
	if c.Nice-c.GuestNice < p.Nice-p.GuestNice {
		return 0
	}
	return r.llSpValue(p.Nice-p.GuestNice, c.Nice-c.GuestNice, intervalSeconds)
}

func (r *CPUStatReader) calcSysPercent(c, p *rawCPUStat, intervalSeconds float64) float64 {
	return r.llSpValue(p.Sys, c.Sys, intervalSeconds)
}

func (r *CPUStatReader) calcIOWaitPercent(c, p *rawCPUStat, intervalSeconds float64) float64 {
	return r.llSpValue(p.IOWait, c.IOWait, intervalSeconds)
}

func (r *CPUStatReader) llSpValue(v1, v2 uint64, intervalSeconds float64) float64 {
	// Workaround for CPU counters read from /proc/stat: Dyn-tick kernels
	// have a race issue that can make those counters go backward.
	if v2 < v1 {
		return 0
	}
	return float64(v2-v1) / intervalSeconds / float64(r.numCPU)
}
