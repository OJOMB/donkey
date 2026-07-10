package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/donkey/internal/compiler"
	"github.com/OJOMB/donkey/internal/lexer"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/internal/parser"
)

type vmTestCase struct {
	input    string
	expected any
}

func runVMTests(t *testing.T, tcs []vmTestCase) {
	t.Helper()

	for _, tc := range tcs {
		l := lexer.New(tc.input, nil)
		p, err := parser.New(l, nil)
		require.NoError(t, err)

		program := p.ParseProgram()

		c := compiler.New()
		err = c.Compile(program)
		require.NoError(t, err)

		vm := New(c.Bytecode())
		err = vm.Run()
		require.NoError(t, err)

		stackElem := vm.StackTop()

		assert.Equal(t, tc.expected, stackElem)
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tcs := []vmTestCase{
		{"1", &objects.Integer{Value: 1}},
		{"2", &objects.Integer{Value: 2}},
		{"1 + 2", &objects.Integer{Value: 2}}, // fix me later
	}

	runVMTests(t, tcs)
}
