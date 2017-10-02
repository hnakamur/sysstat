sysstat
=======

sysstat is a Go package which provides functionalities to monitor system statistics
like CPU, memory, disk, network, and uptime.

This package focuses on the performance, low CPU and memory footprint, to give
least impact on system when we are monitoring.

* We monitor only minimal metrics which we are actually interested in.
* We try to minimize memory allocations during iterations of monitoring in order to
  reduce load from Go runtime GC.

Currently we achieve zero memory allocation during interations of monitoring
for CPUStatReader, DiskStatReader, LoadAvgReader, MemoryStatReader, NetworkStatReader,
and UptimeReader.

Check it out with

```
go test -run NONE -bench . -benchmem
```

This package supports only Linux.
