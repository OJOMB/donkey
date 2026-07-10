package bytecode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewInstruction(t *testing.T) {
	type testCase struct {
		name     string
		op       Opcode
		operands []int
		expected []byte
	}

	testCases := []testCase{
		{"opConstant with 2 byte operand", OpCodeConstant, []int{65534}, []byte{byte(OpCodeConstant), 255, 254}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			instructions := NewInstruction(tc.op, tc.operands...)
			assert.Equal(t, len(tc.expected), len(instructions), "wrong instruction length. expected %d, got %d", len(tc.expected), len(instructions))

			for i, b := range tc.expected {
				assert.Equal(t, b, instructions[i], "wrong byte at pos %d, expected %d, got %d", i, b, instructions[i])
			}
		})
	}
}

func TestInstructionString(t *testing.T) {
	type testCase struct {
		name     string
		input    Instruction
		expected string
	}

	testCases := []testCase{
		{
			name:     "opConstant 2 byte operand",
			input:    NewInstruction(OpCodeConstant, 6553),
			expected: "OpConstant 6553",
		},
		{
			name:     "opAdd no operands",
			input:    NewInstruction(OpCodeAdd),
			expected: "OpAdd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.input.String())
		})
	}
}
