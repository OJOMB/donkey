package evaluator

import (
	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/pkg/logs"
)

type Evaluator struct {
	logger logs.Logger
}

func NewEvaluator(l logs.Logger) *Evaluator {
	if l == nil {
		l = logs.NewNullLogger()
	}

	return &Evaluator{logger: l.With("component", "evaluator")}
}

// Eval evaluates the given AST node and returns the resulting object.
func (e *Evaluator) Eval(node ast.Node) objects.Object {
	return nil
}
