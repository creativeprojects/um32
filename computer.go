package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"math"
	"os"
)

type Computer struct {
	registers      [8]uint32
	arrays         [][]uint32
	preAlloc       []uint32
	preAllocIndex  uint32
	instruction    *Instruction
	bufferCount    uint32
	totalAllocSize uint32
	maxAllocSize   uint32
}

func NewComputer() *Computer {
	return &Computer{
		arrays:      make([][]uint32, 1, 100_000_000),
		instruction: &Instruction{},
	}
}

func (c *Computer) PreAlloc(size uint32) {
	c.preAlloc = make([]uint32, size)
	c.preAllocIndex = 0
}

func (c *Computer) Load(program io.Reader, capacity int) (int, error) {
	memory := make([]uint32, capacity)
	buffer := make([]byte, 4)
	index := 0
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
			return len(memory), fmt.Errorf("expected to read 4 bytes but read %d", read)
		}
		temp := binary.BigEndian.Uint32(buffer)
		if index >= len(memory) {
			memory = append(memory, temp)
		} else {
			memory[index] = temp
		}
		index++
	}
}

func (c *Computer) Run() error {
	var pc uint32
	for pc < uint32(len(c.arrays[0])) {
		c.instruction.Load(c.arrays[0][pc])
		pc++

		opcode := c.instruction.Opcode()
		switch opcode {
		case OpcodeConditionalMove:
			if c.getRegisterC() != 0 {
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
			size := c.getRegisterC()
			if size > c.maxAllocSize {
				c.maxAllocSize = size
			}
			if c.preAllocIndex+size <= uint32(len(c.preAlloc)) {
				c.arrays = append(c.arrays, c.preAlloc[c.preAllocIndex:c.preAllocIndex+size])
				c.preAllocIndex += size
			} else {
				newArray := make([]uint32, size)
				c.arrays = append(c.arrays, newArray)
			}
			c.bufferCount++
			c.totalAllocSize += size
			c.setRegisterB(uint32(len(c.arrays) - 1))

		case OpcodeAbandonment:
			c.arrays[c.getRegisterC()] = nil

		case OpcodeOutput:
			fmt.Print(string(byte(c.getRegisterC())))

		case OpcodeInput:
			char := make([]byte, 1)
			n, err := os.Stdin.Read(char)
			if err == io.EOF {
				c.setRegisterC(math.MaxUint32)
				break
			}
			if err != nil {
				return fmt.Errorf("reading input: %s", err)
			}
			if n != 1 {
				c.setRegisterC(math.MaxUint32)
				break
			}
			if char[0] > 255 {
				log.Printf("invalid input: %c", char[0])
			}
			c.setRegisterC(uint32(char[0]))

		case OpcodeLoadProgram:
			if c.getRegisterB() != 0 {
				// it needs a *copy* of the array
				array := c.arrays[c.getRegisterB()]
				temp := make([]uint32, len(array))
				copy(temp, array)
				c.arrays[0] = temp
			}
			pc = c.getRegisterC()

		case OpcodeOrthography:
			c.setRegisterA(c.instruction.SpecialData())

		default:
			return fmt.Errorf("invalid instruction opcode %d at %d", c.instruction.Opcode(), pc)
		}
	}
	log.Print("program finished (no more instruction)")
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
