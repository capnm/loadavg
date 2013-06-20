package loadavg

import (
	"fmt"
	"os"
)

const procfile = "/proc/loadavg"

var f *os.File = nil

// The first three numbers represent the number of active tasks on the system
//  	– processes that are actually running – averaged over the last 1, 5, and 15 minutes.
//  The next entry shows the instantaneous current number of runnable tasks
//	– processes that are currently scheduled to run rather than being blocked in a system call
//	– and the total number of processes on the system.
// The final entry is the process ID of the process that most recently ran.
//
// TODO: thread save, cache.
//
func LoadAvg2() ([3]float64, [3]int, error) {
	var err error
	if f == nil {
		f, err = os.Open(procfile)

		if err != nil {
			return [3]float64{}, [3]int{},
				fmt.Errorf("loadavg: unable to open procfile %q: %v", procfile, err)
		}
	}

	_, err = f.Seek(0, 0)
	if err != nil {
		return [3]float64{}, [3]int{},
			fmt.Errorf("loadavg: unable to read loadavg: %v", err)
	}

	var loadavg [3]float64
	var pr [3]int

	n, err := fmt.Fscanf(f, "%f %f %f %d/%d %d", &loadavg[0], &loadavg[1], &loadavg[2],
		&pr[0], &pr[1], &pr[2])
	if n != 6 || err != nil {
		return [3]float64{}, [...]int{0, 0, 0},
			fmt.Errorf("loadavg: unable to read loadavg: %v", err)
	}
	return loadavg, pr, nil
}

// LoadAvg returns the traditional 1, 5, and 15 min load averages.
// Not thread save.
func LoadAvg() ([3]float64, error) {
	a, _, c := LoadAvg2()
	return a, c
}

// Donn't bother if you don't have to, the real OS
// will clear all the Go mess at exit anyway.
func Close() error {
	if f != nil {
		f.Close()
		f = nil
	}
	return nil
}
