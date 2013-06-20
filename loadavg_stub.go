// +build !linux,!darwin,!freebsd

package loadavg

import (
	"fmt"
	"runtime"
)

// LoadAvg returns the traditional 1, 5, and 15 min load averages.
func LoadAvg() ([3]float64, error) {
	return [...]float64{-1, -1, -1}, fmt.Errorf("loadavg: unsupported platform %q", runtime.GOOS)
}

// LoadAvg returns the traditional 1, 5, and 15 min load averages.
func LoadAvg2() ([3]float64, [3]int, error) {
	return [3]float64{}, [3]int{}, fmt.Errorf("loadavg: unsupported platform %q", runtime.GOOS)
}

// nop
func Close() error { return nil }
