package sysstat

import (
	"syscall"
	"unsafe"
)

const _AT_FDCWD = -0x64

var _zero uintptr

func open(path []byte, mode int, perm uint32) (fd int, err error) {
	return openat(_AT_FDCWD, path, mode|syscall.O_LARGEFILE, perm)
}

func openat(dirfd int, path []byte, flags int, mode uint32) (fd int, err error) {
	var _p0 unsafe.Pointer
	if len(path) > 0 {
		_p0 = unsafe.Pointer(&path[0])
	} else {
		_p0 = unsafe.Pointer(&_zero)
	}
	r0, _, e1 := syscall.Syscall6(syscall.SYS_OPENAT, uintptr(dirfd), uintptr(unsafe.Pointer(_p0)), uintptr(flags), uintptr(mode), 0, 0)
	fd = int(r0)
	if e1 != 0 {
		err = e1
	}
	return
}
