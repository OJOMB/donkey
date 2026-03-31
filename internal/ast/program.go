package ast

import "strings"

// Program is the root node of the AST. It contains a slice of statements.
type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	return &Program{
		Statements: make([]Statement, 0),
	}
}

// TokenLexeme returns the lexeme of the first statement's token in the program, or an empty string if there are no statements.
func (p *Program) TokenLexeme() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLexeme()
	}

	return ""
}

// String returns a string representation of the program by concatenating the string representations of all its statements.
func (p *Program) String() string {
	var result = strings.Builder{}
	for _, stmt := range p.Statements {
		if _, err := result.WriteString(stmt.String()); err != nil {
			return "failed to write program string representation"
		}
	}

	return result.String()
}
