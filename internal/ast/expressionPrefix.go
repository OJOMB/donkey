package ast

import "github.com/OJOMB/monkey/internal/tokens"

type ExpressionPrefix struct {
	Token    tokens.Token
	Operator string
	Right    Expression
}

func (ep *ExpressionPrefix) expressionNode()     {}
func (ep *ExpressionPrefix) TokenLexeme() string { return ep.Token.Lexeme }

func (ep *ExpressionPrefix) String() string {
	return "(" + ep.Operator + ep.Right.String() + ")"
}
