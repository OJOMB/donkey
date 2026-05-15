package ast

import (
	"strings"

	"github.com/OJOMB/donkey/internal/tokens"
)

// ExpressionLiteralInteger represents an integer literal expression in the Donkey programming language, such as 5 or 10.
// For example, in the expression "var x = 5;", the "5" is an integer literal expression that represents the value being assigned to the variable "x" in the var statement.
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
// For example, in the expression "var greeting = "hello";", the ""hello"" is a string literal expression that represents the value being assigned to the variable "greeting" in the var statement.
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
// For example, in the expression "var isValid = true;", the "true" is a boolean literal expression that represents the value being assigned to the variable "isValid" in the var statement.
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
// For example, in the expression "var add = fn(x) { x + 1 };"
// the "fn(x) { x + 1 }" is a function literal expression that represents the value being assigned to the variable "add" in the var statement.
// not to be confused with ExpressionCall, which represents a function call expression like "add(5)" where "add" is the function being called and "5" is the argument passed to the function.
// fn(<parameters>) { <body> }
// essentially an anonymous function that can be assigned to a variable or passed as an argument to another function, allowing for higher-order programming and functional programming paradigms in the Donkey language.
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

// ExpressionLiteralList represents a list literal expression in the Donkey programming language, such as "[1, 2, 3]" or "["foo", "bar", "baz"]".
// For example, in the expression "var myList = [1, 2, 3];", the "[1, 2, 3]" is a list literal expression that represents the value being assigned to the variable "myList" in the var statement.
type ExpressionLiteralList struct {
	Token    tokens.Token
	Elements []Expression
}

func (ll *ExpressionLiteralList) expressionNode()     {}
func (ll *ExpressionLiteralList) TokenLexeme() string { return ll.Token.Lexeme }

func (ll *ExpressionLiteralList) String() string {
	var out = strings.Builder{}
	_, _ = out.WriteString("[")
	for i, elem := range ll.Elements {
		if i > 0 {
			_, _ = out.WriteString(", ")
		}

		_, _ = out.WriteString(elem.String())
	}

	_, _ = out.WriteString("]")

	return out.String()
}

// ExpressionLiteralMap represents a map literal expression in the Donkey programming language, such as "{"key": "value", "foo": "bar"}".
// For example, in the expression "var myMap = {"key": "value", "foo": "bar"};", the "{"key": "value", "foo": "bar"}" is a map literal expression that represents the value being assigned to the variable "myMap" in the var statement.
type ExpressionLiteralMap struct {
	Token tokens.Token
	// Pairs is a map where the keys are expressions representing the keys of the map, and the values are expressions representing the corresponding values in the map.
	Pairs map[Expression]Expression
}

func (lm *ExpressionLiteralMap) expressionNode()     {}
func (lm *ExpressionLiteralMap) TokenLexeme() string { return lm.Token.Lexeme }

func (lm *ExpressionLiteralMap) String() string {
	var out = strings.Builder{}
	_, _ = out.WriteString("{")
	i := 0
	for key, value := range lm.Pairs {
		if i > 0 {
			_, _ = out.WriteString(", ")
		}
		_, _ = out.WriteString(key.String())
		_, _ = out.WriteString(": ")
		_, _ = out.WriteString(value.String())
		i++
	}
	_, _ = out.WriteString("}")
	return out.String()
}
