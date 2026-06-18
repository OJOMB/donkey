package ast

import (
	"strings"

	"github.com/OJOMB/donkey/internal/tokens"
)

// ExpressionCall represents a function call expression in the AST. It contains a token for the call, an expression for the function being called, and a slice of expressions for the arguments passed to the function.
type ExpressionCall struct {
	// Token is the token associated with the call expression.
	Token tokens.Token
	// Function is an expression representing the function being called. It can be an identifier (for calling a named function) or a function literal (for calling an anonymous function).
	Function Expression
	// Arguments is a slice of expressions representing the arguments passed to the function in the call. Each argument can be any valid expression in the Donkey programming language.
	Arguments []Expression
}

func (ec *ExpressionCall) expressionNode()     {}
func (ec *ExpressionCall) TokenLexeme() string { return ec.Token.Lexeme }

func (ec *ExpressionCall) String() string {
	var out = strings.Builder{}
	_, _ = out.WriteString(ec.Function.String())
	_, _ = out.WriteString("(")
	for i, arg := range ec.Arguments {
		if i > 0 {
			_, _ = out.WriteString(", ")
		}
		_, _ = out.WriteString(arg.String())
	}

	_, _ = out.WriteString(")")
	return out.String()
}

func (ec *ExpressionCall) Type() NodeType {
	return NodeTypeExpressionCall
}
