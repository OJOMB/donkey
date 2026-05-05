package ast

import (
	"github.com/OJOMB/donkey/internal/tokens"
)

// StatementFunctionBind represents a function binding statement in the AST. It contains a token for the binding, an identifier for the function name, and an expression for the function being bound.
// For example, in the function binding statement "fn add(x, y) { return x + y; }",
// the token would be the "fn" token,
// the Name would be an ExpressionIdentifier representing "add"
// the Value would be an ExpressionLiteralFunction representing the function definition.
type StatementFunctionBind struct {
	// Token is the token associated with the function binding statement.
	Token tokens.Token
	// Name is the identifier for the function being bound.
	Name *ExpressionIdentifier
	// Value is the expression representing the function being bound.
	Value *ExpressionLiteralFunction
}

// statementNode is a marker method to indicate that StatementFunctionBind is a statement node in the AST.
func (s *StatementFunctionBind) statementNode() {}

// TokenLexeme returns the lexeme of the token associated with the function binding statement, which is typically "fn" for function declarations.
func (s *StatementFunctionBind) TokenLexeme() string {
	return s.Token.Lexeme
}

func (s *StatementFunctionBind) String() string {
	return s.TokenLexeme() + " " + s.Name.String() + s.Value.String()
}
