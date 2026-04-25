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
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

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
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

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
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

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
			name: "basic string expression program",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token:      tokens.New(tokens.TypeString, "hello"),
						Expression: &ast.ExpressionLiteralString{Value: "hello"}},
				},
			},
			expected: &objects.String{Value: "hello"},
		},
		{
			name: "basic integer expression program",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token:      tokens.New(tokens.TypeInt, "5"),
						Expression: &ast.ExpressionLiteralInteger{Value: 5}},
				},
			},
			expected: &objects.Integer{Value: 5},
		},
		{
			name: "basic boolean expression program",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token:      tokens.New(tokens.TypeTrue, "true"),
						Expression: &ast.ExpressionLiteralBoolean{Value: true}},
				},
			},
			expected: &objects.Boolean{Value: true},
		},
		{
			name: "multiple statements program",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token:      tokens.New(tokens.TypeInt, "5"),
						Expression: &ast.ExpressionLiteralInteger{Value: 5}},
					&ast.StatementExpression{
						Token:      tokens.New(tokens.TypeInt, "10"),
						Expression: &ast.ExpressionLiteralInteger{Value: 10}},
				},
			},
			expected: &objects.Integer{Value: 10}, // the result of evaluating a program is the result of evaluating the last statement in the program
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}

func TestEvaluatorEvalPrefixExpressions(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "bang operator on true",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeBang, "!"),
						Expression: &ast.ExpressionPrefix{
							Token: tokens.New(tokens.TypeBang, "!"),
							Right: &ast.ExpressionLiteralBoolean{Value: true},
						},
					},
				},
			},
			expected: False,
		},
		{
			name: "bang operator on false",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeBang, "!"),
						Expression: &ast.ExpressionPrefix{
							Token: tokens.New(tokens.TypeBang, "!"),
							Right: &ast.ExpressionLiteralBoolean{Value: false},
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "minus operator on positive integer",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeMinus, "-"),
						Expression: &ast.ExpressionPrefix{
							Token: tokens.New(tokens.TypeMinus, "-"),
							Right: &ast.ExpressionLiteralInteger{Value: 5},
						},
					},
				},
			},
			expected: &objects.Integer{Value: -5},
		},
		{
			name: "minus operator on negative integer",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeMinus, "-"),
						Expression: &ast.ExpressionPrefix{
							Token: tokens.New(tokens.TypeMinus, "-"),
							Right: &ast.ExpressionLiteralInteger{Value: -6},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 6},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}

func TestEvaluatorEvalExpressionInfixNumerical(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "1 + 2",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypePlus, "+"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypePlus, "+"),
							Left:     &ast.ExpressionLiteralInteger{Value: 1},
							Right:    &ast.ExpressionLiteralInteger{Value: 2},
							Operator: "+",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 3},
		},
		{
			name: "1 - 2",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeMinus, "-"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeMinus, "-"),
							Left:     &ast.ExpressionLiteralInteger{Value: 1},
							Right:    &ast.ExpressionLiteralInteger{Value: 2},
							Operator: "-",
						},
					},
				},
			},
			expected: &objects.Integer{Value: -1},
		},
		{
			name: "3 * 2",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeAsterisk, "*"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeAsterisk, "*"),
							Left:     &ast.ExpressionLiteralInteger{Value: 3},
							Right:    &ast.ExpressionLiteralInteger{Value: 2},
							Operator: "*",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 6},
		},
		{
			name: "4 / 2",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeForwardSlash, "/"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeForwardSlash, "/"),
							Left:     &ast.ExpressionLiteralInteger{Value: 4},
							Right:    &ast.ExpressionLiteralInteger{Value: 2},
							Operator: "/",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 2},
		},
		{
			name: "4 / 2 * 3 + 1 - 5",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeMinus, "-"),
						Expression: &ast.ExpressionInfix{
							Token: tokens.New(tokens.TypeMinus, "-"),
							Left: &ast.ExpressionInfix{
								Token: tokens.New(tokens.TypePlus, "+"),
								Left: &ast.ExpressionInfix{
									Token: tokens.New(tokens.TypeAsterisk, "*"),
									Left: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeForwardSlash, "/"),
										Left:     &ast.ExpressionLiteralInteger{Value: 4},
										Right:    &ast.ExpressionLiteralInteger{Value: 2},
										Operator: "/",
									},
									Right:    &ast.ExpressionLiteralInteger{Value: 3},
									Operator: "*",
								},
								Right:    &ast.ExpressionLiteralInteger{Value: 1},
								Operator: "+",
							},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "-",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 2},
		},
		{
			name: "10 % 3",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypePercent, "%"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypePercent, "%"),
							Left:     &ast.ExpressionLiteralInteger{Value: 10},
							Right:    &ast.ExpressionLiteralInteger{Value: 3},
							Operator: "%",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 1},
		},
		{
			name: "5 > 3",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeGT, ">"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeGT, ">"),
							Left:     &ast.ExpressionLiteralInteger{Value: 5},
							Right:    &ast.ExpressionLiteralInteger{Value: 3},
							Operator: ">",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "3 > 5",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeGT, ">"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeGT, ">"),
							Left:     &ast.ExpressionLiteralInteger{Value: 3},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: ">",
						},
					},
				},
			},
			expected: False,
		},
		{
			name: "5 < 3",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeLT, "<"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLT, "<"),
							Left:     &ast.ExpressionLiteralInteger{Value: 5},
							Right:    &ast.ExpressionLiteralInteger{Value: 3},
							Operator: "<",
						},
					},
				},
			},
			expected: False,
		},
		{
			name: "3 < 5",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeLT, "<"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLT, "<"),
							Left:     &ast.ExpressionLiteralInteger{Value: 3},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "<",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "5 >= 3",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeGTEQ, ">="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeGTEQ, ">="),
							Left:     &ast.ExpressionLiteralInteger{Value: 5},
							Right:    &ast.ExpressionLiteralInteger{Value: 3},
							Operator: ">=",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "3 >= 5",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeGTEQ, ">="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeGTEQ, ">="),
							Left:     &ast.ExpressionLiteralInteger{Value: 3},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: ">=",
						},
					},
				},
			},
			expected: False,
		},
		{
			name: "5 >= 5",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeGTEQ, ">="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLTEQ, "<="),
							Left:     &ast.ExpressionLiteralInteger{Value: 5},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: ">=",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "3 <= 5",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeLTEQ, "<="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLTEQ, "<="),
							Left:     &ast.ExpressionLiteralInteger{Value: 3},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "<=",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "6 <= 5",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeLTEQ, "<="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLTEQ, "<="),
							Left:     &ast.ExpressionLiteralInteger{Value: 6},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "<=",
						},
					},
				},
			},
			expected: False,
		},
		{
			name: "5 <= 5",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeLTEQ, "<="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLTEQ, "<="),
							Left:     &ast.ExpressionLiteralInteger{Value: 5},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "<=",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "2 ^ 3",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeCaret, "^"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeCaret, "^"),
							Left:     &ast.ExpressionLiteralInteger{Value: 2},
							Right:    &ast.ExpressionLiteralInteger{Value: 3},
							Operator: "^",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 8},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}

func TestEvaluatorEvalExpressionInfixBoolean(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "true == true",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeEq, "=="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeEq, "=="),
							Left:     &ast.ExpressionLiteralBoolean{Value: true},
							Right:    &ast.ExpressionLiteralBoolean{Value: true},
							Operator: "==",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "true == false",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeEq, "=="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeEq, "=="),
							Left:     &ast.ExpressionLiteralBoolean{Value: true},
							Right:    &ast.ExpressionLiteralBoolean{Value: false},
							Operator: "==",
						},
					},
				},
			},
			expected: False,
		},
		{
			name: "true != false",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeNotEq, "!="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeNotEq, "!="),
							Left:     &ast.ExpressionLiteralBoolean{Value: true},
							Right:    &ast.ExpressionLiteralBoolean{Value: false},
							Operator: "!=",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "true != true",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeNotEq, "!="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeNotEq, "!="),
							Left:     &ast.ExpressionLiteralBoolean{Value: true},
							Right:    &ast.ExpressionLiteralBoolean{Value: true},
							Operator: "!=",
						},
					},
				},
			},
			expected: False,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}

func TestEvaluatorEvalExpressionInfixString(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "hello + world",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypePlus, "+"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypePlus, "+"),
							Left:     &ast.ExpressionLiteralString{Value: "hello"},
							Right:    &ast.ExpressionLiteralString{Value: "world"},
							Operator: "+",
						},
					},
				},
			},
			expected: &objects.String{Value: "helloworld"},
		},
		{
			name: "foobar - bar",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeMinus, "-"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeMinus, "-"),
							Left:     &ast.ExpressionLiteralString{Value: "foobar"},
							Right:    &ast.ExpressionLiteralString{Value: "bar"},
							Operator: "-",
						},
					},
				},
			},
			expected: &objects.String{Value: "foo"},
		},
		{
			name: "foobar - foo",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeMinus, "-"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeMinus, "-"),
							Left:     &ast.ExpressionLiteralString{Value: "foobar"},
							Right:    &ast.ExpressionLiteralString{Value: "foo"},
							Operator: "-",
						},
					},
				},
			},
			expected: &objects.String{Value: "foobar"},
		},
		{
			name: "foobar == foobar",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeEq, "=="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeEq, "=="),
							Left:     &ast.ExpressionLiteralString{Value: "foobar"},
							Right:    &ast.ExpressionLiteralString{Value: "foobar"},
							Operator: "==",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "foobar == barfoo",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeEq, "=="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeEq, "=="),
							Left:     &ast.ExpressionLiteralString{Value: "foobar"},
							Right:    &ast.ExpressionLiteralString{Value: "barfoo"},
							Operator: "==",
						},
					},
				},
			},
			expected: False,
		},
		{
			name: "foobar != barfoo",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeNotEq, "!="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeNotEq, "!="),
							Left:     &ast.ExpressionLiteralString{Value: "foobar"},
							Right:    &ast.ExpressionLiteralString{Value: "barfoo"},
							Operator: "!=",
						},
					},
				},
			},
			expected: True,
		},
		{
			name: "foobar != foobar",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeNotEq, "!="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeNotEq, "!="),
							Left:     &ast.ExpressionLiteralString{Value: "foobar"},
							Right:    &ast.ExpressionLiteralString{Value: "foobar"},
							Operator: "!=",
						},
					},
				},
			},
			expected: False,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())

		})
	}
}

func TestEvaluatorEvalConditionals(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "if condition true, elif is also true but should not be evaluated, neither should else",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIf, "if"),
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 1},
										Right:    &ast.ExpressionLiteralInteger{Value: 1},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "if condition was true"},
											},
										},
									},
								},
								{
									Token: tokens.Token{Type: tokens.TypeElif, Lexeme: "elif"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 2},
										Right:    &ast.ExpressionLiteralInteger{Value: 2},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "elif condition 1 was true"},
											},
										},
									},
								},
							},
							Alternative: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementExpression{
										Token:      tokens.New(tokens.TypeReturn, "return"),
										Expression: &ast.ExpressionLiteralString{Value: "else evaluated"},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.String{Value: "if condition was true"},
		},
		{
			name: "if condition false, 1st elif is true, else should not be evaluated",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIf, "if"),
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 1},
										Right:    &ast.ExpressionLiteralInteger{Value: 42},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "if condition was true"},
											},
										},
									},
								},
								{
									Token: tokens.Token{Type: tokens.TypeElif, Lexeme: "elif"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 2},
										Right:    &ast.ExpressionLiteralInteger{Value: 2},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "elif condition 1 was true"},
											},
										},
									},
								},
							},
							Alternative: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementExpression{
										Token:      tokens.New(tokens.TypeReturn, "return"),
										Expression: &ast.ExpressionLiteralString{Value: "else evaluated"},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.String{Value: "elif condition 1 was true"},
		},
		{
			name: "if condition false, 2nd elif is true, else should not be evaluated",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIf, "if"),
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 1},
										Right:    &ast.ExpressionLiteralInteger{Value: 42},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "if condition was true"},
											},
										},
									},
								},
								{
									Token: tokens.Token{Type: tokens.TypeElif, Lexeme: "elif"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralString{Value: "two"},
										Right:    &ast.ExpressionLiteralString{Value: "three"},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "elif condition 1 was true"},
											},
										},
									},
								},
								{
									Token: tokens.Token{Type: tokens.TypeElif, Lexeme: "elif"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 3},
										Right:    &ast.ExpressionLiteralInteger{Value: 3},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "elif condition 2 was true"},
											},
										},
									},
								},
							},
							Alternative: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementExpression{
										Token:      tokens.New(tokens.TypeReturn, "return"),
										Expression: &ast.ExpressionLiteralString{Value: "else evaluated"},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.String{Value: "elif condition 2 was true"},
		},
		{
			name: "if condition false, elif is false, else should be evaluated",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIf, "if"),
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 1},
										Right:    &ast.ExpressionLiteralInteger{Value: 42},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "if condition was true"},
											},
										},
									},
								},
								{
									Token: tokens.Token{Type: tokens.TypeElif, Lexeme: "elif"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 2},
										Right:    &ast.ExpressionLiteralInteger{Value: 3},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "elif condition 1 was true"},
											},
										},
									},
								},
							},
							Alternative: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementExpression{
										Token:      tokens.New(tokens.TypeReturn, "return"),
										Expression: &ast.ExpressionLiteralString{Value: "else evaluated"},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.String{Value: "else evaluated"},
		},
		{
			name: "if condition false, elif is false, no else block, should return Nowt",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIf, "if"),
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 1},
										Right:    &ast.ExpressionLiteralInteger{Value: 42},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "if condition was true"},
											},
										},
									},
								},
								{
									Token: tokens.Token{Type: tokens.TypeElif, Lexeme: "elif"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeEq, "=="),
										Left:     &ast.ExpressionLiteralInteger{Value: 2},
										Right:    &ast.ExpressionLiteralInteger{Value: 3},
										Operator: "==",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token:      tokens.New(tokens.TypeReturn, "return"),
												Expression: &ast.ExpressionLiteralString{Value: "elif condition 1 was true"},
											},
										},
									},
								},
							},
							Alternative: nil,
						},
					},
				},
			},
			expected: Nowt,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}

func TestEvaluatorEvalReturnStatements(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "return 5;",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementReturn{
						Token: tokens.New(tokens.TypeReturn, "return"),
						Value: &ast.ExpressionLiteralInteger{Value: 5},
					},
				},
			},
			expected: &objects.Integer{Value: 5},
		},
		{
			name: "do not eval after return statement",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementReturn{
						Token: tokens.New(tokens.TypeReturn, "return"),
						Value: &ast.ExpressionLiteralInteger{Value: 5},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypePlus, "+"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypePlus, "+"),
							Left:     &ast.ExpressionLiteralInteger{Value: 10},
							Right:    &ast.ExpressionLiteralInteger{Value: 20},
							Operator: "+",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 5},
		},
		{
			name: "return in if consequence",
			// if (10 > 1) {
			//  if (10 > 1) {
			//    return 10
			//  }
			//  return 1
			// } else {
			//  return -1
			// }
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIf, "if"),
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeGT, ">"),
										Left:     &ast.ExpressionLiteralInteger{Value: 10},
										Right:    &ast.ExpressionLiteralInteger{Value: 1},
										Operator: ">",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token: tokens.New(tokens.TypeIf, "if"),
												Expression: &ast.ExpressionIf{
													Branches: []ast.ConditionalBranch{
														{
															Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
															Condition: &ast.ExpressionInfix{
																Token:    tokens.New(tokens.TypeGT, ">"),
																Left:     &ast.ExpressionLiteralInteger{Value: 10},
																Right:    &ast.ExpressionLiteralInteger{Value: 1},
																Operator: ">",
															},
															Consequence: &ast.StatementBlock{
																Statements: []ast.Statement{
																	&ast.StatementReturn{
																		Token: tokens.New(tokens.TypeReturn, "return"),
																		Value: &ast.ExpressionLiteralInteger{Value: 10},
																	},
																},
															},
														},
													},
													Alternative: &ast.StatementBlock{
														Statements: []ast.Statement{
															&ast.StatementReturn{
																Token: tokens.New(tokens.TypeReturn, "return"),
																Value: &ast.ExpressionLiteralInteger{Value: 1},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							Alternative: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementReturn{
										Token: tokens.New(tokens.TypeReturn, "return"),
										Value: &ast.ExpressionLiteralInteger{Value: -1},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 10},
		},
		{
			name: "return in if consequence",
			// if (1 > 10) {
			//  if (10 > 1) {
			//    return 10
			//  }
			//  return 1
			// } else {
			//  return -1
			// }
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIf, "if"),
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeGT, ">"),
										Left:     &ast.ExpressionLiteralInteger{Value: 1},
										Right:    &ast.ExpressionLiteralInteger{Value: 10},
										Operator: ">",
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token: tokens.New(tokens.TypeIf, "if"),
												Expression: &ast.ExpressionIf{
													Branches: []ast.ConditionalBranch{
														{
															Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
															Condition: &ast.ExpressionInfix{
																Token:    tokens.New(tokens.TypeGT, ">"),
																Left:     &ast.ExpressionLiteralInteger{Value: 10},
																Right:    &ast.ExpressionLiteralInteger{Value: 1},
																Operator: ">",
															},
															Consequence: &ast.StatementBlock{
																Statements: []ast.Statement{
																	&ast.StatementReturn{
																		Token: tokens.New(tokens.TypeReturn, "return"),
																		Value: &ast.ExpressionLiteralInteger{Value: 10},
																	},
																},
															},
														},
													},
													Alternative: &ast.StatementBlock{
														Statements: []ast.Statement{
															&ast.StatementReturn{
																Token: tokens.New(tokens.TypeReturn, "return"),
																Value: &ast.ExpressionLiteralInteger{Value: 1},
															},
														},
													},
												},
											},
										},
									},
								},
							},
							Alternative: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementReturn{
										Token: tokens.New(tokens.TypeReturn, "return"),
										Value: &ast.ExpressionLiteralInteger{Value: -1},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.Integer{Value: -1},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}

func TestEvaluatorEvalErrorHandling(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "10 / 0",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeForwardSlash, "/"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeForwardSlash, "/"),
							Left:     &ast.ExpressionLiteralInteger{Value: 10},
							Right:    &ast.ExpressionLiteralInteger{Value: 0},
							Operator: "/",
						},
					},
				},
			},
			expected: &objects.ErrorValue{Message: "division by zero"},
		},
		{
			name: "10 % 0",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypePercent, "%"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypePercent, "%"),
							Left:     &ast.ExpressionLiteralInteger{Value: 10},
							Right:    &ast.ExpressionLiteralInteger{Value: 0},
							Operator: "%",
						},
					},
				},
			},
			expected: &objects.ErrorValue{Message: "modulo by zero"},
		},
		{
			name: "type mismatch: 5 + true",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypePlus, "+"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypePlus, "+"),
							Left:     &ast.ExpressionLiteralInteger{Value: 5},
							Right:    &ast.ExpressionLiteralBoolean{Value: true},
							Operator: "+",
						},
					},
				},
			},
			expected: &objects.ErrorValue{Message: "type mismatch for infix operator: INTEGER + BOOLEAN"},
		},
		{
			name: "ident not found",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "a"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
					},
				},
			},
			expected: &objects.ErrorValue{Message: "identifier not found: a"},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}

func TestEvaluatorEvalBindStatements(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "let a = 5; a;",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "let"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 5},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "a"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 5},
		},
		{
			name: "let a = 5; let b = a; b;",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "let"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 5},
					},
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "let"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "b"),
							Value: "b",
						},
						Value: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "b"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "b"),
							Value: "b",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 5},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}

func TestEvaluatorEvalRebindStatements(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: `
				let a = 5;
				a = 13;
				a;`,
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "let"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 5},
					},
					&ast.StatementRebind{
						Token: tokens.New(tokens.TypeIdent, "a"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 13},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "a"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 13},
		},
		{
			name: `
			let a = 5;
			let b = a;
			b = 109;
			b;`,
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "let"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 5},
					},
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "let"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "b"),
							Value: "b",
						},
						Value: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
					},
					&ast.StatementRebind{
						Token: tokens.New(tokens.TypeIdent, "b"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "b"),
							Value: "b",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 109},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "b"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "b"),
							Value: "b",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 109},
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}
