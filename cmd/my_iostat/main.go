package main

import (
	"flag"
	"log"
	"time"

	"github.com/hnakamur/sysstat"
)

func main() {
	dev := flag.String("dev", "sda", "disk device name (not partition)")
	interval := flag.Duration("interval", 10*time.Second, "interval")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	diskStatReader := sysstat.NewDiskStatReader([]string{*dev})
	err := diskStatReader.InitialRead()
	if err != nil {
		log.Fatal(err)
	}
	diskStats := make([]sysstat.DiskStat, 1)
	diskStats[0].DevName = *dev

	ticker := time.NewTicker(*interval)
	defer ticker.Stop()
	for {
		t := <-ticker.C
		err = diskStatReader.Read(diskStats)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("t=%s, stats=%+v", t, diskStats)
	}
}
