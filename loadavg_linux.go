package loadavg

import (
	"fmt"
	"os"
	"syscall"
)

const procfile = "/proc/loadavg"

var f *os.File = nil

// loadAvgProc extracts load average from the /proc/loadavg file.
//
// The first three numbers represent the number of active tasks on the system
//  	– processes that are actually running – averaged over the last 1, 5, and 15 minutes.
//  The next entry shows the instantaneous current number of runnable tasks
//	– processes that are currently scheduled to run rather than being blocked in a system call
//	– and the total number of processes on the system.
// The final entry is the process ID of the process that most recently ran.
//
// TODO:?? thread safety, cache.
//
func loadAvgProc() ([3]float64, [3]int, error) {
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
func LoadAvg() ([3]float64, error) {
	a, _, c := loadAvgSys()
	return a, c
}

// Don't bother if you don't have to, the real OS
// will clear all the Go mess at exit anyway.
func close() {
	if f != nil {
		f.Close()
		f = nil
	}
}

// loadAvgSys extracts load average from the SYSINFO(2) syscall.
//
// The first three numbers represent the number of active tasks on the system
//  	– processes that are actually running – averaged over the last 1, 5, and 15 minutes.
//  The next entry shows the instantaneous current number of runnable tasks
//	– processes that are currently scheduled to run rather than being blocked in a system call
//	– and the total number of processes on the system.
// The final entry is the process ID of the process that most recently ran.
//
// TODO:?? int / 0 - runable tasks, 2 - last pid
//
func loadAvgSys() ([3]float64, [3]int, error) {
	/*
	   // http://man7.org/linux/man-pages/man2/sysinfo.2.html
	   type Sysinfo_t struct {
	   	Uptime    int64		// Seconds since boot
	   	Loads     [3]uint64	// 1, 5, and 15 minute load averages
	   	Totalram  uint64	// Total usable main memory size
	   	Freeram   uint64	// Available memory size
	   	Sharedram uint64	// Amount of shared memory
	   	Bufferram uint64	// Memory used by buffers
	   	Totalswap uint64	// Total swap space siz
	   	Freeswap  uint64	// swap space still available
	   	Procs     uint16	// Number of current processes
	   	Pad       uint16
	   	Pad_cgo_0 [4]byte
	   	Totalhigh uint64	// Total high memory size
	   	Freehigh  uint64	// Available high memory size
	   	Unit      uint32	// Memory unit size in bytes
	   	X_f       [0]byte
	   	Pad_cgo_1 [4]byte	// Padding to 64 bytes
	   }
	*/

	si := &syscall.Sysinfo_t{}

	err := syscall.Sysinfo(si)
	if err != nil {
		return [3]float64{}, [3]int{}, err
	}
	scale := 65536.0 // magic
	return [3]float64{
			float64(si.Loads[0]) / scale,
			float64(si.Loads[1]) / scale,
			float64(si.Loads[2]) / scale,
		},
		[3]int{-1, int(si.Procs), os.Getpid()}, // XXX:?? 0 - runable tasks, 2 - last pid
		nil
}
