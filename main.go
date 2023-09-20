package main

import (
	"flag"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

func main() {
	var cpuprofile string
	var preAlloc uint
	flag.StringVar(&cpuprofile, "cpu-profile", "", "Saves a CPU profile")
	flag.UintVar(&preAlloc, "pre-alloc", 10_000_000, "Pre-allocate memory")
	flag.Parse()

	args := flag.Args()

	if len(args) <= 0 {
		log.Printf("missing program to run")
		return
	}
	filename := args[0]
	stat, err := os.Stat(filename)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(filename)
	if err != nil {
		log.Printf("%s", err)
		return
	}
	defer file.Close()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Printf("could not create CPU profile: %s", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Printf("could not start CPU profile: %s", err)
		}
		defer pprof.StopCPUProfile()
	}

	start := time.Now()
	computer := NewComputer()
	computer.PreAlloc(uint32(preAlloc))
	loaded, err := computer.Load(file, int(stat.Size()/4+1))
	if err != nil {
		log.Fatalf("loading program: %s", err)
	}
	log.Printf("loaded %d instructions (32bits)\n", loaded)
	err = computer.Run()
	if err != nil {
		log.Fatalf("running program: %s", err)
	}
	log.Printf("finished running in %s", time.Since(start))
}
