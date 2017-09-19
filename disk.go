package sysstat

import (
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/hnakamur/ascii"
)

const sectorBytes = 512

// DiskStat is a statistics about disk.
type DiskStat struct {
	DevName            string
	ReadCountPerSec    float64
	ReadBytesPerSec    float64
	WrittenCountPerSec float64
	WrittenBytesPerSec float64
}

// DiskStat represents I/O statistics of block devices.
// https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
// https://github.com/torvalds/linux/blob/486088bc4689f826b80aa317b45ac9e42e8b25ee/Documentation/iostats.txt
type rawDiskStat struct {
	//  4 - reads completed successfully
	RdIOs uint64
	//  5 - reads merged
	RdMerges uint64
	//  6 - sectors read
	RdSect uint64
	//  7 - time spent reading (ms)
	RdTicks uint64
	//  8 - writes completed
	WrIOs uint64
	//  9 - writes merged
	WrMerges uint64
	// 10 - sectors written
	WrSect uint64
	// 11 - time spent writing (ms)
	WrTicks uint64
}

type lastTwoRawDiskStats struct {
	devName string
	stats   [2]rawDiskStat
}

// DiskStatReader is used for reading disk statistics.
// DiskStatReader is not safe for concurrent acceses from multiple goroutines.
type DiskStatReader struct {
	buf      [8192]byte
	curr     int
	stats    []lastTwoRawDiskStats
	prevTime time.Time
}

// NewDiskStatReader creates a DiskStatReader and does an initial read.
func NewDiskStatReader(devNames []string) (*DiskStatReader, error) {
	r := new(DiskStatReader)
	r.allocStats(devNames)
	err := r.readDiskStat(nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *DiskStatReader) allocStats(devNames []string) {
	stats := make([]lastTwoRawDiskStats, len(devNames))
	for i := 0; i < len(stats); i++ {
		stats[i].devName = devNames[i]
	}
	r.stats = stats
}

// Read reads statistics about disk.
func (r *DiskStatReader) Read(stats []DiskStat) error {
	return r.readDiskStat(stats)
}

func (r *DiskStatReader) readDiskStat(stats []DiskStat) error {
	fd, err := open([]byte("/proc/diskstats"), os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer syscall.Close(fd)

	n, err := syscall.Read(fd, r.buf[:])
	if err != nil {
		return err
	}
	err = r.parse(r.buf[:n], r.stats)
	if err != nil {
		return err
	}

	now := time.Now()
	if stats != nil {
		intervalSeconds := float64(now.Sub(r.prevTime)) / float64(time.Second)
		err = r.fillDiskStats(stats, intervalSeconds)
		if err != nil {
			return err
		}
	}
	r.prevTime = now
	r.switchCurr()
	return nil
}

func (r *DiskStatReader) fillDiskStats(stats []DiskStat, intervalSeconds float64) error {
	for i := 0; i < len(stats); i++ {
		lastTwo := r.findLastTwoRawDiskStats(stats[i].DevName)
		if lastTwo == nil {
			return errors.New("device name not found in disk stats")
		}

		r.fillDiskStat(&stats[i], lastTwo, intervalSeconds)
	}
	return nil
}

func (r *DiskStatReader) fillDiskStat(s *DiskStat, lastTwo *lastTwoRawDiskStats, intervalSeconds float64) {
	c := &lastTwo.stats[r.curr]
	p := &lastTwo.stats[1-r.curr]
	s.ReadCountPerSec = float64(c.RdIOs-p.RdIOs) / intervalSeconds
	s.ReadBytesPerSec = float64(c.RdSect-p.RdSect) * sectorBytes / intervalSeconds
	s.WrittenCountPerSec = float64(c.WrIOs-p.WrIOs) / intervalSeconds
	s.WrittenBytesPerSec = float64(c.WrSect-p.WrSect) * sectorBytes / intervalSeconds
}

func (r *DiskStatReader) findLastTwoRawDiskStats(devName string) *lastTwoRawDiskStats {
	for i := 0; i < len(r.stats); i++ {
		if r.stats[i].devName == devName {
			return &r.stats[i]
		}
	}
	return nil
}

func (r *DiskStatReader) switchCurr() {
	r.curr = 1 - r.curr
}

func (r *DiskStatReader) parse(buf []byte, stats []lastTwoRawDiskStats) error {
	for len(buf) > 0 {
		line := ascii.GetLine(buf)
		start, end := ascii.NthField(line, 2)
		devName := line[start:end]
		for i := 0; i < len(stats); i++ {
			if string(devName) == stats[i].devName {
				err := r.parseLineAfterDevName(line[end+1:], &stats[i].stats[r.curr])
				if err != nil {
					return err
				}
			}
		}

		buf = buf[len(line):]
	}
	return nil
}

func (r *DiskStatReader) parseLineAfterDevName(buf []byte, s *rawDiskStat) error {
	var err error
	s.RdIOs, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.RdMerges, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.RdSect, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.RdTicks, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.WrIOs, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.WrMerges, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.WrSect, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.WrTicks, err = readUint64Field(&buf)
	return err
}
