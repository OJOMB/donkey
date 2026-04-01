package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/OJOMB/donkey/internal/tokens"
)

func TestNextToken1(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput []tokens.Token
	}

	var testCases = []testCase{
		{
			name: "test next token with all token types",
			input: `
				let five = 5;
				let ten10 = 10;

				let add = fn(x, y) {
					let sum = x + y;

					return sum;
				};

				if (five < 10) {
					let __one = false;
				} elif (ten10 > 10) {
				 	let __two = true;
				} elif (ten10 >= 10) {
					let __twoPointFive = true;
				} elif (ten10 <= 10) {
					let __twoPointSix = true;
				} elif (ten10 == 10) {
					let __three = false;
				} elif (ten10 != 10) {
					let __four = true;
				} else {
					let __five = false;
				}

				let result = add(five, ten);`,
			expectedOutput: []tokens.Token{
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "five"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeInt, Lexeme: "5"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},

				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "ten10"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeInt, Lexeme: "10"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},

				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "add"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeFunction, Lexeme: "fn"},
				{Type: tokens.TypeLParen, Lexeme: "("},
				{Type: tokens.TypeIdent, Lexeme: "x"},
				{Type: tokens.TypeComma, Lexeme: ","},
				{Type: tokens.TypeIdent, Lexeme: "y"},
				{Type: tokens.TypeRParen, Lexeme: ")"},
				{Type: tokens.TypeLBrace, Lexeme: "{"},
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "sum"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeIdent, Lexeme: "x"},
				{Type: tokens.TypePlus, Lexeme: "+"},
				{Type: tokens.TypeIdent, Lexeme: "y"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeReturn, Lexeme: "return"},
				{Type: tokens.TypeIdent, Lexeme: "sum"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeRBrace, Lexeme: "}"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},

				{Type: tokens.TypeIf, Lexeme: "if"},
				{Type: tokens.TypeLParen, Lexeme: "("},
				{Type: tokens.TypeIdent, Lexeme: "five"},
				{Type: tokens.TypeLT, Lexeme: "<"},
				{Type: tokens.TypeInt, Lexeme: "10"},
				{Type: tokens.TypeRParen, Lexeme: ")"},
				{Type: tokens.TypeLBrace, Lexeme: "{"},
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "__one"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeFalse, Lexeme: "false"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeRBrace, Lexeme: "}"},
				{Type: tokens.TypeElif, Lexeme: "elif"},
				{Type: tokens.TypeLParen, Lexeme: "("},
				{Type: tokens.TypeIdent, Lexeme: "ten10"},
				{Type: tokens.TypeGT, Lexeme: ">"},
				{Type: tokens.TypeInt, Lexeme: "10"},
				{Type: tokens.TypeRParen, Lexeme: ")"},
				{Type: tokens.TypeLBrace, Lexeme: "{"},
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "__two"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeTrue, Lexeme: "true"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeRBrace, Lexeme: "}"},
				{Type: tokens.TypeElif, Lexeme: "elif"},
				{Type: tokens.TypeLParen, Lexeme: "("},
				{Type: tokens.TypeIdent, Lexeme: "ten10"},
				{Type: tokens.TypeGTEQ, Lexeme: ">="},
				{Type: tokens.TypeInt, Lexeme: "10"},
				{Type: tokens.TypeRParen, Lexeme: ")"},
				{Type: tokens.TypeLBrace, Lexeme: "{"},
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "__twoPointFive"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeTrue, Lexeme: "true"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeRBrace, Lexeme: "}"},
				{Type: tokens.TypeElif, Lexeme: "elif"},
				{Type: tokens.TypeLParen, Lexeme: "("},
				{Type: tokens.TypeIdent, Lexeme: "ten10"},
				{Type: tokens.TypeLTEQ, Lexeme: "<="},
				{Type: tokens.TypeInt, Lexeme: "10"},
				{Type: tokens.TypeRParen, Lexeme: ")"},
				{Type: tokens.TypeLBrace, Lexeme: "{"},
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "__twoPointSix"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeTrue, Lexeme: "true"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeRBrace, Lexeme: "}"},
				{Type: tokens.TypeElif, Lexeme: "elif"},
				{Type: tokens.TypeLParen, Lexeme: "("},
				{Type: tokens.TypeIdent, Lexeme: "ten10"},
				{Type: tokens.TypeEq, Lexeme: "=="},
				{Type: tokens.TypeInt, Lexeme: "10"},
				{Type: tokens.TypeRParen, Lexeme: ")"},
				{Type: tokens.TypeLBrace, Lexeme: "{"},
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "__three"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeFalse, Lexeme: "false"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeRBrace, Lexeme: "}"},
				{Type: tokens.TypeElif, Lexeme: "elif"},
				{Type: tokens.TypeLParen, Lexeme: "("},
				{Type: tokens.TypeIdent, Lexeme: "ten10"},
				{Type: tokens.TypeNotEq, Lexeme: "!="},
				{Type: tokens.TypeInt, Lexeme: "10"},
				{Type: tokens.TypeRParen, Lexeme: ")"},
				{Type: tokens.TypeLBrace, Lexeme: "{"},
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "__four"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeTrue, Lexeme: "true"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeRBrace, Lexeme: "}"},
				{Type: tokens.TypeElse, Lexeme: "else"},
				{Type: tokens.TypeLBrace, Lexeme: "{"},
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "__five"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeFalse, Lexeme: "false"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeRBrace, Lexeme: "}"},

				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "result"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeIdent, Lexeme: "add"},
				{Type: tokens.TypeLParen, Lexeme: "("},
				{Type: tokens.TypeIdent, Lexeme: "five"},
				{Type: tokens.TypeComma, Lexeme: ","},
				{Type: tokens.TypeIdent, Lexeme: "ten"},
				{Type: tokens.TypeRParen, Lexeme: ")"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeEOF, Lexeme: ""},
			},
		},
	}

	for i, tc := range testCases {
		lex := New(tc.input)

		// call NextToken until we get an EOF token - assuming the lexer is working correctly we should get the expected output tokens in order
		var toks []tokens.Token
		for {
			tok := lex.NextToken()
			toks = append(toks, tok)

			if tok.Type == tokens.TypeEOF {
				break
			}
		}

		assert.Equal(t, tc.expectedOutput, toks, "test case %d failed", i)
	}
}

func TestStringLiteral(t *testing.T) {
	type testCase struct {
		name           string
		input          string
		expectedOutput []tokens.Token
	}

	var testCases = []testCase{
		{
			name:  "test string literal token",
			input: `let myString = "foobar";`,
			expectedOutput: []tokens.Token{
				{Type: tokens.TypeBinder, Lexeme: "let"},
				{Type: tokens.TypeIdent, Lexeme: "myString"},
				{Type: tokens.TypeAssign, Lexeme: "="},
				{Type: tokens.TypeString, Lexeme: "foobar"},
				{Type: tokens.TypeSemicolon, Lexeme: ";"},
				{Type: tokens.TypeEOF, Lexeme: ""},
			},
		},
	}

	for i, tc := range testCases {
		lex := New(tc.input)

		var toks []tokens.Token
		for {
			tok := lex.NextToken()
			toks = append(toks, tok)
			if tok.Type == tokens.TypeEOF {
				break
			}
		}

		assert.Equal(t, tc.expectedOutput, toks, "test case %d failed", i)
	}
}
