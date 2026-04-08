package evaluator

import (
	"go/ast"

	"github.com/OJOMB/donkey/internal/objects"
)

type Evaluator struct {
}

func NewEvaluator() *Evaluator {
	return &Evaluator{}
}

func (e *Evaluator) Eval(node ast.Node) objects.Object {
	return nil
}
