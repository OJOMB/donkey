package ast

import "github.com/OJOMB/monkey/internal/tokens"

// ExpressionPrefix represents a prefix expression in the Monkey programming language.
// For example, in the expression "!true", the "!" is the prefix operator and "true" is the RHS expression.
// Similarly, in the expression "-10", the "-" is the prefix operator and "10" is the RHS expression.
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
