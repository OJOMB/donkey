package ast

import "github.com/OJOMB/donkey/internal/tokens"

// ExpressionInfix represents an infix expression in the Donkey programming language, such as 5 + 5 or 10 - 2.
type ExpressionInfix struct {
	Token    tokens.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ei *ExpressionInfix) expressionNode()     {}
func (ei *ExpressionInfix) TokenLexeme() string { return ei.Token.Lexeme }

func (ei *ExpressionInfix) String() string {
	return "(" + ei.Left.String() + " " + ei.Operator + " " + ei.Right.String() + ")"
}

func (ei *ExpressionInfix) Type() NodeType {
	return NodeTypeExpressionInfix
}
