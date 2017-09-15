package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/hnakamur/sysstat"
)

// sysstat.server1.disk.sda.rd_count
// sysstat.server1.disk.sda.rd_bytes
// sysstat.server1.disk.sda.wr_count
// sysstat.server1.disk.sda.wr_bytes

// sysstat.server1.mem.mem_total
// sysstat.server1.mem.mem_free
// sysstat.server1.mem.mem_avail
// sysstat.server1.mem.buffers
// sysstat.server1.mem.cached
// sysstat.server1.mem.swap_cached
// sysstat.server1.mem.swap_total
// sysstat.server1.mem.swap_free

func main() {
	dev := flag.String("dev", "sda", "disk device name (not partition)")
	interval := flag.Duration("interval", 10*time.Second, "interval")
	pprofAddr := flag.String("pprof-addr", ":8088", "pprof listen address")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	diskStatReader, err := sysstat.NewDiskStatReader([]string{*dev})
	if err != nil {
		log.Fatal(err)
	}
	diskStats := make([]sysstat.DiskStat, 1)
	diskStats[0].DevName = *dev

	cpuStatReader, err := sysstat.NewCPUStatReader()
	if err != nil {
		log.Fatal(err)
	}
	var cpuStat sysstat.CPUStat

	go func() {
		log.Println(http.ListenAndServe(*pprofAddr, nil))
	}()

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()
	for {
		t := <-ticker.C
		err = diskStatReader.Read(diskStats)
		if err != nil {
			log.Fatal(err)
		}
		err = cpuStatReader.Read(&cpuStat)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("t=%s, diskStats=%+v, cpuStat=%+v", t, diskStats, cpuStat)
	}
}
