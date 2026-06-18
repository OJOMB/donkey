package ast

import "github.com/OJOMB/donkey/internal/tokens"

// ExpressionIndex represents an index expression in the Donkey programming language.
// It consists of a left expression, which is any expression being indexed (e.g., myList/myMap/"myString"),
// and an index expression, which is the expression representing the index (e.g., "0" or "1+i").
type ExpressionIndex struct {
	// Token is the token associated with the index expression, which is typically the left bracket "[" token.
	Token tokens.Token
	// Left is the expression being indexed, such as a variable or another expression that evaluates to a list or string.
	Left Expression
	// Index is the expression representing the index, which can be any expression that evaluates to an integer or a string.
	Index Expression
}

func (ei *ExpressionIndex) expressionNode()     {}
func (ei *ExpressionIndex) TokenLexeme() string { return ei.Left.TokenLexeme() }

func (ei *ExpressionIndex) String() string {
	return ei.Left.String() + "[" + ei.Index.String() + "]"
}

func (ei *ExpressionIndex) Type() NodeType {
	return NodeTypeExpressionIndex
}
