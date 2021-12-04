package main

const (
	OpcodeConditionalMove uint8 = iota
	OpcodeArrayIndex
	OpcodeArrayAmendment
	OpcodeAddition
	OpcodeMultiplication
	OpcodeDivision
	OpcodeNotAnd
	OpcodeHalt
	OpcodeAllocation
	OpcodeAbandonment
	OpcodeOutput
	OpcodeInput
	OpcodeLoadProgram
	OpcodeOrthography
)

type Instruction struct {
	data             uint32
	specialRegisterA bool
}

func (i *Instruction) Load(data uint32) {
	i.data = data
	i.specialRegisterA = false
}

func (i *Instruction) Opcode() uint8 {
	opcode := uint8(i.data >> 28)
	if opcode == OpcodeOrthography {
		i.specialRegisterA = true
	}
	return opcode
}

func (i *Instruction) RegisterA() uint8 {
	if i.specialRegisterA {
		return uint8((i.data >> 25) & 7)
	}
	return uint8((i.data >> 6) & 7)
}

func (i *Instruction) RegisterB() uint8 {
	return uint8((i.data >> 3) & 7)
}

func (i *Instruction) RegisterC() uint8 {
	return uint8(i.data & 7)
}
