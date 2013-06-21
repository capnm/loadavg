package loadavg

import (
	"fmt"
	"log"
	t "testing"
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

func ExampleLoadAvg() {
	loadavg, err := LoadAvg()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("loadavg: %2.2f, %2.2f, %2.2f", loadavg[MIN_1], loadavg[MIN_5],
		loadavg[MIN_15])
}

func BenchmarkLoadAvg(b *t.B) {
	for i := 0; i < b.N; i++ {
		LoadAvg()
	}
}
