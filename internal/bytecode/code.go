package bytecode

import "encoding/binary"

type Instructions []byte

type Opcode byte

// New creates a new instruction with the given opcode and operands.
// It looks up the definition of the opcode to determine the expected operand widths, and then constructs a byte slice representing the instruction.
// The first byte is the opcode, followed by the operands encoded in big-endian format according to their specified widths.
// If the opcode is not found in the definitions, it returns nil.
func New(op Opcode, operands ...int) Instructions {
	def, ok := Lookup(op)
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

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

func Lookup(op Opcode) (*Definition, bool) {
	def, ok := definitions[op]
	return def, ok
}
