package bytecode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type testCase struct {
		op       Opcode
		operands []int
		expected []byte
	}

	tests := []testCase{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
	}

	for _, tt := range tests {
		instructions := New(tt.op, tt.operands...)
		assert.Equal(t, len(tt.expected), len(instructions), "wrong instruction length. expected %d, got %d", len(tt.expected), len(instructions))

		for i, b := range tt.expected {
			assert.Equal(t, b, instructions[i], "wrong byte at pos %d, expected %d, got %d", i, b, instructions[i])
		}
	}
}
