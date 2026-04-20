package evaluator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/internal/tokens"
)

func TestEvaluatorEvalIntegerExpression(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.ExpressionLiteralInteger
		expected int
	}

	tests := []testCase{
		{name: "zero", input: &ast.ExpressionLiteralInteger{}, expected: 0},
		{name: "positive int", input: &ast.ExpressionLiteralInteger{Value: 5}, expected: 5},
		{name: "negative int", input: &ast.ExpressionLiteralInteger{Value: -5}, expected: -5},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input)

			require.IsType(t, &objects.Integer{}, result)
			intResult := result.(*objects.Integer)
			assert.Equal(t, tc.expected, intResult.Value)

			assert.Equal(t, objects.TypeInteger, result.Type())
			assert.Equal(t, fmt.Sprintf("%d", tc.expected), result.Inspect())
		})
	}
}

func TestEvaluatorEvalBooleanExpression(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.ExpressionLiteralBoolean
		expected bool
	}

	tests := []testCase{
		{name: "true", input: &ast.ExpressionLiteralBoolean{Value: true}, expected: true},
		{name: "false", input: &ast.ExpressionLiteralBoolean{Value: false}, expected: false},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input)

			require.IsType(t, &objects.Boolean{}, result)
			boolResult := result.(*objects.Boolean)
			assert.Equal(t, tc.expected, boolResult.Value)

			assert.Equal(t, objects.TypeBoolean, result.Type())
			assert.Equal(t, fmt.Sprintf("%t", tc.expected), result.Inspect())
		})
	}
}

func TestEvaluatorEvalStringExpression(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.ExpressionLiteralString
		expected string
	}

	tests := []testCase{
		{name: "empty string", input: &ast.ExpressionLiteralString{Value: ""}, expected: ""},
		{name: "non-empty string", input: &ast.ExpressionLiteralString{Value: "hello"}, expected: "hello"},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input)

			require.IsType(t, &objects.String{}, result)
			stringResult := result.(*objects.String)
			assert.Equal(t, tc.expected, stringResult.Value)

			assert.Equal(t, objects.TypeString, result.Type())
			assert.Equal(t, fmt.Sprintf("%s", tc.expected), result.Inspect())
		})
	}
}

func TestEvaluatorEvalProgram(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		// {name: "empty program", input: &ast.Program{}, expected: ""},
		{
			name: "non-empty program",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token:      tokens.New(tokens.TypeString, "hello"),
						Expression: &ast.ExpressionLiteralString{Value: "hello"}},
				},
			},
			expected: &objects.String{Value: "hello"},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input)

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}
