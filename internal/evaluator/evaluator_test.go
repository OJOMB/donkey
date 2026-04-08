package evaluator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/objects"
)

func TestEvaluatorEvalIntegerExpression(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.ExpressionLiteralInteger
		expected int64
	}

	tests := []testCase{
		{name: "zero", input: &ast.ExpressionLiteralInteger{}, expected: 0},
		{name: "positive int", input: &ast.ExpressionLiteralInteger{Value: 5}, expected: 5},
		{name: "negative int", input: &ast.ExpressionLiteralInteger{Value: -5}, expected: -5},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := NewEvaluator(nil)
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
			evaluator := NewEvaluator(nil)
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
			evaluator := NewEvaluator(nil)
			result := evaluator.Eval(tc.input)

			require.IsType(t, &objects.String{}, result)
			stringResult := result.(*objects.String)
			assert.Equal(t, tc.expected, stringResult.Value)

			assert.Equal(t, objects.TypeString, result.Type())
			assert.Equal(t, fmt.Sprintf("%q", tc.expected), result.Inspect())
		})
	}
}
