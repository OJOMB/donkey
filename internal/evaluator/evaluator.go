package evaluator

import (
	"fmt"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/internal/tokens"
	"github.com/OJOMB/donkey/pkg/logs"
)

var (
	// Nowt is the singleton Nowt object that represents the absence of a value in the Donkey programming language.
	Nowt  = &objects.Nowt{}
	True  = &objects.Boolean{Value: true}
	False = &objects.Boolean{Value: false}
)

type Evaluator struct {
	logger logs.Logger
}

func New(l logs.Logger) *Evaluator {
	if l == nil {
		l = logs.NewNullLogger()
	}

	return &Evaluator{logger: l.With("component", "evaluator")}
}

// Eval evaluates the given AST node and returns the resulting object.
func (e *Evaluator) Eval(node ast.Node) objects.Object {
	switch nt := node.(type) {
	case *ast.Program:
		return e.evalStatements(nt)
	case *ast.StatementExpression:
		return e.Eval(nt.Expression)
	case *ast.ExpressionPrefix:
		right := e.Eval(nt.Right)
		if right == nil {
			e.logger.Error("prefix operator right-hand side evaluated to nil", "operator", nt.Token.Lexeme)
			return Nowt
		}

		switch nt.Token.Type {
		case tokens.TypeBang:
			return e.evalExpressionPrefixBang(right)
		case tokens.TypeMinus:
			return evalExpressionPrefixMinus(right, e)
		default:
			e.logger.Error("unsupported prefix operator", "operator", nt.Token.Lexeme)
			return Nowt
		}
	case *ast.ExpressionLiteralInteger:
		return &objects.Integer{Value: nt.Value}
	case *ast.ExpressionLiteralBoolean:
		if nt.Value {
			return True
		}
		return False
	case *ast.ExpressionLiteralString:
		return &objects.String{Value: nt.Value}
	case *ast.ExpressionLiteralFunction:
		// function literals are not evaluated to a value until they are called, so we return a Nowt object for now
		return Nowt
	default:
		e.logger.Warn("unsupported AST node type", "type", fmt.Sprintf("%T", nt))
		return Nowt
	}
}

func (e *Evaluator) evalStatements(program *ast.Program) objects.Object {
	var result objects.Object
	for i, stmt := range program.Statements {
		e.logger.Debug("evaluating statement", "index", i, "statement", stmt.String())
		result = e.Eval(stmt)
	}

	return result
}

func (e *Evaluator) evalExpressionPrefixBang(right objects.Object) objects.Object {
	if right.Type() != objects.TypeBoolean {
		e.logger.Warn("unsupported operand type for ! operator", "type", right.Type())
		return Nowt
	}

	switch right {
	case True:
		return False
	case False:
		return True
	default:
		return Nowt
	}
}

func evalExpressionPrefixMinus(right objects.Object, e *Evaluator) objects.Object {
	if right.Type() != objects.TypeInteger {
		e.logger.Warn("unsupported operand type for - operator", "type", right.Type())
		return Nowt
	}

	value := right.(*objects.Integer).Value
	return &objects.Integer{Value: -value}
}
