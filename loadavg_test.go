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

func TestLoadAvg2(t *testing.T) {
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

func ExampleLoadAvg() {
	loadavg, err := LoadAvg()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("loadavg: %2.2f, %2.2f, %2.2f", loadavg[MIN_1], loadavg[MIN_5],
		loadavg[MIN_15])
}

func TestLookPath(t *testing.T) {
	path, err := exec.LookPath("sleep")
	if err == nil {
		fmt.Println("LookPath: sleep:", path)
	} else {
		fmt.Println("LookPath:", err)
	}
}


func TestProcFileRead(t *testing.T) {
	f, err := os.Open(procfile)
	if err != nil {
		t.Error("Open:", err)
	}
	defer f.Close()
	for i := 0; i < 10; i++ {
		b, err := ioutil.ReadAll(f)
		if err != nil {
			t.Error("ReadAll:", err)
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
			t.Error("Seek:", err)
		}
	}

	fmt.Println("wait...")
	time.Sleep(16 * time.Second)

	_, err = f.Seek(0, 0)
	if err != nil {
		t.Error("Seek:", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error("ReadAll:", err)
	}
	fmt.Print(string(b))

}
