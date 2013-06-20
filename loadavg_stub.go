// +build !linux,!darwin,!freebsd

package loadavg

import (
	"fmt"
	"runtime"
)

// LoadAvg returns the traditional 1, 5, and 15 min load averages.
func LoadAvg() ([3]float64, error) {
	return [3]float64{}, fmt.Errorf("LoadAvg: unsupported platform %q", runtime.GOOS)
}
