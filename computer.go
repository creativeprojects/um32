package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
)

type Computer struct {
	arrays map[uint32][]uint32
}

func NewComputer() *Computer {
	return &Computer{
		arrays: make(map[uint32][]uint32, 1),
	}
}

func (c *Computer) Load(program io.Reader) (int, error) {
	memory := make([]uint32, 0)
	buffer := make([]byte, 4)
	for {
		read, err := program.Read(buffer)
		if err == io.EOF {
			c.arrays[0] = memory
			return len(memory), nil
		}
		if err != nil {
			return len(memory), err
		}
		if read != 4 {
			c.arrays[0] = memory
			return len(memory), nil
		}
		temp := binary.BigEndian.Uint32(buffer)
		memory = append(memory, temp)
	}
}

func (c *Computer) Run() error {
	var pc uint32
	instruction := &Instruction{}
	programSize := uint32(len(c.arrays[0]))
	for pc < programSize {
		instruction.Load(c.arrays[0][pc])

		switch instruction.Opcode() {
		case OpcodeConditionalMove:
		case OpcodeArrayIndex:
		case OpcodeArrayAmendment:
		case OpcodeAddition:
		case OpcodeMultiplication:
		case OpcodeDivision:
		case OpcodeNotAnd:
		case OpcodeHalt:
		case OpcodeAllocation:
		case OpcodeAbandonment:
		case OpcodeOutput:
		case OpcodeInput:
		case OpcodeLoadProgram:

		case OpcodeOrthography:
		default:
			return fmt.Errorf("invalid instruction opcode %d", instruction.Opcode())
		}
		pc++
	}
	log.Print("program finished")
	return nil
}
