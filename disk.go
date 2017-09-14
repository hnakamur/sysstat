package sysstat

import (
	"bytes"
	"time"

	"github.com/hnakamur/ascii"
	"github.com/hnakamur/bytesconv"
)

type ioStat struct {
	DevName            string
	ReadCountPerSec    float64
	ReadBytesPerSec    float64
	WrittenCountPerSec float64
	WrittenBytesPerSec float64
}

// /* major minor name rio rmerge rsect ruse wio wmerge wsect wuse running use aveq */
//     i = sscanf(line, "%u %u %s %lu %lu %lu %lu %lu %lu %lu %u %u %u %u",
//          &major, &minor, dev_name,
//          &rd_ios, &rd_merges_or_rd_sec, &rd_sec_or_wr_ios, &rd_ticks_or_wr_sec,
//          &wr_ios, &wr_merges, &wr_sec, &wr_ticks, &ios_pgr, &tot_ticks, &rq_ticks);

// DiskStat represents I/O statistics of block devices.
// https://www.kernel.org/doc/Documentation/ABI/testing/procfs-diskstats
// https://github.com/torvalds/linux/blob/486088bc4689f826b80aa317b45ac9e42e8b25ee/Documentation/iostats.txt
type diskStat struct {
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

type lastTwoDiskStats struct {
	devName []byte
	stats   [2]diskStat
}

type ioStatReader struct {
	buf      [8192]byte
	curr     int
	stats    []lastTwoDiskStats
	prevTime time.Time
}

var gIOStatReader ioStatReader

//func NewIOStatReader(devNames [][]byte) (*ioStatReader, error) {
//}

func (r *ioStatReader) switchCurr() {
	r.curr = 1 - r.curr
}

func (r *ioStatReader) parse(buf []byte, stats []lastTwoDiskStats) error {
	for len(buf) > 0 {
		line := ascii.GetLine(buf)
		start, end := ascii.NthField(line, 2)
		devName := line[start:end]

		for i := 0; i < len(stats); i++ {
			if bytes.Equal(devName, stats[i].devName) {
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

func (r *ioStatReader) parseLineAfterDevName(buf []byte, s *diskStat) error {
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
	if err != nil {
		return err
	}

	return nil
}

func readUint64Field(buf *[]byte) (uint64, error) {
	start, end := ascii.NextField(*buf)
	val, err := bytesconv.ParseUint((*buf)[start:end], 10, 64)
	if err != nil {
		return 0, err
	}
	*buf = (*buf)[end+1:]
	return val, nil
}

// common.h
// 128:#define S_VALUE(m,n,p)      (((double) ((n) - (m))) / (p) * HZ)

//  /*       rrq/s wrq/s   r/s   w/s  rsec  wsec  rqsz  qusz await r_await w_await svctm %util */
//   cprintf_f(2, 8, 2,
//       S_VALUE(ioj->rd_merges, ioi->rd_merges, itv),
//       S_VALUE(ioj->wr_merges, ioi->wr_merges, itv));
//   cprintf_f(2, 7, 2,
//       S_VALUE(ioj->rd_ios, ioi->rd_ios, itv),
//       S_VALUE(ioj->wr_ios, ioi->wr_ios, itv));
//   cprintf_f(4, 8, 2,
//       S_VALUE(ioj->rd_sectors, ioi->rd_sectors, itv) / fctr,
//       S_VALUE(ioj->wr_sectors, ioi->wr_sectors, itv) / fctr,
//       xds.arqsz,
//       S_VALUE(ioj->rq_ticks, ioi->rq_ticks, itv) / 1000.0);
//
// /* Structure for block devices statistics */
// struct stats_disk {
//         unsigned long long nr_ios       __attribute__ ((aligned (16)));
//         unsigned long rd_sect           __attribute__ ((aligned (16)));
//         unsigned long wr_sect           __attribute__ ((aligned (8)));
//         unsigned int rd_ticks           __attribute__ ((aligned (8)));
//         unsigned int wr_ticks           __attribute__ ((packed));
//         unsigned int tot_ticks          __attribute__ ((packed));
//         unsigned int rq_ticks           __attribute__ ((packed));
//         unsigned int major              __attribute__ ((packed));
//         unsigned int minor              __attribute__ ((packed));
// };

// /* Structure for block devices statistics */
// struct stats_disk {
//         unsigned long long nr_ios       __attribute__ ((aligned (16)));
//         unsigned long rd_sect           __attribute__ ((aligned (16)));
//         unsigned long wr_sect           __attribute__ ((aligned (8)));
//         unsigned int rd_ticks           __attribute__ ((aligned (8)));
//         unsigned int wr_ticks           __attribute__ ((packed));
//         unsigned int tot_ticks          __attribute__ ((packed));
//         unsigned int rq_ticks           __attribute__ ((packed));
//         unsigned int major              __attribute__ ((packed));
//         unsigned int minor              __attribute__ ((packed));
// };

// /* # of sectors read */
// unsigned long rd_sectors        __attribute__ ((aligned (8)));
// /* # of sectors written */
// unsigned long wr_sectors        __attribute__ ((packed));
// /* # of read operations issued to the device */
// unsigned long rd_ios            __attribute__ ((packed));
// /* # of read requests merged */
// unsigned long rd_merges         __attribute__ ((packed));
// /* # of write operations issued to the device */
// unsigned long wr_ios            __attribute__ ((packed));
// /* # of write requests merged */
// unsigned long wr_merges         __attribute__ ((packed));
// /* Time of read requests in queue */
// unsigned int  rd_ticks          __attribute__ ((packed));
// /* Time of write requests in queue */
// unsigned int  wr_ticks          __attribute__ ((packed));
// /* # of I/Os in progress */
// unsigned int  ios_pgr           __attribute__ ((packed));
// /* # of ticks total (for this device) for I/O */
// unsigned int  tot_ticks         __attribute__ ((packed));
// /* # of ticks requests spent in queue */
// unsigned int  rq_ticks          __attribute__ ((packed));
