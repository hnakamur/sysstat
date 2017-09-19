// Package sysstat provides functionalities to monitor system statistics
// like CPU, memory, disk, network, and uptime.
//
// This package focuses on the performance.
//
//   * Only minimal metrics which we are actually interested in are read.
//   * We try to minimize memory allocations during iterations of monitoring.
//
// Currently we achieve zero memory allocation count during interations of monitoring.
// Check it out with
//   go test -run NONE -bench . -benchmem
//
// This package supports only Linux.
package sysstat
