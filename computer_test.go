package main

import (
	"log"
	"os"
	"runtime/pprof"
	"testing"
	"time"
)

func TestComputer(t *testing.T) {
	programs := []struct {
		filename   string
		preAlloc   uint32
		cpuprofile string
	}{
		{"midmark.um", 10_000_000, "midmark.prof"},
		{"sandmark.umz", 400_000_000, "sandmark.prof"},
	}
	for _, program := range programs {
		t.Run(program.filename, func(t *testing.T) {
			stat, err := os.Stat(program.filename)
			if err != nil {
				t.Skipf("could not stat program: %s", err)
			}
			file, err := os.Open(program.filename)
			if err != nil {
				t.Skipf("could not open program: %s", err)
			}
			defer file.Close()

			if program.cpuprofile != "" {
				f, err := os.Create(program.cpuprofile)
				if err != nil {
					t.Skipf("could not create CPU profile: %s", err)
				}
				defer f.Close()
				if err := pprof.StartCPUProfile(f); err != nil {
					t.Skipf("could not start CPU profile: %s", err)
				}
				defer pprof.StopCPUProfile()
			}

			computer := NewComputer()
			computer.PreAlloc(400_000_000)
			start := time.Now()
			loaded, err := computer.Load(file, int(stat.Size()/4+1))
			if err != nil {
				log.Printf("loading program: %s", err)
				return
			}
			log.Printf("loaded %d instructions (32bits) in %s\n", loaded, time.Since(start))
			start = time.Now()
			err = computer.Run()
			if err != nil {
				log.Printf("running program: %s", err)
				return
			}
			log.Printf("finished running in %s", time.Since(start))
			log.Printf("buffer allocations: %d, total size = %d, biggest buffer size = %d\n",
				computer.bufferCount,
				computer.totalAllocSize,
				computer.maxAllocSize)
		})
	}
}

func TestLoadBigProgram(t *testing.T) {
	filename := "output.um"
	stat, err := os.Stat(filename)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	defer file.Close()

	f, err := os.Create("load.prof")
	if err != nil {
		t.Errorf("could not create CPU profile: %s", err)
	}
	defer f.Close()
	if err := pprof.StartCPUProfile(f); err != nil {
		t.Errorf("could not start CPU profile: %s", err)
	}
	defer pprof.StopCPUProfile()

	computer := NewComputer()
	loaded, err := computer.Load(file, int(stat.Size()/4+1))
	if err != nil {
		log.Printf("loading program: %s", err)
		return
	}
	log.Printf("loaded %d instructions (32bits)\n", loaded)
}
