package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/donkey/internal/lexer"
)

func TestOperatorPrecedenceParsing(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput string
	}

	var testCases = []testCase{
		// {
		// 	name:           "operator precedence parsing 1",
		// 	input:          `-a * b`,
		// 	expectedOutput: `((-a) * b)`,
		// },
		// {
		// 	name:           "operator precedence parsing 2",
		// 	input:          `!-a`,
		// 	expectedOutput: `(!(-a))`,
		// },
		// {
		// 	name:           "operator precedence parsing 3",
		// 	input:          `a + b + c`,
		// 	expectedOutput: `((a + b) + c)`,
		// },
		// {
		// 	name:           "operator precedence parsing 4",
		// 	input:          `a + b - c;`,
		// 	expectedOutput: `((a + b) - c)`,
		// },
		// {
		// 	name:           "operator precedence parsing 5",
		// 	input:          `a * b / c`,
		// 	expectedOutput: `((a * b) / c)`,
		// },
		// {
		// 	name:           "operator precedence parsing 6",
		// 	input:          `a + b * c + d / e - f`,
		// 	expectedOutput: `(((a + (b * c)) + (d / e)) - f)`,
		// },
		// {
		// 	name:           "operator precedence parsing 7",
		// 	input:          `3 + 4; -5 * 5`,
		// 	expectedOutput: "(3 + 4)\n((-5) * 5)",
		// },
		// {
		// 	name:           "operator precedence parsing 8",
		// 	input:          `5 > 4 == 3 < 4`,
		// 	expectedOutput: `((5 > 4) == (3 < 4))`,
		// },
		// {
		// 	name:           "operator precedence parsing 9",
		// 	input:          `5 < 4 != 3 > 4`,
		// 	expectedOutput: `((5 < 4) != (3 > 4))`,
		// },
		// {
		// 	name:           "operator precedence parsing 10",
		// 	input:          `3 + 4 * 5 == 3 * 1 + 4 * 5`,
		// 	expectedOutput: `((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))`,
		// },
		// {
		// 	name:           "operator precedence parsing 11",
		// 	input:          `3 + 4 * 5 == add(add(3 * 1), 4 * 5)`,
		// 	expectedOutput: `((3 + (4 * 5)) == add(add((3 * 1)), (4 * 5)))`,
		// },
		// {
		// 	name:           "operator precedence parsing 12",
		// 	input:          `-a ** b`,
		// 	expectedOutput: `((-a) ** b)`,
		// },
		// {
		// 	name:           "operator precedence parsing 13",
		// 	input:          `a ** b + c`,
		// 	expectedOutput: `((a ** b) + c)`,
		// },
		// {
		// 	name:           "operator precedence parsing 14",
		// 	input:          `a + b ** c + d`,
		// 	expectedOutput: `((a + (b ** c)) + d)`,
		// },
		// {
		// 	name:           "operator precedence parsing 15",
		// 	input:          `a && b || c && d `,
		// 	expectedOutput: `((a && b) || (c && d))`,
		// },
		// {
		// 	name:           "operator precedence parsing 16",
		// 	input:          `a || b && c || d && e`,
		// 	expectedOutput: `((a || (b && c)) || (d && e))`,
		// },
		// {
		// 	name:           "operator precedence parsing 17",
		// 	input:          `a | b & c | d`,
		// 	expectedOutput: `((a | (b & c)) | d)`,
		// },
		// {
		// 	name:           "operator precedence parsing 18",
		// 	input:          `a & b | c & d ^ e`,
		// 	expectedOutput: `((a & b) | ((c & d) ^ e))`,
		// },
		// {
		// 	name:           "operator precedence parsing 19",
		// 	input:          `[1, 2, 3][0] + 4 * 5`,
		// 	expectedOutput: `([1, 2, 3][0] + (4 * 5))`,
		// },
		{
			name:           "operator precedence parsing 20",
			input:          `{"key1": [1, 2, 3]}["key1"][0] + 4 * 5`,
			expectedOutput: `({key1: [1, 2, 3]}[key1][0] + (4 * 5))`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			require.NotNil(t, program)

			prgrmStr := program.String()
			assert.Equal(t, tc.expectedOutput, prgrmStr)
		})
	}
}
