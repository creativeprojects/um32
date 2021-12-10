package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	var trace bool
	flag.BoolVar(&trace, "trace", false, "Display information on each instruction executed")
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

	computer := NewComputer()
	if trace {
		computer.SetTrace(log.Default())
	}
	loaded, err := computer.Load(file, int(stat.Size()/4+1))
	if err != nil {
		log.Fatalf("loading program: %s", err)
	}
	log.Printf("loaded %d instructions (32bits)\n", loaded)
	err = computer.Run()
	if err != nil {
		log.Fatalf("running program: %s", err)
	}
}
