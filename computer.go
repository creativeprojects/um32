package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

type Computer struct {
	arrays      [][]uint32
	registers   []uint32
	instruction *Instruction
}

func NewComputer() *Computer {
	return &Computer{
		arrays:      make([][]uint32, 1),
		registers:   make([]uint32, 8),
		instruction: &Instruction{},
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
	programSize := uint32(len(c.arrays[0]))
	for pc < programSize {
		c.instruction.Load(c.arrays[0][pc])

		switch c.instruction.Opcode() {
		case OpcodeConditionalMove:
			if c.getRegisterC() > 0 {
				c.setRegisterA(c.getRegisterB())
			}

		case OpcodeArrayIndex:
			c.setRegisterA(c.arrays[c.getRegisterB()][c.getRegisterC()])

		case OpcodeArrayAmendment:
			c.arrays[c.getRegisterA()][c.getRegisterB()] = c.getRegisterC()

		case OpcodeAddition:
			c.setRegisterA(c.getRegisterB() + c.getRegisterC())

		case OpcodeMultiplication:
			c.setRegisterA(c.getRegisterB() * c.getRegisterC())

		case OpcodeDivision:
			if c.getRegisterC() == 0 {
				return fmt.Errorf("division by zero")
			}
			c.setRegisterA(c.getRegisterB() / c.getRegisterC())

		case OpcodeNotAnd:
			c.setRegisterA(^(c.getRegisterB() & c.getRegisterC()))

		case OpcodeHalt:
			log.Print("program halted")
			return nil

		case OpcodeAllocation:
			newArray := make([]uint32, c.getRegisterC())
			c.arrays = append(c.arrays, newArray)
			c.setRegisterB(uint32(len(c.arrays) - 1))

		case OpcodeAbandonment:
			c.arrays[c.getRegisterC()] = nil

		case OpcodeOutput:
			fmt.Print(string(c.getRegisterC()))

		case OpcodeInput:
			char := make([]byte, 1)
			for {
				n, err := os.Stdin.Read(char)
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Printf("reading input: %s", err)
					break
				}
				if n != 1 {
					break
				}
			}
		case OpcodeLoadProgram:
			log.Print("loadProgram: not implemented yet")

		case OpcodeOrthography:
			c.setRegisterA(c.instruction.SpecialData())

		default:
			return fmt.Errorf("invalid instruction opcode %d", c.instruction.Opcode())
		}
		pc++
	}
	log.Print("program finished")
	return nil
}

func (c *Computer) getRegister(register uint8) uint32 {
	return c.registers[register]
}

func (c *Computer) getRegisterA() uint32 {
	return c.getRegister(c.instruction.RegisterA())
}

func (c *Computer) getRegisterB() uint32 {
	return c.getRegister(c.instruction.RegisterB())
}

func (c *Computer) getRegisterC() uint32 {
	return c.getRegister(c.instruction.RegisterC())
}

func (c *Computer) setRegister(register uint8, value uint32) {
	c.registers[register] = value
}

func (c *Computer) setRegisterA(value uint32) {
	c.setRegister(c.instruction.RegisterA(), value)
}
func (c *Computer) setRegisterB(value uint32) {
	c.setRegister(c.instruction.RegisterB(), value)
}
func (c *Computer) setRegisterC(value uint32) {
	c.setRegister(c.instruction.RegisterC(), value)
}
