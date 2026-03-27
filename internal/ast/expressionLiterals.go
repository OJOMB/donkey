package ast

import "github.com/OJOMB/monkey/internal/tokens"

// ExpressionLiteralInteger represents an integer literal expression in the Monkey programming language, such as 5 or 10.
// For example, in the expression "let x = 5;", the "5" is an integer literal expression that represents the value being assigned to the variable "x" in the let statement.
type ExpressionLiteralInteger struct {
	Token tokens.Token
	Value int
}

func (li *ExpressionLiteralInteger) expressionNode()     {}
func (li *ExpressionLiteralInteger) TokenLexeme() string { return li.Token.Lexeme }

func (li *ExpressionLiteralInteger) String() string {
	return li.Token.Lexeme
}

type ExpressionLiteralString struct {
	Token tokens.Token
	Value string
}

func (ls *ExpressionLiteralString) expressionNode()     {}
func (ls *ExpressionLiteralString) TokenLexeme() string { return ls.Token.Lexeme }

func (ls *ExpressionLiteralString) String() string {
	return ls.Token.Lexeme
}

type ExpressionLiteralBoolean struct {
	Token tokens.Token
	Value bool
}

func (lb *ExpressionLiteralBoolean) expressionNode()     {}
func (lb *ExpressionLiteralBoolean) TokenLexeme() string { return lb.Token.Lexeme }

func (lb *ExpressionLiteralBoolean) String() string {
	return lb.Token.Lexeme
}
