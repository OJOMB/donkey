package evaluator

import (
	"fmt"

	"github.com/OJOMB/donkey/internal/objects"
)

var (
	ErrInvalidForLoopInitializer = fmt.Errorf("invalid for loop initializer: expected a statement")
	ErrInvalidForLoopStep        = fmt.Errorf("invalid for loop step: expected a statement")
	ErrInvalidLoopCondition      = fmt.Errorf("invalid for loop condition: expected an expression")
	ErrInvalidLoopConditionType  = fmt.Errorf("invalid for loop condition: expected a boolean expression")
)

func newError(format string, args ...any) *objects.ErrorValue {
	return &objects.ErrorValue{Message: fmt.Sprintf(format, args...)}
}
