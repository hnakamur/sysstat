package sysstat

import (
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/hnakamur/ascii"
)

// NetworkStat is a statistics of network per device.
type NetworkStat struct {
	DevName            string
	RecvBytesPerSec    float64
	RecvPacketsPerSec  float64
	RecvErrsPerSec     float64
	RecvDropsPerSec    float64
	TransBytesPerSec   float64
	TransPacketsPerSec float64
	TransErrsPerSec    float64
	TransDropsPerSec   float64
	TransCollsPerSec   float64
}

// NetworkStat represents I/O statistics of block devices.
// https://github.com/torvalds/linux/blob/486088bc4689f826b80aa317b45ac9e42e8b25ee/Documentation/filesystems/proc.txt#L1152-L1169
type rawNetworkStat struct {
	// 2 - receive bytes
	RecvBytes uint64
	// 3 - receive packets
	RecvPackets uint64
	// 4 - receive errors
	RecvErrs uint64
	// 5 - receive drops
	RecvDrops uint64
	// 10 - transmit bytes
	TransBytes uint64
	// 11 - transmit packets
	TransPackets uint64
	// 12 - transmit errors
	TransErrs uint64
	// 13 - transmit drops
	TransDrops uint64
	// 14 - transmit collisions
	TransColls uint64
}

type lastTwoRawNetworkStats struct {
	devName string
	stats   [2]rawNetworkStat
}

// NetworkStatReader is used for reading disk statistics.
// NetworkStatReader is not safe for concurrent acceses from multiple goroutines.
type NetworkStatReader struct {
	buf      [8192]byte
	curr     int
	stats    []lastTwoRawNetworkStats
	prevTime time.Time
}

// NewNetworkStatReader creates a NetworkStatReader and does an initial read.
func NewNetworkStatReader(devNames []string) (*NetworkStatReader, error) {
	r := new(NetworkStatReader)
	r.allocStats(devNames)
	err := r.readNetworkStat(nil)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (r *NetworkStatReader) allocStats(devNames []string) {
	stats := make([]lastTwoRawNetworkStats, len(devNames))
	for i := 0; i < len(stats); i++ {
		stats[i].devName = devNames[i]
	}
	r.stats = stats
}

// Read reads network statistics.
func (r *NetworkStatReader) Read(stats []NetworkStat) error {
	return r.readNetworkStat(stats)
}

func (r *NetworkStatReader) readNetworkStat(stats []NetworkStat) error {
	fd, err := open([]byte("/proc/net/dev"), os.O_RDONLY, 0)
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
		err = r.fillNetworkStats(stats, intervalSeconds)
		if err != nil {
			return err
		}
	}
	r.prevTime = now
	r.switchCurr()
	return nil
}

func (r *NetworkStatReader) fillNetworkStats(stats []NetworkStat, intervalSeconds float64) error {
	for i := 0; i < len(stats); i++ {
		lastTwo := r.findLastTwoRawNetworkStats(stats[i].DevName)
		if lastTwo == nil {
			return errors.New("device name not found in disk stats")
		}

		r.fillNetworkStat(&stats[i], lastTwo, intervalSeconds)
	}
	return nil
}

func (r *NetworkStatReader) fillNetworkStat(s *NetworkStat, lastTwo *lastTwoRawNetworkStats, intervalSeconds float64) {
	c := &lastTwo.stats[r.curr]
	p := &lastTwo.stats[1-r.curr]
	s.RecvBytesPerSec = r.llSpValue(p.RecvBytes, c.RecvBytes, intervalSeconds)
	s.RecvPacketsPerSec = r.llSpValue(p.RecvPackets, c.RecvPackets, intervalSeconds)
	s.RecvErrsPerSec = r.llSpValue(p.RecvErrs, c.RecvErrs, intervalSeconds)
	s.RecvDropsPerSec = r.llSpValue(p.RecvDrops, c.RecvDrops, intervalSeconds)
	s.TransBytesPerSec = r.llSpValue(p.TransBytes, c.TransBytes, intervalSeconds)
	s.TransPacketsPerSec = r.llSpValue(p.TransPackets, c.TransPackets, intervalSeconds)
	s.TransErrsPerSec = r.llSpValue(p.TransErrs, c.TransErrs, intervalSeconds)
	s.TransDropsPerSec = r.llSpValue(p.TransDrops, c.TransDrops, intervalSeconds)
	s.TransCollsPerSec = r.llSpValue(p.TransColls, c.TransColls, intervalSeconds)
}

func (r *NetworkStatReader) llSpValue(v1, v2 uint64, intervalSeconds float64) float64 {
	if v2 < v1 {
		return 0
	}
	return float64(v2-v1) / intervalSeconds
}

func (r *NetworkStatReader) findLastTwoRawNetworkStats(devName string) *lastTwoRawNetworkStats {
	for i := 0; i < len(r.stats); i++ {
		if r.stats[i].devName == devName {
			return &r.stats[i]
		}
	}
	return nil
}

func (r *NetworkStatReader) switchCurr() {
	r.curr = 1 - r.curr
}

func (r *NetworkStatReader) parse(buf []byte, stats []lastTwoRawNetworkStats) error {
	for len(buf) > 0 {
		line := ascii.GetLine(buf)
		start, end := ascii.NextField(line)
		devName := line[start : end-1] // -1 for colon character in eth0:
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

func (r *NetworkStatReader) parseLineAfterDevName(buf []byte, s *rawNetworkStat) error {
	var err error
	s.RecvBytes, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.RecvPackets, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.RecvErrs, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.RecvDrops, err = readUint64Field(&buf)
	if err != nil {
		return err
	}

	s.TransBytes, err = readNthUint64Field(&buf, 4) // jump from 5 to 10
	if err != nil {
		return err
	}
	s.TransPackets, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.TransErrs, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.TransDrops, err = readUint64Field(&buf)
	if err != nil {
		return err
	}
	s.TransColls, err = readUint64Field(&buf)
	return err
}
