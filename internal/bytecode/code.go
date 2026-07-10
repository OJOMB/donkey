package bytecode

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const (
	OpNameConstant = "opConstant"
)

const (
	OpCodeConstant Opcode = iota + 1
)

var definitions = map[Opcode]*Definition{
	OpCodeConstant: {OpNameConstant, byte(OpCodeConstant), []int{2}},
}

type Opcode byte

type Instruction []byte

func (i Instruction) String() string {
	if len(i) == 0 {
		return ""
	}

	def, ok := lookup(Opcode(i[0]))
	if !ok {
		return "Unknown opcode"
	}

	operands, _ := ReadOperands(def, i[1:])
	if len(operands) != len(def.OperandWidths) {
		return "Wrong number of operands"
	}

	var out strings.Builder
	out.WriteString(def.Name)
	for _, operand := range operands {
		fmt.Fprintf(&out, " %d", operand)
	}

	return out.String()
}

func ReadOperands(def *Definition, ins Instruction) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(binary.BigEndian.Uint16(ins[offset:]))
		}
		offset += width
	}

	return operands, offset
}

// NewInstruction creates a new instruction with the given opcode and operands.
// It looks up the definition of the opcode to determine the expected operand widths, and then constructs a byte slice representing the instruction.
// The first byte is the opcode, followed by the operands encoded in big-endian format according to their specified widths.
// If the opcode is not found in the definitions, it returns nil.
func NewInstruction(op Opcode, operands ...int) Instruction {
	def, ok := lookup(op)
	if !ok {
		return nil
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	instructions := make([]byte, instructionLen)
	instructions[0] = def.Code
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

func lookup(opCode Opcode) (*Definition, bool) {
	def, ok := definitions[opCode]
	return def, ok
}

type Definition struct {
	// Name is the human-readable name of the instruction, used for debugging and disassembly purposes.
	Name string
	// Code is the byte value representing the opcode for this instruction.
	Code byte
	// OperandWidths is a slice of integers representing the widths (in bytes) of each operand for this instruction.
	// many opcodes have different numbers of operands, and each operand can have a different width.
	OperandWidths []int
}
