package evaluator

import (
	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/pkg/logs"
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
	case *ast.ExpressionLiteralInteger:
		return &objects.Integer{Value: nt.Value}
	case *ast.ExpressionLiteralBoolean:
		return &objects.Boolean{Value: nt.Value}
	case *ast.ExpressionLiteralString:
		return &objects.String{Value: nt.Value}
	case *ast.ExpressionLiteralFunction:
		// function literals are not evaluated to a value until they are called, so we return a Nowt object for now
		return &objects.Nowt{}
	default:
		e.logger.Warn("unsupported AST node type", "type: %T", nt)
		return &objects.Nowt{}
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
