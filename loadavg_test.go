package loadavg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	t "testing"
	"time"
)

func TestLoadAvg(t *t.T) {
	for i := 0; i < 3; i++ {
		loadavg, err := LoadAvg()
		if err != nil {
			t.Fatal(err)
		}
		defer close()
		for _, l := range loadavg[:] {
			if l < 0 {
				t.Errorf("expected loadavg >= 0, got %v", l)
			}
		}
		fmt.Printf("LoadAvg:\t%2.2f, %2.2f, %2.2f\n", loadavg[MIN_1],
			loadavg[MIN_5], loadavg[MIN_15])
	}
}

func TestLoadAvgProc(t *t.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test")
	}

	for i := 0; i < 3; i++ {
		loadavg, pr, err := loadAvgProc()
		if err != nil {
			t.Fatal(err)
		}
		defer close()
		for _, l := range loadavg[:] {
			if l < 0 {
				t.Errorf("expected loadavg >= 0, got %v", l)
			}
		}
		fmt.Printf("loadAvgProc:\t%2.2f, %2.2f, %2.2f, %v\n", loadavg[MIN_1],
			loadavg[MIN_5], loadavg[MIN_15], pr)
	}
}

func TestLoadAvgSys(t *t.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test")
	}

	for i := 0; i < 3; i++ {
		loadavg, pr, err := loadAvgSys()
		if err != nil {
			t.Fatal(err)
		}
		defer close() // nop
		for _, l := range loadavg[:] {
			if l < 0 {
				t.Errorf("expected loadavg >= 0, got %v", l)
			}
		}
		fmt.Printf("loadAvgSys:\t%2.2f, %2.2f, %2.2f, %v\n", loadavg[MIN_1],
			loadavg[MIN_5], loadavg[MIN_15], pr)
	}
}

func ExampleLoadAvg() {
	loadavg, err := LoadAvg()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("loadavg: %2.2f, %2.2f, %2.2f", loadavg[MIN_1], loadavg[MIN_5],
		loadavg[MIN_15])
}

// sleep?
func TestLookPath(t *t.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test")
	}

	path, err := exec.LookPath("sleep")
	if err == nil {
		fmt.Println("LookPath: sleep:", path)
	} else {
		t.Fatal("LookPath:", err)
	}
}

// stress test
func TestProcFileRead(t *t.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test")
	}

	f, err := os.Open(procfile)
	if err != nil {
		t.Fatal("Open:", err)
	}
	defer f.Close()
	for i := 0; i < 10; i++ {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			t.Fatal("ReadAll:", err)
		}
		fmt.Print(string(b))

		for i := 0; i < 100; i++ {
			// permanently ++ the go thread pool
			go exec.Command("sleep", "8").Run()
		}
		runtime.GC()
		time.Sleep(time.Second)
		_, err = f.Seek(0, 0)
		if err != nil {
			t.Fatal("Seek:", err)
		}
	}

	fmt.Println("TestProcFileRead: wait...")
	time.Sleep(16 * time.Second)

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Fatal("Seek:", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal("ReadAll:", err)
	}
	fmt.Println(string(b))

}

// stress test
func TestLinuxSyscall(t *t.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test")
	}

	for i := 0; i < 10; i++ {
		loadavg, pr, err := loadAvgSys()
		if err != nil {
			t.Fatal("loadAvgSys:", err)
		}
		fmt.Printf("%2.2f, %2.2f, %2.2f, %d/%d %d\n", loadavg[MIN_1],
			loadavg[MIN_5], loadavg[MIN_15], pr[0], pr[1], pr[2])

		for i := 0; i < 100; i++ {
			// permanently ++ the go thread pool
			go exec.Command("sleep", "8").Run()
		}
		runtime.GC()
		time.Sleep(time.Second)
	}

	fmt.Println("TestLinuxSyscall: wait...")
	time.Sleep(16 * time.Second)

	loadavg, pr, err := loadAvgSys()
	if err != nil {
		t.Fatal("loadAvgSys:", err)
	}
	fmt.Printf("%2.2f, %2.2f, %2.2f, %d/%d %d\n", loadavg[MIN_1],
		loadavg[MIN_5], loadavg[MIN_15], pr[0], pr[1], pr[2])

}

/*

BenchmarkLoadAvg	 	1000000	      2214 ns/op
BenchmarkLoadAvgSyscall	 	1000000	      2165 ns/op
BenchmarkLoadAvgProcFile	  50000	     40244 ns/op

*/

func BenchmarkLoadAvg(b *t.B) {
	for i := 0; i < b.N; i++ {
		LoadAvg()
	}
}

func BenchmarkLoadAvgSyscall(b *t.B) {
	if runtime.GOOS != "linux" {
		b.Skip("skipping linux benchmark")
	}
	for i := 0; i < b.N; i++ {
		loadAvgSys()
	}
}

func BenchmarkLoadAvgProcFile(b *t.B) {
	if runtime.GOOS != "linux" {
		b.Skip("skipping linux benchmark")
	}
	for i := 0; i < b.N; i++ {
		loadAvgProc()
	}
}
