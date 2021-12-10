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
	arrays      [][]uint32
	registers   []uint32
	instruction *Instruction
	trace       Logger
}

func NewComputer() *Computer {
	return &Computer{
		arrays:      make([][]uint32, 1),
		registers:   make([]uint32, 8),
		instruction: &Instruction{},
		trace:       &dummyLogger{},
	}
}

func (c *Computer) Load(program io.Reader, capacity int) (int, error) {
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
			return len(memory), fmt.Errorf("expected to read 4 bytes but read %d", read)
		}
		temp := binary.BigEndian.Uint32(buffer)
		memory = append(memory, temp)
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
			// c.trace.Printf("[PC %6d] if C(%s), A(%s) = B(%s)", pc, c.DebugRegisterC(), c.DebugRegisterA(), c.DebugRegisterB())
			if c.getRegisterC() != 0 {
				c.setRegisterA(c.getRegisterB())
			}

		case OpcodeArrayIndex:
			// c.trace.Printf("A = B[C]")
			c.setRegisterA(c.arrays[c.getRegisterB()][c.getRegisterC()])

		case OpcodeArrayAmendment:
			// c.trace.Printf("A[B] = C")
			c.arrays[c.getRegisterA()][c.getRegisterB()] = c.getRegisterC()

		case OpcodeAddition:
			// c.trace.Printf("[PC %6d] A(%s) = B(%s) + C(%s)", pc, c.DebugRegisterA(), c.DebugRegisterB(), c.DebugRegisterC())
			c.setRegisterA(c.getRegisterB() + c.getRegisterC())

		case OpcodeMultiplication:
			// c.trace.Printf("A = B * C")
			c.setRegisterA(c.getRegisterB() * c.getRegisterC())

		case OpcodeDivision:
			// c.trace.Printf("A = B / C")
			if c.getRegisterC() == 0 {
				return fmt.Errorf("division by zero")
			}
			c.setRegisterA(c.getRegisterB() / c.getRegisterC())

		case OpcodeNotAnd:
			// c.trace.Printf("A = B nand C")
			c.setRegisterA(^(c.getRegisterB() & c.getRegisterC()))

		case OpcodeHalt:
			log.Print("program halted")
			return nil

		case OpcodeAllocation:
			// c.trace.Printf("B = allocate C words")
			newArray := make([]uint32, c.getRegisterC())
			c.arrays = append(c.arrays, newArray)
			c.setRegisterB(uint32(len(c.arrays) - 1))

		case OpcodeAbandonment:
			// c.trace.Printf("deallocate C")
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
			// c.trace.Printf("[PC %6d] A(%s) = load %x", pc, c.DebugRegisterA(), c.instruction.SpecialData())
			c.setRegisterA(c.instruction.SpecialData())

		default:
			return fmt.Errorf("invalid instruction opcode %d at %d", c.instruction.Opcode(), pc)
		}
	}
	log.Print("program finished (no more instruction)")
	return nil
}

func (c *Computer) SetTrace(trace Logger) {
	c.trace = trace
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

func (c *Computer) DebugRegisterA() string {
	return fmt.Sprintf("%d:%x", c.instruction.RegisterA(), c.getRegisterA())
}

func (c *Computer) DebugRegisterB() string {
	return fmt.Sprintf("%d:%x", c.instruction.RegisterB(), c.getRegisterB())
}
func (c *Computer) DebugRegisterC() string {
	return fmt.Sprintf("%d:%x", c.instruction.RegisterC(), c.getRegisterC())
}
