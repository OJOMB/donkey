package ast

type NodeType string

const (
	// NodeTypeProgram represents the node type for the root of the AST, which is a program.
	NodeTypeProgram NodeType = "Program"
	// NodeTypeStatementExpression represents the node type for an expression statement in the AST.
	NodeTypeStatementExpression NodeType = "StatementExpression"
	// NodeTypeStatementBind represents the node type for a bind statement in the AST.
	NodeTypeStatementBind NodeType = "StatementBind"
	// NodeTypeStatementRebind represents the node type for a rebind statement in the AST.
	NodeTypeStatementRebind NodeType = "StatementRebind"
	// NodeTypeStatementReturn represents the node type for a return statement in the AST.
	NodeTypeStatementReturn NodeType = "StatementReturn"
	// NodeTypeStatementFor represents the node type for a for statement in the AST.
	NodeTypeStatementFor NodeType = "StatementFor"
	// NodeTypeStatementWhile represents the node type for a while statement in the AST.
	NodeTypeStatementWhile NodeType = "StatementWhile"
	// NodeTypeStatementBlock represents the node type for a block statement in the AST.
	NodeTypeStatementBlock NodeType = "StatementBlock"
	// NodeTypeStatementIndexBind represents the node type for an index bind statement in the AST.
	NodeTypeStatementIndexBind NodeType = "StatementIndexBind"
	// NodeTypeStatementFunctionBind represents the node type for a function bind statement in the AST.
	NodeTypeStatementFunctionBind NodeType = "StatementFunctionBind"
	// NodeTypeExpressionCall represents the node type for a function call expression in the AST.
	NodeTypeExpressionCall NodeType = "ExpressionCall"
	// NodeTypeExpressionIdentifier represents the node type for an identifier expression in the AST.
	NodeTypeExpressionIdentifier NodeType = "ExpressionIdentifier"
	// NodeTypeExpressionIf represents the node type for an if expression in the AST.
	NodeTypeExpressionIf NodeType = "ExpressionIf"
	// NodeTypeExpressionInfix represents the node type for an infix expression in the AST.
	NodeTypeExpressionInfix NodeType = "ExpressionInfix"
	// NodeTypeExpressionPrefix represents the node type for a prefix expression in the AST.
	NodeTypeExpressionPrefix NodeType = "ExpressionPrefix"
	// NodeTypeExpressionKeyword represents the node type for a keyword expression in the AST.
	NodeTypeExpressionKeyword NodeType = "ExpressionKeyword"
	// NodeTypeExpressionLiteralNull represents the node type for a literal null expression in the AST.
	NodeTypeExpressionLiteralNull NodeType = "ExpressionLiteralNull"
	// NodeTypeExpressionLiteralBoolean represents the node type for a literal boolean expression in the AST.
	NodeTypeExpressionLiteralBoolean NodeType = "ExpressionLiteralBoolean"
	// NodeTypeExpressionLiteralFloat represents the node type for a literal float expression in the AST.
	NodeTypeExpressionLiteralFloat NodeType = "ExpressionLiteralFloat"
	// NodeTypeExpressionLiteralInteger represents the node type for a literal integer expression in the AST.
	NodeTypeExpressionLiteralInteger NodeType = "ExpressionLiteralInteger"
	// NodeTypeExpressionLiteralString represents the node type for a literal string expression in the AST.
	NodeTypeExpressionLiteralString NodeType = "ExpressionLiteralString"
	// NodeTypeExpressionLiteralList represents the node type for a literal list expression in the AST.
	NodeTypeExpressionLiteralList NodeType = "ExpressionLiteralList"
	// NodeTypeExpressionLiteralMap represents the node type for a literal map expression in the AST.
	NodeTypeExpressionLiteralMap NodeType = "ExpressionLiteralMap"
	// NodeTypeExpressionLiteralFunction represents the node type for a function literal expression in the AST.
	NodeTypeExpressionLiteralFunction NodeType = "ExpressionLiteralFunction"
	// NodeTypeExpressionList represents the node type for a list expression in the AST.
	NodeTypeExpressionList NodeType = "ExpressionList"
	// NodeTypeExpressionMap represents the node type for a map expression in the AST.
	NodeTypeExpressionMap NodeType = "ExpressionMap"
	// NodeTypeExpressionIndex represents the node type for an index expression in the AST.
	NodeTypeExpressionIndex NodeType = "ExpressionIndex"
)

// Node is the base interface for all nodes in the AST.
// It has a method TokenLexeme() that returns the lexeme of the token associated with the node.
type Node interface {
	// TokenLexeme returns the lexeme of the token associated with the node.
	TokenLexeme() string
	// String returns a string representation of the node.
	String() string
	// Type returns the type of the node as a NodeType.
	Type() NodeType
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
