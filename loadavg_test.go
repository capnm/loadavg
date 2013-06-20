package loadavg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestLoadAvg(t *testing.T) {
	for i := 0; i < 3; i++ {
		loadavg, err := LoadAvg()
		if err != nil {
			t.Fatal(err)
		}
		defer Close()
		for _, l := range loadavg[:] {
			if l < 0 {
				t.Errorf("expected loadavg >= 0, got %v", l)
			}
		}
		fmt.Printf("loadavg: %2.2f, %2.2f, %2.2f\n", loadavg[MIN_1],
			loadavg[MIN_5], loadavg[MIN_15])
	}
}

func TestLoadAvg2(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test")
	}

	for i := 0; i < 3; i++ {
		loadavg, pr, err := LoadAvg2()
		if err != nil {
			t.Fatal(err)
		}
		defer Close()
		for _, l := range loadavg[:] {
			if l < 0 {
				t.Errorf("expected loadavg >= 0, got %v", l)
			}
		}
		fmt.Printf("loadavg: %2.2f, %2.2f, %2.2f, %v\n", loadavg[MIN_1],
			loadavg[MIN_5], loadavg[MIN_15], pr)
	}
}

func TestLoadAvg3(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test")
	}

	for i := 0; i < 3; i++ {
		loadavg, pr, err := LoadAvg3()
		if err != nil {
			t.Fatal(err)
		}
		defer Close() // nop
		for _, l := range loadavg[:] {
			if l < 0 {
				t.Errorf("expected loadavg >= 0, got %v", l)
			}
		}
		fmt.Printf("loadavg: %2.2f, %2.2f, %2.2f, %v\n", loadavg[MIN_1],
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

func TestLookPath(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test.")
	}

	path, err := exec.LookPath("sleep")
	if err == nil {
		fmt.Println("LookPath: sleep:", path)
	} else {
		t.Fatal("LookPath:", err)
	}
}

func TestProcFileRead(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test.")
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

	fmt.Println("wait...")
	time.Sleep(16 * time.Second)

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Fatal("Seek:", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal("ReadAll:", err)
	}
	fmt.Print(string(b))

}

func TestLinuxSyscall(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("skipping linux test.")
	}

	for i := 0; i < 10; i++ {
		loadavg, pr, err := LoadAvg3()
		if err != nil {
			t.Fatal("LoadAvg3:", err)
		}
		fmt.Printf("loadavg: %2.2f, %2.2f, %2.2f, %v\n", loadavg[MIN_1],
			loadavg[MIN_5], loadavg[MIN_15], pr)

		for i := 0; i < 100; i++ {
			// permanently ++ the go thread pool
			go exec.Command("sleep", "8").Run()
		}
		runtime.GC()
		time.Sleep(time.Second)
	}

	fmt.Println("wait...")
	time.Sleep(16 * time.Second)

	loadavg, pr, err := LoadAvg3()
	if err != nil {
		t.Fatal("LoadAvg3:", err)
	}
	fmt.Printf("loadavg: %2.2f, %2.2f, %2.2f, %v\n", loadavg[MIN_1],
		loadavg[MIN_5], loadavg[MIN_15], pr)

}
