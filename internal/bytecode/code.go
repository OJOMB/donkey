package bytecode

import "encoding/binary"

type InstructionFactory struct {
	opCodeToDefinition map[Opcode]*Definition
}

func NewInstructionFactory(definitions map[Opcode]*Definition) *InstructionFactory {
	if definitions == nil {
		definitions = defaultDefinitions
	}

	return &InstructionFactory{
		opCodeToDefinition: definitions,
	}
}

type Instruction []byte

type Opcode byte

// NewInstruction creates a new instruction with the given opcode and operands.
// It looks up the definition of the opcode to determine the expected operand widths, and then constructs a byte slice representing the instruction.
// The first byte is the opcode, followed by the operands encoded in big-endian format according to their specified widths.
// If the opcode is not found in the definitions, it returns nil.
func (i *InstructionFactory) NewInstruction(op Opcode, operands ...int) Instruction {
	def, ok := i.Lookup(op)
	if !ok {
		return nil
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instructions := make([]byte, instructionLen)
	instructions[0] = byte(op)

	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instructions[offset:], uint16(o))
		}

		offset += width
	}

	return instructions
}

func (i InstructionFactory) Lookup(opCode Opcode) (*Definition, bool) {
	def, ok := i.opCodeToDefinition[opCode]
	return def, ok
}

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name string
	// OperandWidths specifies the widths of the operands for this instruction.
	// Each width is an integer representing the number of bytes used to encode the corresponding operand.
	// For example, a width of 2 means that the operand is encoded using 2 bytes in big-endian format.
	OperandWidths []int
}

var defaultDefinitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}
