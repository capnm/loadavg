// +build !linux,!darwin,!freebsd

package loadavg

import (
	"fmt"
	"runtime"
)

// LoadAvg returns the traditional 1, 5, and 15 min load averages, i.e.
// processes that are actually running – averaged over the last 1, 5, and 15 minutes.
func LoadAvg() ([3]float64, error) {
	return [3]float64{}, fmt.Errorf("LoadAvg: unsupported platform %q", runtime.GOOS)
}

func loadAvgSys() ([3]float64, [3]int, error) {
	return [3]float64{}, [3]int{}, nil
}
func loadAvgProc() ([3]float64, [3]int, error) {
	return [3]float64{}, [3]int{}, nil
}

func close() {}
