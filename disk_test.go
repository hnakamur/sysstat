package sysstat

import "testing"

func TestDiskStatReader_parse(t *testing.T) {
	buf := []byte(`   7       0 loop0 10373 0 25194 476 0 0 0 0 0 124 476
   7       1 loop1 29322 0 63156 1128940 0 0 0 0 0 31244 1128936
   7       2 loop2 9086 0 21024 352 0 0 0 0 0 36 352
   7       3 loop3 18765 0 38148 2685988 0 0 0 0 0 81508 2685988
   7       4 loop4 22015 0 46504 21200 0 0 0 0 0 1120 21196
   7       5 loop5 190 0 1004 32 0 0 0 0 0 24 32
   7       6 loop6 20573 0 522236 3376 0 0 0 0 0 304 3376
   7       7 loop7 10243 0 22546 852 0 0 0 0 0 60 852
   8       0 sda 18115828 30368 4557439074 74915904 2432358 2630 480699421 5326092 0 42767900 80233776
   8       1 sda1 18115502 30368 4557434718 74913284 2405484 2630 480699421 2968916 0 40412152 77873980
   8      16 sdb 5544832 2815913 247978488 36492504 58421052 61635299 4134504954 1305616336 0 25602936 1343425972
   8      17 sdb1 3285 715 48488 1004 4204 2111 3698202 104560 0 10492 105560
   8      18 sdb2 124 0 138 12 0 0 0 0 0 12 12
   8      21 sdb5 5539944 2815198 247862418 36490540 51132600 61633188 4130805136 1304560156 0 24726036 1342364316
   8      32 sdc 18260428 34965 4556272723 31999648 2446024 2481 480699421 3506836 0 23124200 35499712
   8      33 sdc1 18260102 34965 4556268367 31998856 2419150 2481 480699421 2080560 0 21699464 34072652
  11       0 sr0 0 0 0 0 0 0 0 0 0 0 0
 252       0 dm-0 8392092 0 247857842 85348760 120144807 0 4130805136 1889017424 0 27119592 2011492480
 252       1 dm-1 4856165 0 219533850 21092312 109671043 0 4070864152 2020215052 0 26305540 2041433468
 252       2 dm-2 3535546 0 28307568 64269284 7492624 0 59940984 4164029824 0 2293272 4264445128
 252       3 dm-3 3535054 0 28292160 98918464 7492624 0 59940984 1936335464 0 12677308 2083039984
   7       8 loop8 1745 0 3670 36 0 0 0 0 0 4 36`)
	// Note: To avoid actual disk read, we construct reader manually here
	reader := new(DiskStatReader)
	reader.allocStats([]string{"sda", "sdb"})
	stats := make([]lastTwoRawDiskStats, 2)
	stats[0].devName = "sda"
	stats[1].devName = "sdb"
	err := reader.parse(buf, stats)
	if err != nil {
		t.Fatal(err)
	}
	testCases := []struct {
		devIndex int
		name     string
		ptr      *uint64
		want     uint64
	}{
		{0, "RdIOs", &stats[0].stats[reader.curr].RdIOs, 18115828},
		{0, "RdMerges", &stats[0].stats[reader.curr].RdMerges, 30368},
		{0, "RdSect", &stats[0].stats[reader.curr].RdSect, 4557439074},
		{0, "RdTicks", &stats[0].stats[reader.curr].RdTicks, 74915904},
		{0, "WrIOs", &stats[0].stats[reader.curr].WrIOs, 2432358},
		{0, "WrMerges", &stats[0].stats[reader.curr].WrMerges, 2630},
		{0, "WrSect", &stats[0].stats[reader.curr].WrSect, 480699421},
		{0, "WrTicks", &stats[0].stats[reader.curr].WrTicks, 5326092},
		{1, "RdIOs", &stats[1].stats[reader.curr].RdIOs, 5544832},
		{1, "RdMerges", &stats[1].stats[reader.curr].RdMerges, 2815913},
		{1, "RdSect", &stats[1].stats[reader.curr].RdSect, 247978488},
		{1, "RdTicks", &stats[1].stats[reader.curr].RdTicks, 36492504},
		{1, "WrIOs", &stats[1].stats[reader.curr].WrIOs, 58421052},
		{1, "WrMerges", &stats[1].stats[reader.curr].WrMerges, 61635299},
		{1, "WrSect", &stats[1].stats[reader.curr].WrSect, 4134504954},
		{1, "WrTicks", &stats[1].stats[reader.curr].WrTicks, 1305616336},
	}
	for _, c := range testCases {
		if *c.ptr != c.want {
			t.Errorf("dev %d %s unmatch, got %d, want %d", c.devIndex, c.name, *c.ptr, c.want)
		}
	}
}

func BenchmarkDiskStatReader_parse(b *testing.B) {
	buf := []byte(`   7       0 loop0 10373 0 25194 476 0 0 0 0 0 124 476
   7       1 loop1 29322 0 63156 1128940 0 0 0 0 0 31244 1128936
   7       2 loop2 9086 0 21024 352 0 0 0 0 0 36 352
   7       3 loop3 18765 0 38148 2685988 0 0 0 0 0 81508 2685988
   7       4 loop4 22015 0 46504 21200 0 0 0 0 0 1120 21196
   7       5 loop5 190 0 1004 32 0 0 0 0 0 24 32
   7       6 loop6 20573 0 522236 3376 0 0 0 0 0 304 3376
   7       7 loop7 10243 0 22546 852 0 0 0 0 0 60 852
   8       0 sda 18115828 30368 4557439074 74915904 2432358 2630 480699421 5326092 0 42767900 80233776
   8       1 sda1 18115502 30368 4557434718 74913284 2405484 2630 480699421 2968916 0 40412152 77873980
   8      16 sdb 5544832 2815913 247978488 36492504 58421052 61635299 4134504954 1305616336 0 25602936 1343425972
   8      17 sdb1 3285 715 48488 1004 4204 2111 3698202 104560 0 10492 105560
   8      18 sdb2 124 0 138 12 0 0 0 0 0 12 12
   8      21 sdb5 5539944 2815198 247862418 36490540 51132600 61633188 4130805136 1304560156 0 24726036 1342364316
   8      32 sdc 18260428 34965 4556272723 31999648 2446024 2481 480699421 3506836 0 23124200 35499712
   8      33 sdc1 18260102 34965 4556268367 31998856 2419150 2481 480699421 2080560 0 21699464 34072652
  11       0 sr0 0 0 0 0 0 0 0 0 0 0 0
 252       0 dm-0 8392092 0 247857842 85348760 120144807 0 4130805136 1889017424 0 27119592 2011492480
 252       1 dm-1 4856165 0 219533850 21092312 109671043 0 4070864152 2020215052 0 26305540 2041433468
 252       2 dm-2 3535546 0 28307568 64269284 7492624 0 59940984 4164029824 0 2293272 4264445128
 252       3 dm-3 3535054 0 28292160 98918464 7492624 0 59940984 1936335464 0 12677308 2083039984
   7       8 loop8 1745 0 3670 36 0 0 0 0 0 4 36`)
	// Note: To avoid actual disk read, we construct reader manually here
	reader := new(DiskStatReader)
	reader.allocStats([]string{"sda", "sdb"})
	stats := make([]lastTwoRawDiskStats, 2)
	stats[0].devName = "sda"
	stats[1].devName = "sdb"
	for i := 0; i < b.N; i++ {
		err := reader.parse(buf, stats)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDiskStatReader_Read(b *testing.B) {
	reader, err := NewDiskStatReader([]string{"sda", "sdb"})
	if err != nil {
		b.Fatal(err)
	}
	stats := make([]DiskStat, 2)
	stats[0].DevName = "sda"
	stats[1].DevName = "sdb"
	for i := 0; i < b.N; i++ {
		err := reader.Read(stats)
		if err != nil {
			b.Fatal(err)
		}
	}
}
