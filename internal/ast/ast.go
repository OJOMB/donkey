package ast

import "strings"

// Node is the base interface for all nodes in the AST.
// It has a method TokenLexeme() that returns the lexeme of the token associated with the node.
type Node interface {
	// TokenLexeme returns the lexeme of the token associated with the node.
	TokenLexeme() string
	// String returns a string representation of the node.
	String() string
}

// Statement represents a statement in the AST.
// It embeds the Node interface and has an additional method statementNode() that is used to differentiate it from expressions.
type Statement interface {
	Node
	// statementNode is a marker method to indicate that a struct is a Statement.
	statementNode()
}

// Expression represents an expression in the AST.
// It embeds the Node interface and has an additional method expressionNode() that is used to differentiate it from statements.
type Expression interface {
	Node
	// expressionNode is a marker method to indicate that a struct is an Expression.
	expressionNode()
}

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
		_, _ = result.WriteString(stmt.String())
	}

	return result.String()
}
