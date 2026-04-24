package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/lexer"
	"github.com/OJOMB/donkey/internal/tokens"
)

func TestParseStatements(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name: "test bind statements - no errors",
			input: `
					var x = 5;
					var y = "hello";
					var __foobar__ = false;
					var myFunction = fn(x) { return x + 1; };
				`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBinder, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "x"},
							Value: "x",
						},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: "INT", Lexeme: "5"},
							Value: 5,
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBinder, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "y"},
							Value: "y",
						},
						Value: &ast.ExpressionLiteralString{
							Token: tokens.Token{Type: "STRING", Lexeme: "hello"},
							Value: "hello",
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBinder, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "__foobar__"},
							Value: "__foobar__",
						},
						Value: &ast.ExpressionLiteralBoolean{
							Token: tokens.Token{Type: "FALSE", Lexeme: "false"},
							Value: false,
						},
					},
					&ast.StatementBind{
						Token: tokens.Token{Type: tokens.TypeBinder, Lexeme: "var"},
						Name: &ast.ExpressionIdentifier{
							Token: tokens.Token{Type: "IDENT", Lexeme: "myFunction"},
							Value: "myFunction",
						},
						Value: &ast.ExpressionLiteralFunction{
							Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
							Parameters: []*ast.ExpressionIdentifier{
								{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
									Value: "x",
								},
							},
							Body: &ast.StatementBlock{Statements: []ast.Statement{
								&ast.StatementReturn{
									Token: tokens.Token{Type: tokens.TypeReturn, Lexeme: "return"},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionLiteralInteger{
											Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "1"},
											Value: 1,
										},
									},
								},
							}},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name: "test return statements",
			input: `
				return 5;
				return "kool";
				return true;
			`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementReturn{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						Value: &ast.ExpressionLiteralInteger{
							Token: tokens.Token{Type: "INT", Lexeme: "5"},
							Value: 5,
						},
					},
					&ast.StatementReturn{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						Value: &ast.ExpressionLiteralString{
							Token: tokens.Token{Type: "STRING", Lexeme: "kool"},
							Value: "kool",
						},
					},
					&ast.StatementReturn{
						Token: tokens.Token{Type: "RETURN", Lexeme: "return"},
						Value: &ast.ExpressionLiteralBoolean{
							Token: tokens.Token{Type: "TRUE", Lexeme: "true"},
							Value: true,
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		// TODO see how complex it is to add support for rebind statements
		// {
		// 	name: "rebind statements",
		// 	input: `
		// 		x = 5;
		// 		y = "kool";
		// 		__foobar__ = true;
		// 	`,
		// 	expectedOutput: &ast.Program{
		// 		Statements: []ast.Statement{},
		// 	},
		// 	expectedErrs: []string{
		// 		"unexpected token: =",
		// 		"unexpected token: =",
		// 		"unexpected token: =",
		// 	},
		// },
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p, err := New(lexer.New(input, nil), nil)
	require.NoError(t, err)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.StatementExpression)
	require.True(t, ok)

	ident, ok := stmt.Expression.(*ast.ExpressionIdentifier)
	require.True(t, ok)
	assert.Equal(t, "foobar", ident.Value)
	assert.Equal(t, "foobar", ident.TokenLexeme())
}

func TestIntegerExpression(t *testing.T) {
	input := "5;"

	p, err := New(lexer.New(input, nil), nil)
	require.NoError(t, err)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.StatementExpression)
	require.True(t, ok)

	intLiteral, ok := stmt.Expression.(*ast.ExpressionLiteralInteger)
	require.True(t, ok)
	assert.Equal(t, 5, intLiteral.Value)
	assert.Equal(t, "5", intLiteral.TokenLexeme())
	assert.Equal(t, "5", intLiteral.String())
}

func TestStringExpression(t *testing.T) {
	input := `"hello";`

	p, err := New(lexer.New(input, nil), nil)
	require.NoError(t, err)

	program := p.ParseProgram()
	assert.NotNil(t, program)
	require.Len(t, program.Statements, 1)

	stmt, ok := program.Statements[0].(*ast.StatementExpression)
	require.True(t, ok)

	strLiteral, ok := stmt.Expression.(*ast.ExpressionLiteralString)
	require.True(t, ok)
	assert.Equal(t, "hello", strLiteral.Value)
	assert.Equal(t, "hello", strLiteral.TokenLexeme())
	assert.Equal(t, "hello", strLiteral.String())
}

func TestExpressionPrefix(t *testing.T) {
	t.Run("prefix expression: !", func(t *testing.T) {
		input := "!5;"

		p, err := New(lexer.New(input, nil), nil)
		require.NoError(t, err)

		program := p.ParseProgram()
		assert.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.StatementExpression)
		require.True(t, ok)

		assert.IsType(t, &ast.ExpressionPrefix{}, stmt.Expression)

		prefixExp, ok := stmt.Expression.(*ast.ExpressionPrefix)
		require.True(t, ok)
		assert.Equal(t, "!", prefixExp.Operator)
		assert.Equal(t, 5, prefixExp.Right.(*ast.ExpressionLiteralInteger).Value)
		assert.Equal(t, "!", prefixExp.TokenLexeme())
	})

	t.Run("prefix expression: -", func(t *testing.T) {
		input := "-15;"

		p, err := New(lexer.New(input, nil), nil)
		require.NoError(t, err)

		program := p.ParseProgram()
		assert.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.StatementExpression)
		require.True(t, ok)

		assert.IsType(t, &ast.ExpressionPrefix{}, stmt.Expression)

		prefixExp, ok := stmt.Expression.(*ast.ExpressionPrefix)
		require.True(t, ok)
		assert.Equal(t, "-", prefixExp.Operator)
		assert.Equal(t, 15, prefixExp.Right.(*ast.ExpressionLiteralInteger).Value)
		assert.Equal(t, "-", prefixExp.TokenLexeme())
	})
}

func TestExpressionStatementBool(t *testing.T) {
	t.Run("simple boolean literal: true", func(t *testing.T) {
		input := "true;"

		p, err := New(lexer.New(input, nil), nil)
		require.NoError(t, err)

		program := p.ParseProgram()
		assert.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.StatementExpression)
		require.True(t, ok)

		assert.IsType(t, &ast.ExpressionLiteralBoolean{}, stmt.Expression)

		boolExp, ok := stmt.Expression.(*ast.ExpressionLiteralBoolean)
		require.True(t, ok)
		assert.Equal(t, true, boolExp.Value)
		assert.Equal(t, "true", boolExp.TokenLexeme())
	})

	t.Run("simple boolean literal: false", func(t *testing.T) {
		input := "false;"

		p, err := New(lexer.New(input, nil), nil)
		require.NoError(t, err)

		program := p.ParseProgram()
		assert.NotNil(t, program)
		require.Len(t, program.Statements, 1)

		stmt, ok := program.Statements[0].(*ast.StatementExpression)
		require.True(t, ok)

		assert.IsType(t, &ast.ExpressionLiteralBoolean{}, stmt.Expression)

		boolExp, ok := stmt.Expression.(*ast.ExpressionLiteralBoolean)
		require.True(t, ok)
		assert.Equal(t, false, boolExp.Value)
		assert.Equal(t, "false", boolExp.TokenLexeme())
	})
}

func TestParsingInfixExpressions(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "simple infix expressions - no errors",
			input: `5 + 122;`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
								Value: 5,
							},
							Right: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "122"},
								Value: 122,
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "slightly more complex infix expression - no errors",
			input: `5 + 5 / 10 * 4;`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionLiteralInteger{
								Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
								Value: 5,
							},
							Right: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypeAsterisk, Lexeme: "*"},
								Operator: "*",
								Left: &ast.ExpressionInfix{
									Token:    tokens.Token{Type: tokens.TypeForwardSlash, Lexeme: "/"},
									Operator: "/",
									Left: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
										Value: 5,
									},
									Right: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
										Value: 10,
									},
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "4"},
									Value: 4,
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "slightly more complex infix expression - no errors",
			input: `5 * 5 + 10 / 4;`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
							Operator: "+",
							Left: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypeAsterisk, Lexeme: "*"},
								Operator: "*",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
									Value: 5,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
									Value: 5,
								},
							},
							Right: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypeForwardSlash, Lexeme: "/"},
								Operator: "/",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
									Value: 10,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "4"},
									Value: 4,
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "grouped infix expression - no errors",
			input: `(5 + 5) * (10 / 4);`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeLParen, Lexeme: "("},
						Expression: &ast.ExpressionInfix{
							Token:    tokens.Token{Type: tokens.TypeAsterisk, Lexeme: "*"},
							Operator: "*",
							Left: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
								Operator: "+",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
									Value: 5,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
									Value: 5,
								},
							},
							Right: &ast.ExpressionInfix{
								Token:    tokens.Token{Type: tokens.TypeForwardSlash, Lexeme: "/"},
								Operator: "/",
								Left: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "10"},
									Value: 10,
								},
								Right: &ast.ExpressionLiteralInteger{
									Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "4"},
									Value: 4,
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestIfExpression(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "if expression without else - no errors",
			input: `if (x < y) { x }`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypeLT, Lexeme: "<"},
										Operator: "<",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
												Expression: &ast.ExpressionIdentifier{
													Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
													Value: "x",
												},
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
			expectedErrs: []string{},
		},
		{
			name:  "if expression with elif and else - no errors",
			input: `if (x < y) { x } elif (x > y) { y } else { 5 }`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
						Expression: &ast.ExpressionIf{
							Branches: []ast.ConditionalBranch{
								{
									Token: tokens.Token{Type: tokens.TypeIf, Lexeme: "if"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypeLT, Lexeme: "<"},
										Operator: "<",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
												Expression: &ast.ExpressionIdentifier{
													Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
													Value: "x",
												},
											},
										},
									},
								},
								{
									Token: tokens.Token{Type: tokens.TypeElif, Lexeme: "elif"},
									Condition: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypeGT, Lexeme: ">"},
										Operator: ">",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
									Consequence: &ast.StatementBlock{
										Statements: []ast.Statement{
											&ast.StatementExpression{
												Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
												Expression: &ast.ExpressionIdentifier{
													Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
													Value: "y",
												},
											},
										},
									},
								},
							},
							Alternative: &ast.StatementBlock{
								Statements: []ast.Statement{
									&ast.StatementExpression{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
										Expression: &ast.ExpressionLiteralInteger{
											Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
											Value: 5,
										},
									},
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestFunctionLiterals(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "function literal with no parameters - no errors",
			input: `fn() { return 5; }`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
						Expression: &ast.ExpressionLiteralFunction{
							Token:      tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
							Parameters: []*ast.ExpressionIdentifier{},
							Body: &ast.StatementBlock{Statements: []ast.Statement{
								&ast.StatementReturn{
									Token: tokens.Token{Type: tokens.TypeReturn, Lexeme: "return"},
									Value: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "5"},
										Value: 5,
									},
								},
							}},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
		{
			name:  "function literal with parameters - no errors",
			input: `fn(x, y) { return x + y; }`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
						Expression: &ast.ExpressionLiteralFunction{
							Token: tokens.Token{Type: tokens.TypeFunction, Lexeme: "fn"},
							Parameters: []*ast.ExpressionIdentifier{
								{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
									Value: "x",
								},
								{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
									Value: "y",
								},
							},
							Body: &ast.StatementBlock{Statements: []ast.Statement{
								&ast.StatementReturn{
									Token: tokens.Token{Type: tokens.TypeReturn, Lexeme: "return"},
									Value: &ast.ExpressionInfix{
										Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
										Operator: "+",
										Left: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "x"},
											Value: "x",
										},
										Right: &ast.ExpressionIdentifier{
											Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "y"},
											Value: "y",
										},
									},
								},
							}},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestParsingFunctionCalls(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput *ast.Program
		expectedErrs   []string
	}

	var testCases = []testCase{
		{
			name:  "function call with args - no errors",
			input: `myFunction(2+3, param2, false);`,
			expectedOutput: &ast.Program{
				Statements: []ast.Statement{
					&ast.StatementExpression{
						Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myFunction"},
						Expression: &ast.ExpressionCall{
							Token: tokens.Token{Type: tokens.TypeLParen, Lexeme: "("},
							Function: &ast.ExpressionIdentifier{
								Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "myFunction"},
								Value: "myFunction",
							},
							Arguments: []ast.Expression{
								&ast.ExpressionInfix{
									Token:    tokens.Token{Type: tokens.TypePlus, Lexeme: "+"},
									Operator: "+",
									Left: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "2"},
										Value: 2,
									},
									Right: &ast.ExpressionLiteralInteger{
										Token: tokens.Token{Type: tokens.TypeInt, Lexeme: "3"},
										Value: 3,
									},
								},
								&ast.ExpressionIdentifier{
									Token: tokens.Token{Type: tokens.TypeIdent, Lexeme: "param2"},
									Value: "param2",
								},
								&ast.ExpressionLiteralBoolean{
									Token: tokens.Token{Type: tokens.TypeFalse, Lexeme: "false"},
									Value: false,
								},
							},
						},
					},
				},
			},
			expectedErrs: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expectedOutput, program)
			assert.Equal(t, tc.expectedErrs, p.Errors)
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput string
	}

	var testCases = []testCase{
		{
			name:           "operator precedence parsing 1",
			input:          `-a * b`,
			expectedOutput: `((-a) * b)`,
		},
		{
			name:           "operator precedence parsing 2",
			input:          `!-a`,
			expectedOutput: `(!(-a))`,
		},
		{
			name:           "operator precedence parsing 3",
			input:          `a + b + c`,
			expectedOutput: `((a + b) + c)`,
		},
		{
			name:           "operator precedence parsing 4",
			input:          `a + b - c;`,
			expectedOutput: `((a + b) - c)`,
		},
		{
			name:           "operator precedence parsing 5",
			input:          `a * b / c`,
			expectedOutput: `((a * b) / c)`,
		},
		{
			name:           "operator precedence parsing 6",
			input:          `a + b * c + d / e - f`,
			expectedOutput: `(((a + (b * c)) + (d / e)) - f)`,
		},
		{
			name:           "operator precedence parsing 7",
			input:          `3 + 4; -5 * 5`,
			expectedOutput: "(3 + 4)\n((-5) * 5)",
		},
		{
			name:           "operator precedence parsing 8",
			input:          `5 > 4 == 3 < 4`,
			expectedOutput: `((5 > 4) == (3 < 4))`,
		},
		{
			name:           "operator precedence parsing 9",
			input:          `5 < 4 != 3 > 4`,
			expectedOutput: `((5 < 4) != (3 > 4))`,
		},
		{
			name:           "operator precedence parsing 10",
			input:          `3 + 4 * 5 == 3 * 1 + 4 * 5`,
			expectedOutput: `((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))`,
		},
		{
			name:           "operator precedence parsing 11",
			input:          `3 + 4 * 5 == add(add(3 * 1), 4 * 5)`,
			expectedOutput: `((3 + (4 * 5)) == add(add((3 * 1)), (4 * 5)))`,
		},
		{
			name:           "operator precedence parsing 12",
			input:          `-a ^ b`,
			expectedOutput: `((-a) ^ b)`,
		},
		{
			name:           "operator precedence parsing 13",
			input:          `a ^ b + c`,
			expectedOutput: `((a ^ b) + c)`,
		},
		{
			name:           "operator precedence parsing 14",
			input:          `a + b ^ c + d`,
			expectedOutput: `((a + (b ^ c)) + d)`,
		},
		{
			name:           "operator precedence parsing 15",
			input:          `a && b || c && d `,
			expectedOutput: `((a && b) || (c && d))`,
		},
		{
			name:           "operator precedence parsing 16",
			input:          `a || b && c || d && e`,
			expectedOutput: `((a || (b && c)) || (d && e))`,
		},
		{
			name:           "operator precedence parsing 17",
			input:          `a | b & c | d`,
			expectedOutput: `((a | (b & c)) | d)`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p, err := New(lexer.New(tc.input, nil), nil)
			require.NoError(t, err)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			prgrmStr := program.String()
			assert.Equal(t, tc.expectedOutput, prgrmStr)
		})
	}
}
