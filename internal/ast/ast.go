package ast

// Node is the base interface for all nodes in the AST.
// It has a method TokenLexeme() that returns the lexeme of the token associated with the node.
type Node interface {
	// TokenLexeme returns the lexeme of the token associated with the node.
	TokenLexeme() string
	// String returns a string representation of the node.
	String() string
}

// Statement represents a statement in the AST.
// It embeds the Node interface and has an additional method statementNode() that is used to differentiate it from expressions.
type Statement interface {
	Node
	// statementNode is a marker method to indicate that a struct is a Statement.
	statementNode()
}

// Expression represents an expression in the AST.
// It embeds the Node interface and has an additional method expressionNode() that is used to differentiate it from statements.
type Expression interface {
	Node
	// expressionNode is a marker method to indicate that a struct is an Expression.
	expressionNode()
}
