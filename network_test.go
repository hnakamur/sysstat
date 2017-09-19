package sysstat

import "testing"

func TestNetworkStatReader_parse(t *testing.T) {
	buf := []byte(`Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
veth72350E: 8630270  113545    0    0    0     0          0         0 251802663  733421    0    0    0     0       0          0
vethY31EET: 183260208 1847380    0    0    0     0          0         0 15465774460 2985241    0    0    0     0       0          0
    lo: 17899045627 119002139    0    0    0     0          0         0 17899045627 119002139    0    0    0     0       0          0
veth0H4TIQ: 8387552  107860    0    0    0     0          0         0 246946911  725424    0    0    0     0       0          0
   br0: 329426402871 130478210    2    1    0     0          0         0 27152202131 88015716    3    4    5     0       0          0
enp0s25: 344775743869 253048085    0  139    0     0          0   3121901 29493822351 102872359    0    0    0     0       0          0
lxdbr0: 1388839751 14443788    0    0    0     0          0         0 31841956536 17300891    0    0    0     0       0          0
virbr0: 426736914 6484607    0    0    0     0          0         0 4196860926 2335725    0    0    0     0       0          0
lxcbr0: 18901220  288404    0    0    0     0          0         0 480430052  317139    0    0    0     0       0          0
vethTYTEJM: 103777555  356369    0    0    0     0          0         0 686545884 3270868    0    0    0     0       0          0
vethEBQIJM:  703932    4810    0    0    0     0          0         0 39280360  608558    0    0    0     0       0          0
veth881N5U:  847446    8167    0    0    0     0          0         0 39178176  605238    0    0    0     0       0          0
vethN5B6AL: 25420504  340655    0    0    0     0          0         0 679492038  986355    0    0    0     0       0          0
virbr0-nic:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
veth6NF9WL:  924275   11075    0    0    0     0          0         0 57471677  604281    0    0    0     0       0          0
docker0: 1087198656 9306285    0    0    0     0          0         0 17169034878 10455215    0    0    0     0       0          0
vethWD6643: 1781176   21853    0    0    0     0          0         0 69057145  611780    0    0    0     0       0          0
vethJ0LPGB: 451533637 3565555    0    0    0     0          0         0 1233918243 4449055    0    0    0     0       0          0`)
	// Note: To avoid actual disk read, we construct reader manually here
	reader := new(NetworkStatReader)
	reader.allocStats([]string{"br0", "enp0s25"})
	stats := make([]lastTwoRawNetworkStats, 2)
	stats[0].devName = "br0"
	stats[1].devName = "enp0s25"
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
		{0, "RecvBytes", &stats[0].stats[reader.curr].RecvBytes, 329426402871},
		{0, "RecvPackets", &stats[0].stats[reader.curr].RecvPackets, 130478210},
		{0, "RecvErrs", &stats[0].stats[reader.curr].RecvErrs, 2},
		{0, "RecvDrops", &stats[0].stats[reader.curr].RecvDrops, 1},
		{0, "TransBytes", &stats[0].stats[reader.curr].TransBytes, 27152202131},
		{0, "TransPackets", &stats[0].stats[reader.curr].TransPackets, 88015716},
		{0, "TransErrs", &stats[0].stats[reader.curr].TransErrs, 3},
		{0, "TransDrops", &stats[0].stats[reader.curr].TransDrops, 4},
		{0, "TransColls", &stats[0].stats[reader.curr].TransColls, 5},
		{1, "RecvBytes", &stats[1].stats[reader.curr].RecvBytes, 344775743869},
		{1, "RecvPackets", &stats[1].stats[reader.curr].RecvPackets, 253048085},
		{1, "RecvErrs", &stats[1].stats[reader.curr].RecvErrs, 0},
		{1, "RecvDrops", &stats[1].stats[reader.curr].RecvDrops, 139},
		{1, "TransBytes", &stats[1].stats[reader.curr].TransBytes, 29493822351},
		{1, "TransPackets", &stats[1].stats[reader.curr].TransPackets, 102872359},
		{1, "TransErrs", &stats[1].stats[reader.curr].TransErrs, 0},
		{1, "TransDrops", &stats[1].stats[reader.curr].TransDrops, 0},
		{1, "TransColls", &stats[1].stats[reader.curr].TransColls, 0},
	}
	for _, c := range testCases {
		if *c.ptr != c.want {
			t.Errorf("dev %d %s unmatch, got %d, want %d", c.devIndex, c.name, *c.ptr, c.want)
		}
	}
}

func BenchmarkNetworkStatReader_parse(b *testing.B) {
	buf := []byte(`Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
veth72350E: 8630270  113545    0    0    0     0          0         0 251802663  733421    0    0    0     0       0          0
vethY31EET: 183260208 1847380    0    0    0     0          0         0 15465774460 2985241    0    0    0     0       0          0
    lo: 17899045627 119002139    0    0    0     0          0         0 17899045627 119002139    0    0    0     0       0          0
veth0H4TIQ: 8387552  107860    0    0    0     0          0         0 246946911  725424    0    0    0     0       0          0
   br0: 329426402871 130478210    2    1    0     0          0         0 27152202131 88015716    3    4    5     0       0          0
enp0s25: 344775743869 253048085    0  139    0     0          0   3121901 29493822351 102872359    0    0    0     0       0          0
lxdbr0: 1388839751 14443788    0    0    0     0          0         0 31841956536 17300891    0    0    0     0       0          0
virbr0: 426736914 6484607    0    0    0     0          0         0 4196860926 2335725    0    0    0     0       0          0
lxcbr0: 18901220  288404    0    0    0     0          0         0 480430052  317139    0    0    0     0       0          0
vethTYTEJM: 103777555  356369    0    0    0     0          0         0 686545884 3270868    0    0    0     0       0          0
vethEBQIJM:  703932    4810    0    0    0     0          0         0 39280360  608558    0    0    0     0       0          0
veth881N5U:  847446    8167    0    0    0     0          0         0 39178176  605238    0    0    0     0       0          0
vethN5B6AL: 25420504  340655    0    0    0     0          0         0 679492038  986355    0    0    0     0       0          0
virbr0-nic:       0       0    0    0    0     0          0         0        0       0    0    0    0     0       0          0
veth6NF9WL:  924275   11075    0    0    0     0          0         0 57471677  604281    0    0    0     0       0          0
docker0: 1087198656 9306285    0    0    0     0          0         0 17169034878 10455215    0    0    0     0       0          0
vethWD6643: 1781176   21853    0    0    0     0          0         0 69057145  611780    0    0    0     0       0          0
vethJ0LPGB: 451533637 3565555    0    0    0     0          0         0 1233918243 4449055    0    0    0     0       0          0`)
	// Note: To avoid actual disk read, we construct reader manually here
	reader := new(NetworkStatReader)
	reader.allocStats([]string{"br0", "enp0s25"})
	stats := make([]lastTwoRawNetworkStats, 2)
	stats[0].devName = "br0"
	stats[1].devName = "enp0s25"
	for i := 0; i < b.N; i++ {
		err := reader.parse(buf, stats)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNetworkStatReader_Read(b *testing.B) {
	reader, err := NewNetworkStatReader([]string{"br0", "enp0s25"})
	if err != nil {
		b.Fatal(err)
	}

	stats := make([]NetworkStat, 2)
	stats[0].DevName = "br0"
	stats[1].DevName = "enp0s25"
	for i := 0; i < b.N; i++ {
		err := reader.Read(stats)
		if err != nil {
			b.Fatal(err)
		}
	}
}
