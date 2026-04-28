package ast

import "github.com/OJOMB/donkey/internal/tokens"

// ExpressionKeyword represents an expression that is a keyword in the Donkey programming language, such as "break" or "continue".
type ExpressionKeyword struct {
	Token   tokens.Token
	Keyword string
}

func (ei *ExpressionKeyword) expressionNode()     {}
func (ei *ExpressionKeyword) TokenLexeme() string { return ei.Token.Lexeme }

// String returns the string representation of the ExpressionKeyword.
func (ek *ExpressionKeyword) String() string {
	return ek.Keyword
}
