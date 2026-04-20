package ast

import (
	"strings"

	"github.com/OJOMB/donkey/internal/tokens"
)

// ExpressionLiteralInteger represents an integer literal expression in the Donkey programming language, such as 5 or 10.
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

// ExpressionLiteralString represents a string literal expression in the Donkey programming language, such as "hello" or "world".
// For example, in the expression "let greeting = "hello";", the ""hello"" is a string literal expression that represents the value being assigned to the variable "greeting" in the let statement.
type ExpressionLiteralString struct {
	Token tokens.Token
	Value string
}

func (ls *ExpressionLiteralString) expressionNode()     {}
func (ls *ExpressionLiteralString) TokenLexeme() string { return ls.Token.Lexeme }

func (ls *ExpressionLiteralString) String() string {
	return ls.Token.Lexeme
}

// ExpressionLiteralBoolean represents a boolean literal expression in the Donkey programming language, such as true or false.
// For example, in the expression "let isValid = true;", the "true" is a boolean literal expression that represents the value being assigned to the variable "isValid" in the let statement.
type ExpressionLiteralBoolean struct {
	Token tokens.Token
	// Value is the boolean value of the expression, which can be either true or false.
	Value bool
}

func (lb *ExpressionLiteralBoolean) expressionNode()     {}
func (lb *ExpressionLiteralBoolean) TokenLexeme() string { return lb.Token.Lexeme }

func (lb *ExpressionLiteralBoolean) String() string {
	return lb.Token.Lexeme
}

// ExpressionLiteralFunction represents a function literal expression in the Donkey programming language, such as "fn(x) { x + 1 }".
// For example, in the expression "let add = fn(x) { x + 1 };"
// the "fn(x) { x + 1 }" is a function literal expression that represents the value being assigned to the variable "add" in the let statement.
// not to be confused with ExpressionCall, which represents a function call expression like "add(5)" where "add" is the function being called and "5" is the argument passed to the function.
// fn(<parameters>) { <body> }
type ExpressionLiteralFunction struct {
	// Token is the token associated with the function literal, which is the "fn" keyword.
	Token tokens.Token
	// Parameters is a slice of pointers to ExpressionIdentifier nodes representing the parameters of the function.
	Parameters []*ExpressionIdentifier
	// Body is a pointer to a StatementBlock node representing the body of the function, which contains the statements that will be executed when the function is called.
	Body *StatementBlock
}

func (lf *ExpressionLiteralFunction) expressionNode()     {}
func (lf *ExpressionLiteralFunction) TokenLexeme() string { return lf.Token.Lexeme }

func (lf *ExpressionLiteralFunction) String() string {
	var out = strings.Builder{}
	_, _ = out.WriteString(lf.Token.Lexeme)
	_, _ = out.WriteString("(")
	for i, param := range lf.Parameters {
		if i > 0 {
			_, _ = out.WriteString(", ")
		}

		_, _ = out.WriteString(param.String())
	}
	_, _ = out.WriteString(") ")
	_, _ = out.WriteString(lf.Body.String())

	return out.String()
}
