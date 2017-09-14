package sysstat

import "errors"

// ErrUnexpectedFormat is an error which is returned when the output format is unexpected.
var ErrUnexpectedFormat = errors.New("unexpected format")
