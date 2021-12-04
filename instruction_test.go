package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstruction(t *testing.T) {
	fixtures := []struct {
		instruction uint32
		opcode      uint8
		registerA   uint8
		registerB   uint8
		registerC   uint8
	}{
		{0x300000c0, 0x03, 3, 0, 0},
		{0xd2000014, 0x0d, 1, 0b10, 0b100}, // special case on register A
		{0xa000004f, 0b1010, 1, 1, 0b111},
		{0x200001a3, 0b10, 0b110, 0b100, 0b11},
	}

	instruction := &Instruction{}
	for _, fixture := range fixtures {
		t.Run(fmt.Sprintf("%x", fixture.instruction), func(t *testing.T) {
			instruction.Load(fixture.instruction)
			assert.Equal(t, fixture.opcode, instruction.Opcode())
			assert.Equal(t, fixture.registerA, instruction.RegisterA())
			assert.Equal(t, fixture.registerB, instruction.RegisterB())
			assert.Equal(t, fixture.registerC, instruction.RegisterC())
		})
	}
}
