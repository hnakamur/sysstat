// Package sysstat provides readers to monitor statistics about system,
// like CPU, memory, disk, network, and uptime.
//
// This package focuses on the performance.
//
//   * Only minimal metrics which we are actually interested are read.
//   * We try to minimize memory allocations during iterations of monitoring.
//
// Currently we achieve zero memory allocation count during interatoins of monitoring.
// Check it out with
//   go test -run NONE -bench . -benchmem
//
// This package supports only Linux.
package sysstat
