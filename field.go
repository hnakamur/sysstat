package sysstat

import (
	"github.com/hnakamur/ascii"
	"github.com/hnakamur/bytesconv"
)

func readFloat64Field(buf *[]byte) (float64, error) {
	start, end := ascii.NextField(*buf)
	val, err := bytesconv.ParseFloat((*buf)[start:end], 64)
	if err != nil {
		return 0, err
	}
	*buf = (*buf)[end+1:]
	return val, nil
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
