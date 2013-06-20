// +build darwin freebsd

package loadavg

import (
	"fmt"
	"syscall"
	"unsafe"
)

const sysctl = "vm.loadavg"

type loadavg struct {
	ldavg [3]uint32
	scale uint64
}

// LoadAvg returns the traditional 1, 5, and 15 min load averages.
func LoadAvg() ([3]float64, error) {
	v, err := syscall.Sysctl(sysctl)
	if err != nil {
		return [...]float64{-1, -1, -1}, fmt.Errorf("loadavg: sysctl failed: %v", err)
	}
	b := []byte(v)
	var l loadavg = *(*loadavg)(unsafe.Pointer(&b[0]))

	// BUG(dfc) syscall.Sysctl truncates the last byte (expecting a null terminated string)
	// so we have no access to the last byte returned. However it looks like on 64bit kernels
	// this byte is padding, so it all works anyway.

	scale := float64(l.scale)
	return [...]float64{
		float64(l.ldavg[0]) / scale,
		float64(l.ldavg[1]) / scale,
		float64(l.ldavg[2]) / scale,
	}, nil
}

// LoadAvg returns the traditional 1, 5, and 15 min load averages.
func LoadAvg2() ([3]float64, [3]int, error) {
	return [3]float64{}, [3]int{}, fmt.Errorf("loadavg: unsupported platform %q", runtime.GOOS)
}

// nop
func Close() error { return nil }
