package main

import (
	"log"
	"os"
	"testing"
)

func TestComputer(t *testing.T) {
	file, err := os.Open("sandmark.umz")
	if err != nil {
		log.Printf("%s", err)
		return
	}
	defer file.Close()

	computer := NewComputer()
	computer.SetTrace(log.Default())
	loaded, err := computer.Load(file)
	if err != nil {
		log.Printf("loading program: %s", err)
		return
	}
	log.Printf("loaded %d instructions (32bits)\n", loaded)
	err = computer.Run()
	if err != nil {
		log.Printf("running program: %s", err)
		return
	}
}
