package parser

import "fmt"

var (
	// ErrLexerUnitialized is the error message returned when the parser is initialized with a nil lexer.
	ErrLexerUnitialized = fmt.Errorf("lexer is uninitialized")
)
