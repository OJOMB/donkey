package lexer

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/OJOMB/monkey/internal/lexer/tokens"
)

func TestNextToken1(t *testing.T) {
	type testCase struct {
		input          string
		expectedOutput []tokens.Token
	}

	var testCases = []testCase{
		{
			input: `
				let five = 5;
				let ten10 = 10;

				let add = fn(x, y) {
					x + y;
				};

				let result = add(five, ten);`,
			expectedOutput: []tokens.Token{
				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "five"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeInt, Lexeme: "5"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},

				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "ten10"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeInt, Lexeme: "10"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},

				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "add"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeFunction, Lexeme: "fn"},
				{Type: tokens.TokenTypeLParen, Lexeme: "("},
				{Type: tokens.TokenTypeIdent, Lexeme: "x"},
				{Type: tokens.TokenTypeComma, Lexeme: ","},
				{Type: tokens.TokenTypeIdent, Lexeme: "y"},
				{Type: tokens.TokenTypeRParen, Lexeme: ")"},
				{Type: tokens.TokenTypeLBrace, Lexeme: "{"},
				{Type: tokens.TokenTypeIdent, Lexeme: "x"},
				{Type: tokens.TokenTypePlus, Lexeme: "+"},
				{Type: tokens.TokenTypeIdent, Lexeme: "y"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeRBrace, Lexeme: "}"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},

				{Type: tokens.TokenTypeLet, Lexeme: "let"},
				{Type: tokens.TokenTypeIdent, Lexeme: "result"},
				{Type: tokens.TokenTypeAssign, Lexeme: "="},
				{Type: tokens.TokenTypeIdent, Lexeme: "add"},
				{Type: tokens.TokenTypeLParen, Lexeme: "("},
				{Type: tokens.TokenTypeIdent, Lexeme: "five"},
				{Type: tokens.TokenTypeComma, Lexeme: ","},
				{Type: tokens.TokenTypeIdent, Lexeme: "ten"},
				{Type: tokens.TokenTypeRParen, Lexeme: ")"},
				{Type: tokens.TokenTypeSemicolon, Lexeme: ";"},
				{Type: tokens.TokenTypeEOF, Lexeme: ""},
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

			if tok.Type == tokens.TokenTypeEOF {
				break
			}
		}

		assert.Equal(t, tc.expectedOutput, toks, "test case %d failed", i)
	}
}
