package evaluator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/lexer"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/internal/parser"
	"github.com/OJOMB/donkey/internal/tokens"
)

func TestLexParseEval(t *testing.T) {
	type testCase struct {
		name     string
		input    string
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "statementBind function definition and call",
			input: `
			var addOne = fn(x) { return x + 1; };
			addOne(5);`,
			expected: &objects.Integer{Value: 6},
		},
		{
			name: "bitwise operators",
			input: `
			var a = 5;
			var b = 3;
			a = a & b;
			a = a | b;
			a;`,
			expected: &objects.Integer{Value: 3},
		},
		{
			name: "list ints",
			input: `
				var list = [1, 2, 3];
				list;`,
			expected: &objects.List{Elements: []objects.Object{
				&objects.Integer{Value: 1},
				&objects.Integer{Value: 2},
				&objects.Integer{Value: 3},
			}},
		},
		{
			name: "list strings",
			input: `
				var list = ["foo", "bar", "baz"];
				list;`,
			expected: &objects.List{Elements: []objects.Object{
				&objects.String{Value: "foo"},
				&objects.String{Value: "bar"},
				&objects.String{Value: "baz"},
			}},
		},
		{
			name: "list indexed",
			input: `
				var list = [1, 2, 3];
				list[2];`,
			expected: &objects.Integer{Value: 3},
		},
		{
			name: "statement block that looks like a map (i.e. does not start with an explicit statement token)",
			input: `
				var foo = 0;

				{
					foo = 10;
				}

				foo;`,
			expected: &objects.Integer{Value: 10},
		},
		{
			name: "check that variables declared in a block do not leak out into the surrounding code",
			input: `
				{
					var foo = 10;
				}
				foo;`,
			expected: &objects.ErrorValue{Message: "identifier not found: foo"},
		},
	}

	evaluator := New(nil)

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			l := lexer.New(tc.input, nil)
			p, err := parser.New(l, nil)
			require.NoError(t, err)
			program := p.ParseProgram()

			result := evaluator.Eval(program, objects.NewEnvironment())

			assert.Equal(t, tc.expected.Type(), result.Type())
			assert.Equal(t, tc.expected.Inspect(), result.Inspect())
		})
	}
}

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
						Token: tokens.New(tokens.TypeExponent, "^"),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeExponent, "^"),
							Left:     &ast.ExpressionLiteralInteger{Value: 2},
							Right:    &ast.ExpressionLiteralInteger{Value: 3},
							Operator: "^",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 1},
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
						Token: tokens.New(tokens.TypeEQ, "=="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeEQ, "=="),
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
						Token: tokens.New(tokens.TypeEQ, "=="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeEQ, "=="),
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
						Token: tokens.New(tokens.TypeNotEQ, "!="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeNotEQ, "!="),
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
						Token: tokens.New(tokens.TypeNotEQ, "!="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeNotEQ, "!="),
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
						Token: tokens.New(tokens.TypeEQ, "=="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeEQ, "=="),
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
						Token: tokens.New(tokens.TypeEQ, "=="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeEQ, "=="),
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
						Token: tokens.New(tokens.TypeNotEQ, "!="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeNotEQ, "!="),
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
						Token: tokens.New(tokens.TypeNotEQ, "!="),
						Expression: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeNotEQ, "!="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
										Token:    tokens.New(tokens.TypeEQ, "=="),
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
			name: "var a = 5; a;",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
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
			name: "var a = 5; var b = a; b;",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 5},
					},
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
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
				var a = 5;
				a = 13;
				a;`,
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
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
			var a = 5;
			var b = a;
			b = 109;
			b;`,
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "a"),
							Value: "a",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 5},
					},
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
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

func TestEvaluatorEvalWhileLoops(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: `
				var i = 0;
				while (i < 5) {
					i = i + 1;
				}

				i;`,
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "i"),
							Value: "i",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 0},
					},
					&ast.StatementWhile{
						Token: tokens.New(tokens.TypeWhile, "while"),
						Condition: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLT, "<"),
							Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "<",
						},
						Body: &ast.StatementBlock{
							Statements: []ast.Statement{
								&ast.StatementRebind{
									Token: tokens.New(tokens.TypeIdent, "i"),
									Name: &ast.ExpressionIdentifier{
										Token: tokens.New(tokens.TypeIdent, "i"),
										Value: "i",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypePlus, "+"),
										Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
										Right:    &ast.ExpressionLiteralInteger{Value: 1},
										Operator: "+",
									},
								},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "i"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "i"),
							Value: "i",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 5},
		},
		// TODO: re-enable when we have break statements implemented
		// {
		// 	name: `while loop with no conditional
		// 		var i = 0;
		// 		while {
		// 			i = i + 1;
		// 			if (i == 5) {
		// 				break;
		// 			}
		// 		}

		// 		i;
		// 	`,
		// 	input: &ast.Program{
		// 		Statements: []ast.Statement{
		// 			&ast.StatementBind{
		// 				Token: tokens.New(tokens.TypeBind, "var"),
		// 				Name: &ast.ExpressionIdentifier{
		// 					Token: tokens.New(tokens.TypeIdent, "i"),
		// 					Value: "i",
		// 				},
		// 				Value: &ast.ExpressionLiteralInteger{Value: 0},
		// 			},
		// 			&ast.StatementWhile{
		// 				Token:     tokens.New(tokens.TypeWhile, "while"),
		// 				Condition: nil,
		// 				Body: &ast.StatementBlock{
		// 					Statements: []ast.Statement{
		// 						&ast.StatementRebind{
		// 							Token: tokens.New(tokens.TypeIdent, "i"),
		// 							Name: &ast.ExpressionIdentifier{
		// 								Token: tokens.New(tokens.TypeIdent, "i"),
		// 								Value: "i",
		// 							},
		// 							Value: &ast.ExpressionInfix{
		// 								Token:    tokens.New(tokens.TypePlus, "+"),
		// 								Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
		// 								Right:    &ast.ExpressionLiteralInteger{Value: 1},
		// 								Operator: "+",
		// 							},
		// 						},
		// 						&ast.StatementExpression{
		// 							Token: tokens.New(tokens.TypeIf, "if"),
		// 							Expression: &ast.ExpressionIf{
		// 								Branches: []ast.ConditionalBranch{
		// 									{
		// 										Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
		// 										Condition: &ast.ExpressionInfix{
		// 											Token:    tokens.New(tokens.TypeEq, "=="),
		// 											Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
		// 											Right:    &ast.ExpressionLiteralInteger{Value: 5},
		// 											Operator: "==",
		// 										},
		// 										Consequence: &ast.StatementBlock{
		// 											Statements: []ast.Statement{
		// 												&ast.StatementExpression{
		// 													Token: tokens.New(tokens.TypeBreak, "break"),
		// 													Expression: &ast.ExpressionKeyword{
		// 														Token:   tokens.New(tokens.TypeBreak, "break"),
		// 														Keyword: "break",
		// 													},
		// 												},
		// 											},
		// 										},
		// 									},
		// 								},
		// 								Alternative: nil,
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 			&ast.StatementExpression{
		// 				Token: tokens.New(tokens.TypeIdent, "i"),
		// 				Expression: &ast.ExpressionIdentifier{
		// 					Token: tokens.New(tokens.TypeIdent, "i"),
		// 					Value: "i",
		// 				},
		// 			},
		// 		},
		// 	},
		// 	expected: &objects.Integer{Value: 5},
		// },
		{
			name: `double while loop
				var i = 0;
				while (i < 5) {
					while (i < 5) {
						i = i + 1;
					}
				}

				i;
			`,
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "i"),
							Value: "i",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 0},
					},
					&ast.StatementWhile{
						Token: tokens.New(tokens.TypeWhile, "while"),
						Condition: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLT, "<"),
							Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "<",
						},
						Body: &ast.StatementBlock{
							Statements: []ast.Statement{
								&ast.StatementWhile{
									Token: tokens.New(tokens.TypeWhile, "while"),
									Condition: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypeLT, "<"),
										Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
										Right:    &ast.ExpressionLiteralInteger{Value: 5},
										Operator: "<",
									},
									Body: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementRebind{
												Token: tokens.New(tokens.TypeIdent, "i"),
												Name: &ast.ExpressionIdentifier{
													Token: tokens.New(tokens.TypeIdent, "i"),
													Value: "i",
												},
												Value: &ast.ExpressionInfix{
													Token:    tokens.New(tokens.TypePlus, "+"),
													Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
													Right:    &ast.ExpressionLiteralInteger{Value: 1},
													Operator: "+",
												},
											},
										},
									},
								},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "i"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "i"),
							Value: "i",
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

func TestEvaluastorEvalForLoops(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: `
				var result = 0;
				for (var i = 0; i < 5; i = i + 1) {
					result = result + i;
				}

				result;`,
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "result"),
							Value: "result",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 0},
					},
					&ast.StatementFor{
						Token: tokens.New(tokens.TypeFor, "for"),
						Initializer: &ast.StatementBind{
							Token: tokens.New(tokens.TypeBind, "var"),
							Name: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "i"),
								Value: "i",
							},
							Value: &ast.ExpressionLiteralInteger{Value: 0},
						},
						Condition: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLT, "<"),
							Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "<",
						},
						Step: &ast.StatementRebind{
							Token: tokens.New(tokens.TypeIdent, "i"),
							Name: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "i"),
								Value: "i",
							},
							Value: &ast.ExpressionInfix{
								Token:    tokens.New(tokens.TypePlus, "+"),
								Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
								Right:    &ast.ExpressionLiteralInteger{Value: 1},
								Operator: "+",
							},
						},
						Body: &ast.StatementBlock{
							Statements: []ast.Statement{
								&ast.StatementRebind{
									Token: tokens.New(tokens.TypeIdent, "result"),
									Name: &ast.ExpressionIdentifier{
										Token: tokens.New(tokens.TypeIdent, "result"),
										Value: "result",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypePlus, "+"),
										Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "result"), Value: "result"},
										Right:    &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
										Operator: "+",
									},
								},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "result"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "result"),
							Value: "result",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 10},
		},
		{
			name: `
				var result = 0;
				for (var i = 0; i < 6; i = i + 1) {
					if (i % 2 == 0) {
						continue;
					}

					result = result + i;
				}

				result;`,
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "result"),
							Value: "result",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 0},
					},
					&ast.StatementFor{
						Token: tokens.New(tokens.TypeFor, "for"),
						Initializer: &ast.StatementBind{
							Token: tokens.New(tokens.TypeBind, "var"),
							Name: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "i"),
								Value: "i",
							},
							Value: &ast.ExpressionLiteralInteger{Value: 0},
						},
						Condition: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLT, "<"),
							Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "<",
						},
						Step: &ast.StatementRebind{
							Token: tokens.New(tokens.TypeIdent, "i"),
							Name: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "i"),
								Value: "i",
							},
							Value: &ast.ExpressionInfix{
								Token:    tokens.New(tokens.TypePlus, "+"),
								Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
								Right:    &ast.ExpressionLiteralInteger{Value: 1},
								Operator: "+",
							},
						},
						Body: &ast.StatementBlock{
							Statements: []ast.Statement{
								&ast.StatementExpression{
									Token: tokens.New(tokens.TypeIf, "if"),
									Expression: &ast.ExpressionIf{
										Branches: []ast.ConditionalBranch{
											{
												Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
												Condition: &ast.ExpressionInfix{
													Token: tokens.New(tokens.TypeEQ, "=="),
													Left: &ast.ExpressionInfix{
														Token:    tokens.New(tokens.TypePercent, "%"),
														Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
														Right:    &ast.ExpressionLiteralInteger{Value: 2},
														Operator: "%",
													},
													Right:    &ast.ExpressionLiteralInteger{Value: 0},
													Operator: "==",
												},
												Consequence: &ast.StatementBlock{
													Statements: []ast.Statement{
														&ast.StatementExpression{
															Token: tokens.New(tokens.TypeContinue, "continue"),
															Expression: &ast.ExpressionKeyword{
																Token:   tokens.New(tokens.TypeContinue, "continue"),
																Keyword: "continue",
															},
														},
													},
												},
											},
										},
										Alternative: nil,
									},
								},
								&ast.StatementRebind{
									Token: tokens.New(tokens.TypeIdent, "result"),
									Name: &ast.ExpressionIdentifier{
										Token: tokens.New(tokens.TypeIdent, "result"),
										Value: "result",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypePlus, "+"),
										Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "result"), Value: "result"},
										Right:    &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
										Operator: "+",
									},
								},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "result"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "result"),
							Value: "result",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 4},
		},
		{
			name: `
				var result = 0;
				for (var i = 0; i < 6; i = i + 1) {
					if (result >= 3) {
						break;
					}

					result = result + 1;
				}

				result;`,
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "result"),
							Value: "result",
						},
						Value: &ast.ExpressionLiteralInteger{Value: 0},
					},
					&ast.StatementFor{
						Token: tokens.New(tokens.TypeFor, "for"),
						Initializer: &ast.StatementBind{
							Token: tokens.New(tokens.TypeBind, "var"),
							Name: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "i"),
								Value: "i",
							},
							Value: &ast.ExpressionLiteralInteger{Value: 0},
						},
						Condition: &ast.ExpressionInfix{
							Token:    tokens.New(tokens.TypeLT, "<"),
							Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
							Right:    &ast.ExpressionLiteralInteger{Value: 5},
							Operator: "<",
						},
						Step: &ast.StatementRebind{
							Token: tokens.New(tokens.TypeIdent, "i"),
							Name: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "i"),
								Value: "i",
							},
							Value: &ast.ExpressionInfix{
								Token:    tokens.New(tokens.TypePlus, "+"),
								Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "i"), Value: "i"},
								Right:    &ast.ExpressionLiteralInteger{Value: 1},
								Operator: "+",
							},
						},
						Body: &ast.StatementBlock{
							Statements: []ast.Statement{
								&ast.StatementExpression{
									Token: tokens.New(tokens.TypeIf, "if"),
									Expression: &ast.ExpressionIf{
										Branches: []ast.ConditionalBranch{
											{
												Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
												Condition: &ast.ExpressionInfix{
													Token:    tokens.New(tokens.TypeGTEQ, ">="),
													Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "result"), Value: "result"},
													Right:    &ast.ExpressionLiteralInteger{Value: 3},
													Operator: ">=",
												},
												Consequence: &ast.StatementBlock{
													Statements: []ast.Statement{
														&ast.StatementExpression{
															Token: tokens.New(tokens.TypeBreak, "break"),
															Expression: &ast.ExpressionKeyword{
																Token:   tokens.New(tokens.TypeBreak, "break"),
																Keyword: "break",
															},
														},
													},
												},
											},
										},
										Alternative: nil,
									},
								},
								&ast.StatementRebind{
									Token: tokens.New(tokens.TypeIdent, "result"),
									Name: &ast.ExpressionIdentifier{
										Token: tokens.New(tokens.TypeIdent, "result"),
										Value: "result",
									},
									Value: &ast.ExpressionInfix{
										Token:    tokens.New(tokens.TypePlus, "+"),
										Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "result"), Value: "result"},
										Right:    &ast.ExpressionLiteralInteger{Value: 1},
										Operator: "+",
									},
								},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "result"),
						Expression: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "result"),
							Value: "result",
						},
					},
				},
			},
			expected: &objects.Integer{Value: 3},
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

func TestEvalatorEvalFunctions(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "function application with return in if consequence",
			// fn isGTFive(x) {
			//     if (x > 5) {
			//        return true;
			//      }
			//
			//  	return false;
			// }
			//
			// isGTFive(10);
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementFunctionBind{
						Token: tokens.New(tokens.TypeFunction, "fn"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "isGTFive"),
							Value: "isGTFive",
						},
						Value: &ast.ExpressionLiteralFunction{
							Token: tokens.New(tokens.TypeFunction, "fn"),
							Parameters: []*ast.ExpressionIdentifier{
								{
									Token: tokens.New(tokens.TypeIdent, "x"),
									Value: "x",
								},
							},
							Body: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementExpression{
										Token: tokens.New(tokens.TypeIf, "if"),
										Expression: &ast.ExpressionIf{
											Branches: []ast.ConditionalBranch{
												{
													Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
													Condition: &ast.ExpressionInfix{
														Token:    tokens.New(tokens.TypeGT, ">"),
														Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "x"), Value: "x"},
														Right:    &ast.ExpressionLiteralInteger{Value: 5},
														Operator: ">",
													},
													Consequence: &ast.StatementBlock{
														Statements: []ast.Statement{
															&ast.StatementReturn{
																Token: tokens.New(tokens.TypeReturn, "return"),
																Value: &ast.ExpressionLiteralBoolean{Value: true},
															},
														},
													},
												},
											},
										},
									},
									&ast.StatementReturn{
										Token: tokens.New(tokens.TypeReturn, "return"),
										Value: &ast.ExpressionLiteralBoolean{Value: false},
									},
								},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "isGTFive"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "isGTFive"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "isGTFive"),
								Value: "isGTFive",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralInteger{Value: 10},
							},
						},
					},
				},
			},
			expected: &objects.Boolean{Value: true},
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

func TestEvaluatorEvalClosures(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "simple adder closure",
			// fn newAdder(x) {
			//     fn adder(y) {
			//         return x + y;
			//     }
			//
			//     return adder;
			// }
			//
			// var addTwo = newAdder(2);
			// addTwo(3);
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementFunctionBind{
						Token: tokens.New(tokens.TypeFunction, "fn"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "newAdder"),
							Value: "newAdder",
						},
						Value: &ast.ExpressionLiteralFunction{
							Token: tokens.New(tokens.TypeFunction, "fn"),
							Parameters: []*ast.ExpressionIdentifier{
								{
									Token: tokens.New(tokens.TypeIdent, "x"),
									Value: "x",
								},
							},
							Body: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementFunctionBind{
										Token: tokens.New(tokens.TypeFunction, "fn"),
										Name: &ast.ExpressionIdentifier{
											Token: tokens.New(tokens.TypeIdent, "adder"),
											Value: "adder",
										},
										Value: &ast.ExpressionLiteralFunction{
											Token: tokens.New(tokens.TypeFunction, "fn"),
											Parameters: []*ast.ExpressionIdentifier{
												{
													Token: tokens.New(tokens.TypeIdent, "y"),
													Value: "y",
												},
											},
											Body: &ast.StatementBlock{
												Statements: []ast.Statement{
													&ast.StatementReturn{
														Token: tokens.New(tokens.TypeReturn, "return"),
														Value: &ast.ExpressionInfix{
															Token:    tokens.New(tokens.TypePlus, "+"),
															Left:     &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "x"), Value: "x"},
															Right:    &ast.ExpressionIdentifier{Token: tokens.New(tokens.TypeIdent, "y"), Value: "y"},
															Operator: "+",
														},
													},
												},
											},
										},
									},
									&ast.StatementReturn{
										Token: tokens.New(tokens.TypeReturn, "return"),
										Value: &ast.ExpressionIdentifier{
											Token: tokens.New(tokens.TypeIdent, "adder"),
											Value: "adder",
										},
									},
								},
							},
						},
					},
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "addTwo"),
							Value: "addTwo",
						},
						Value: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "newAdder"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "newAdder"),
								Value: "newAdder",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralInteger{Value: 2},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "addTwo"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "addTwo"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "addTwo"),
								Value: "addTwo",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralInteger{Value: 3},
							},
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

func TestEvaluatorEvalBuiltinFunctions(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "len with a string argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "len"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "len"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "len"),
								Value: "len",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralString{Value: "hello"},
							},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 5},
		},
		{
			name: "len with empty string argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "len"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "len"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "len"),
								Value: "len",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralString{Value: ""},
							},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 0},
		},
		{
			name: "len with a list argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "len"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "len"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "len"),
								Value: "len",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralList{
									Token: tokens.New(tokens.TypeLBracket, "["),
									Elements: []ast.Expression{
										&ast.ExpressionLiteralInteger{Value: 1},
										&ast.ExpressionLiteralInteger{Value: 2},
										&ast.ExpressionLiteralInteger{Value: 3},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 3},
		},
		{
			name: "print with a single argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "print"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "print"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "print"),
								Value: "print",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralString{Value: "hello world"},
							},
						},
					},
				},
			},
			expected: &objects.String{Value: "hello world"},
		},
		{
			name: "print with multiple arguments",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "print"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "print"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "print"),
								Value: "print",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralString{Value: "hallo"},
								&ast.ExpressionLiteralString{Value: "welt"},
							},
						},
					},
				},
			},
			expected: &objects.String{Value: "hallo welt"},
		},
		{
			name: "cdr with a non-empty list argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "cdr"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "cdr"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "cdr"),
								Value: "cdr",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralList{
									Token: tokens.New(tokens.TypeLBracket, "["),
									Elements: []ast.Expression{
										&ast.ExpressionLiteralInteger{Value: 1},
										&ast.ExpressionLiteralInteger{Value: 2},
										&ast.ExpressionLiteralInteger{Value: 3},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.List{
				Elements: []objects.Object{
					&objects.Integer{Value: 2},
					&objects.Integer{Value: 3},
				},
			},
		},
		{
			name: "cdr with a empty list argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "cdr"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "cdr"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "cdr"),
								Value: "cdr",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralList{
									Token:    tokens.New(tokens.TypeLBracket, "["),
									Elements: []ast.Expression{},
								},
							},
						},
					},
				},
			},
			expected: &objects.List{
				Elements: []objects.Object{},
			},
		},
		{
			name: "car with a non-empty list argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "car"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "car"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "car"),
								Value: "car",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralList{
									Token: tokens.New(tokens.TypeLBracket, "["),
									Elements: []ast.Expression{
										&ast.ExpressionLiteralInteger{Value: 1},
										&ast.ExpressionLiteralInteger{Value: 2},
										&ast.ExpressionLiteralInteger{Value: 3},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 1},
		},
		{
			name: "car with a empty list argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "car"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "car"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "car"),
								Value: "car",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralList{
									Token:    tokens.New(tokens.TypeLBracket, "["),
									Elements: []ast.Expression{},
								},
							},
						},
					},
				},
			},
			expected: &objects.Nowt{},
		},
		{
			name: "cons with a non-empty list argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "cons"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "cons"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "cons"),
								Value: "cons",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralInteger{Value: 1},
								&ast.ExpressionLiteralList{
									Token: tokens.New(tokens.TypeLBracket, "["),
									Elements: []ast.Expression{
										&ast.ExpressionLiteralInteger{Value: 2},
										&ast.ExpressionLiteralInteger{Value: 3},
									},
								},
							},
						},
					},
				},
			},
			expected: &objects.List{
				Elements: []objects.Object{
					&objects.Integer{Value: 1},
					&objects.Integer{Value: 2},
					&objects.Integer{Value: 3},
				},
			},
		},
		{
			name: "cons with an empty list argument",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "cons"),
						Expression: &ast.ExpressionCall{
							Token: tokens.New(tokens.TypeIdent, "cons"),
							Function: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "cons"),
								Value: "cons",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionLiteralInteger{Value: 1},
								&ast.ExpressionLiteralList{
									Token:    tokens.New(tokens.TypeLBracket, "["),
									Elements: []ast.Expression{},
								},
							},
						},
					},
				},
			},
			expected: &objects.List{
				Elements: []objects.Object{
					&objects.Integer{Value: 1},
				},
			},
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

func TestEvaluatorEvalIndexExpressionsList(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "indexing into an array",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "arr"),
							Value: "arr",
						},
						Value: &ast.ExpressionLiteralList{
							Token: tokens.New(tokens.TypeLBracket, "["),
							Elements: []ast.Expression{
								&ast.ExpressionLiteralInteger{Value: 1},
								&ast.ExpressionLiteralInteger{Value: 2},
								&ast.ExpressionLiteralInteger{Value: 3},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "arr"),
						Expression: &ast.ExpressionIndex{
							Token: tokens.New(tokens.TypeLBracket, "["),
							Left: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "arr"),
								Value: "arr",
							},
							Index: &ast.ExpressionLiteralInteger{Value: 1},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 2},
		},
		{
			name: "negative indexing into an array within bounds",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "arr"),
							Value: "arr",
						},
						Value: &ast.ExpressionLiteralList{
							Token: tokens.New(tokens.TypeLBracket, "["),
							Elements: []ast.Expression{
								&ast.ExpressionLiteralInteger{Value: 1},
								&ast.ExpressionLiteralInteger{Value: 2},
								&ast.ExpressionLiteralInteger{Value: 3},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "arr"),
						Expression: &ast.ExpressionIndex{
							Token: tokens.New(tokens.TypeLBracket, "["),
							Left: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "arr"),
								Value: "arr",
							},
							Index: &ast.ExpressionLiteralInteger{Value: -1},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 3},
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

func TestEvaluatorEvalMaps(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected string
	}

	tests := []testCase{
		{
			name: "map literal with string keys and integer values",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "myMap"),
						Expression: &ast.ExpressionLiteralMap{
							Token: tokens.New(tokens.TypeLBrace, "{"),
							Pairs: []ast.MapPair{
								{
									Key:   &ast.ExpressionLiteralString{Value: "one"},
									Value: &ast.ExpressionLiteralInteger{Value: 1},
								},
								{
									Key:   &ast.ExpressionLiteralString{Value: "two"},
									Value: &ast.ExpressionLiteralInteger{Value: 2},
								},
							},
						},
					},
				},
			},
			expected: `{"one": 1, "two": 2}`,
		},
		{
			name: "map literal with mixed key and value types",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "myMap"),
						Expression: &ast.ExpressionLiteralMap{
							Token: tokens.New(tokens.TypeLBrace, "{"),
							Pairs: []ast.MapPair{
								{
									Key:   &ast.ExpressionLiteralInteger{Value: 1},
									Value: &ast.ExpressionLiteralString{Value: "one"},
								},
								{
									Key:   &ast.ExpressionLiteralString{Value: "two"},
									Value: &ast.ExpressionLiteralBoolean{Value: true},
								},
							},
						},
					},
				},
			},
			expected: `{"1": "one", "two": true}`,
		},
		{
			name: "empty map literal",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "myMap"),
						Expression: &ast.ExpressionLiteralMap{
							Token: tokens.New(tokens.TypeLBrace, "{"),
							Pairs: []ast.MapPair{},
						},
					},
				},
			},
			expected: `{}`,
		},
	}

	for i, tc := range tests {
		t.Run(fmt.Sprintf("test %d: %s", i, tc.name), func(t *testing.T) {
			evaluator := New(nil)
			result := evaluator.Eval(tc.input, objects.NewEnvironment())

			mapObj, ok := result.(*objects.Map)
			assert.True(t, ok, "result should be of type *objects.Map")

			mapJSON := mapObj.Inspect()
			assert.JSONEq(t, tc.expected, mapJSON)
		})
	}
}

func TestEvaluatorEvalIndexExpressionsMap(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "indexing into a map with a string key",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.New(tokens.TypeBind, "var"),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "myMap"),
							Value: "myMap",
						},
						Value: &ast.ExpressionLiteralMap{
							Token: tokens.New(tokens.TypeLBrace, "{"),
							Pairs: []ast.MapPair{
								{
									Key:   &ast.ExpressionLiteralString{Value: "one"},
									Value: &ast.ExpressionLiteralInteger{Value: 1},
								},
							},
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeIdent, "myMap"),
						Expression: &ast.ExpressionIndex{
							Token: tokens.New(tokens.TypeLBracket, "["),
							Left: &ast.ExpressionIdentifier{
								Token: tokens.New(tokens.TypeIdent, "myMap"),
								Value: "myMap",
							},
							Index: &ast.ExpressionLiteralString{Value: "one"},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 1},
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

func TestEvaluateIndexBinding(t *testing.T) {
	type testCase struct {
		name     string
		input    *ast.Program
		expected objects.Object
	}

	tests := []testCase{
		{
			name: "nested map key binding to list index and getter - {key1: {key2: 5}}[key1][key2] = 42",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.NewStatic(tokens.TypeBind),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "m"),
							Value: "m",
						},
						Value: &ast.ExpressionLiteralMap{
							Token: tokens.New(tokens.TypeLBrace, "{"),
							Pairs: []ast.MapPair{
								{
									Key: &ast.ExpressionLiteralString{Value: "key1"},
									Value: &ast.ExpressionLiteralMap{
										Token: tokens.New(tokens.TypeLBrace, "{"),
										Pairs: []ast.MapPair{
											{
												Key:   &ast.ExpressionLiteralString{Value: "key2"},
												Value: &ast.ExpressionLiteralInteger{Value: 5},
											},
										},
									},
								},
							},
						},
					},
					&ast.StatementIndexBind{
						Token: tokens.NewStatic(tokens.TypeBind),
						Left: &ast.ExpressionIndex{
							Token: tokens.New(tokens.TypeLBracket, "["),
							Left: &ast.ExpressionIndex{
								Token: tokens.New(tokens.TypeLBracket, "["),
								Left: &ast.ExpressionIdentifier{
									Token: tokens.New(tokens.TypeIdent, "m"),
									Value: "m",
								},
								Index: &ast.ExpressionLiteralString{
									Token: tokens.New(tokens.TypeString, "key1"),
									Value: "key1"},
							},
							Index: &ast.ExpressionLiteralString{
								Token: tokens.New(tokens.TypeString, "key2"),
								Value: "key2",
							},
						},
						Right: &ast.ExpressionLiteralInteger{
							Token: tokens.New(tokens.TypeInt, "5"),
							Value: 42,
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeLBracket, "["),
						Expression: &ast.ExpressionIndex{
							Token: tokens.New(tokens.TypeLBracket, "["),
							Left: &ast.ExpressionIndex{
								Token: tokens.New(tokens.TypeLBracket, "["),
								Left: &ast.ExpressionIdentifier{
									Token: tokens.New(tokens.TypeIdent, "m"),
									Value: "m",
								},
								Index: &ast.ExpressionLiteralString{
									Token: tokens.New(tokens.TypeString, "key1"),
									Value: "key1",
								},
							},
							Index: &ast.ExpressionLiteralString{
								Token: tokens.New(tokens.TypeString, "key2"),
								Value: "key2",
							},
						},
					},
				},
			},
			expected: &objects.Integer{Value: 42},
		},
		{
			name: "nested map key binding to list index and getter - {key1: {key2: [5]}}[key1][key2][0] = 42",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.NewStatic(tokens.TypeBind),
						Name: &ast.ExpressionIdentifier{
							Token: tokens.New(tokens.TypeIdent, "m"),
							Value: "m",
						},
						Value: &ast.ExpressionLiteralMap{
							Token: tokens.New(tokens.TypeLBrace, "{"),
							Pairs: []ast.MapPair{
								{
									Key: &ast.ExpressionLiteralString{Value: "key1"},
									Value: &ast.ExpressionLiteralMap{
										Token: tokens.New(tokens.TypeLBrace, "{"),
										Pairs: []ast.MapPair{
											{
												Key: &ast.ExpressionLiteralString{Value: "key2"},
												Value: &ast.ExpressionLiteralList{
													Token: tokens.New(tokens.TypeLBracket, "["),
													Elements: []ast.Expression{
														&ast.ExpressionLiteralInteger{Value: 5},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					&ast.StatementIndexBind{
						Token: tokens.NewStatic(tokens.TypeBind),
						Left: &ast.ExpressionIndex{
							Token: tokens.NewStatic(tokens.TypeLBracket),
							Left: &ast.ExpressionIndex{
								Token: tokens.NewStatic(tokens.TypeLBracket),
								Left: &ast.ExpressionIndex{
									Token: tokens.NewStatic(tokens.TypeLBracket),
									Left: &ast.ExpressionIdentifier{
										Token: tokens.New(tokens.TypeIdent, "m"),
										Value: "m",
									},
									Index: &ast.ExpressionLiteralString{
										Token: tokens.New(tokens.TypeString, "key1"),
										Value: "key1",
									},
								},
								Index: &ast.ExpressionLiteralString{
									Token: tokens.New(tokens.TypeString, "key2"),
									Value: "key2",
								},
							},
							Index: &ast.ExpressionLiteralInteger{
								Token: tokens.New(tokens.TypeInt, "0"),
								Value: 0,
							},
						},
						Right: &ast.ExpressionLiteralInteger{
							Token: tokens.New(tokens.TypeInt, ""),
							Value: 42,
						},
					},
					&ast.StatementExpression{
						Token: tokens.New(tokens.TypeLBracket, "["),
						Expression: &ast.ExpressionIndex{
							Token: tokens.New(tokens.TypeLBracket, "["),
							Left: &ast.ExpressionIndex{
								Token: tokens.New(tokens.TypeLBracket, "["),
								Left: &ast.ExpressionIdentifier{
									Token: tokens.New(tokens.TypeIdent, "m"),
									Value: "m",
								},
								Index: &ast.ExpressionLiteralString{
									Token: tokens.New(tokens.TypeString, "key1"),
									Value: "key1",
								},
							},
							Index: &ast.ExpressionLiteralString{
								Token: tokens.New(tokens.TypeString, "key2"),
								Value: "key2",
							},
						},
					},
				},
			},
			expected: &objects.List{
				Elements: []objects.Object{
					&objects.Integer{Value: 42},
				},
			},
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
