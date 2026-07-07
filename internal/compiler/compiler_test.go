package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/donkey/internal/bytecode"
	"github.com/OJOMB/donkey/internal/lexer"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/internal/parser"
)

func TestCompilerCompile(t *testing.T) {
	type testCase struct {
		name                 string
		input                string
		expectedConstants    []objects.Object
		expectedInstructions []bytecode.Instruction
	}

	testCases := []testCase{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := lexer.New(tc.input, nil)
			p, err := parser.New(l, nil)
			require.NoError(t, err)

			program := p.ParseProgram()

			c := New()
			err = c.Compile(program)
			require.NoError(t, err)

			bc := c.Bytecode()

			assert.Equal(t, tc.expectedConstants, bc.Constants)
			assert.Equal(t, tc.expectedInstructions, bc.Instructions)
		})
	}
}
