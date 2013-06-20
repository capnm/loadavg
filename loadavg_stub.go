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

func LoadAvg2() ([3]float64, [3]int, error) {
	return [3]float64{}, [3]int{}, fmt.Errorf("LoadAvg2: unsupported platform %q", runtime.GOOS)
}

func LoadAvg3() ([3]float64, [3]int, error) {
	return [3]float64{}, [3]int{}, fmt.Errorf("LoadAvg3: unsupported platform %q", runtime.GOOS)
}

func Close() {}
