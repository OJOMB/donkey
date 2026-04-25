package evaluator

import (
	"fmt"
	"strings"

	"github.com/OJOMB/donkey/internal/ast"
	"github.com/OJOMB/donkey/internal/objects"
	"github.com/OJOMB/donkey/internal/tokens"
	"github.com/OJOMB/donkey/pkg/logs"
)

var (
	// Nowt is the singleton Nowt object that represents the absence of a value in the Donkey programming language.
	Nowt = &objects.Nowt{}
	// True is the singleton Boolean object that represents the boolean value true in the Donkey programming language.
	True = &objects.Boolean{Value: true}
	// False is the singleton Boolean object that represents the boolean value false in the Donkey programming language.
	False = &objects.Boolean{Value: false}
)

// Evaluator is responsible for evaluating input AST nodes and producing the corresponding objects in the Donkey programming language.
type Evaluator struct {
	logger logs.Logger
}

// New creates a new Evaluator instance with the provided logger. If the logger is nil, a null logger will be used.
func New(l logs.Logger) *Evaluator {
	if l == nil {
		l = logs.NewNullLogger()
	}

	return &Evaluator{logger: l.With("component", "evaluator")}

}

// Eval evaluates the given AST node and returns the resulting object.
func (e *Evaluator) Eval(node ast.Node, env *objects.Environment) objects.Object {
	switch nt := node.(type) {
	case *ast.ExpressionLiteralInteger, *ast.ExpressionLiteralBoolean, *ast.ExpressionLiteralString, *ast.ExpressionLiteralFunction:
		return e.evalLiteral(nt)
	case *ast.Program:
		return e.evalStatements(nt, env)
	case *ast.StatementExpression:
		return e.Eval(nt.Expression, env)
	case *ast.ExpressionPrefix:
		right := e.Eval(nt.Right, env)
		if right == nil {
			e.logger.Error("prefix operator right-hand side evaluated to nil", "operator", nt.Token.Lexeme)
			return newError("prefix operator right-hand side evaluated to nil: op:%s r:%v", nt.Token.Lexeme, right)
		}

		switch nt.Token.Type {
		case tokens.TypeBang:
			return e.evalExpressionPrefixBang(right)
		case tokens.TypeMinus:
			return e.evalExpressionPrefixMinus(right)
		default:
			e.logger.Error("unsupported prefix operator", "operator", nt.Token.Lexeme)
			return newError("unsupported prefix operator: %s", nt.Token.Lexeme)
		}
	case *ast.ExpressionIdentifier:
		obj, ok := env.Get(nt.Value)
		if !ok {
			e.logger.Warn("identifier not found in environment", "name", nt.Value)
			return newError("identifier not found: %s", nt.Value)
		}

		return obj
	case *ast.ExpressionInfix:
		l := e.Eval(nt.Left, env)
		if l == nil {
			e.logger.Error("infix operator left-hand side evaluated to nil", "operator", nt.Token.Lexeme)
			return newError("infix operator left-hand side evaluated to nil: op:%s l:%v", nt.Token.Lexeme, l)
		}

		r := e.Eval(nt.Right, env)
		if r == nil {
			e.logger.Error("infix operator right-hand side evaluated to nil", "operator", nt.Token.Lexeme)
			return newError("infix operator right-hand side evaluated to nil: op:%s r:%v", nt.Token.Lexeme, r)
		}

		return e.evalExpressionInfix(nt.Operator, l, r)
	case *ast.ExpressionIf:
		return e.evalExpressionIf(nt, env)
	case *ast.StatementBlock:
		return e.evalStatementBlock(nt, env)
	case *ast.StatementBind:
		value := e.Eval(nt.Value, env)
		if value == nil {
			e.logger.Error("bind statement value evaluated to nil", "name", nt.Name.Value)
			return newError("bind statement value evaluated to nil: name:%s v:%v", nt.Name.Value, value)
		}

		return env.Bind(nt.Name.Value, value)
	case *ast.StatementRebind:
		value := e.Eval(nt.Value, env)
		if value == nil {
			e.logger.Error("rebind statement value evaluated to nil", "name", nt.Name.Value)
			return newError("rebind statement value evaluated to nil: name:%s v:%v", nt.Name.Value, value)
		}

		if _, err := env.Set(nt.Name.Value, value); err != nil {
			e.logger.Warn("failed to rebind variable in environment", "name", nt.Name.Value, "error", err)
			return newError("uninitialized variable: %s", nt.Name.Value)
		}

		return value
	case *ast.StatementReturn:
		value := e.Eval(nt.Value, env)
		if value == nil {
			e.logger.Error("return statement value evaluated to nil")
			return Nowt
		}

		return &objects.ReturnValue{Value: value}
	default:
		e.logger.Error("unsupported AST node type", "type", fmt.Sprintf("%T", nt))
		return newError("unsupported AST node type: %T", nt)
	}
}

func (e *Evaluator) evalStatements(program *ast.Program, env *objects.Environment) objects.Object {
	var result objects.Object
	for i, stmt := range program.Statements {
		e.logger.Debug("evaluating statement", "index", i, "statement", stmt.String())
		result = e.Eval(stmt, env)

		if returnValue, ok := result.(*objects.ReturnValue); ok {
			return returnValue.Value
		}

		if _, ok := result.(*objects.ErrorValue); ok {
			return result
		}
	}

	return result
}

func (e *Evaluator) evalExpressionPrefixBang(right objects.Object) objects.Object {
	if right.Type() != objects.TypeBoolean {
		e.logger.Warn("unsupported operand type for ! operator", "type", right.Type())
		return newError("unsupported operand type for ! operator: %s", right.Type())
	}

	switch right {
	case True:
		return False
	case False:
		return True
	default:
		return newError("unsupported boolean value: %s", right.Inspect())
	}
}

func (e *Evaluator) evalExpressionPrefixMinus(right objects.Object) objects.Object {
	if right.Type() != objects.TypeInteger {
		e.logger.Warn("unsupported operand type for - operator", "type", right.Type())
		return newError("unsupported operand type for - operator: %s", right.Type())
	}

	value := right.(*objects.Integer).Value
	return &objects.Integer{Value: -value}
}

func (e *Evaluator) evalLiteral(node ast.Node) objects.Object {
	switch nt := node.(type) {
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
		e.logger.Warn("unsupported literal type", "type", fmt.Sprintf("%T", nt))
		return Nowt
	}
}

func (e *Evaluator) evalExpressionInfix(operator string, left, right objects.Object) objects.Object {
	if left == nil || right == nil {
		e.logger.Error("infix operator operands evaluated to nil", "operator", operator, "leftNil", left == nil, "rightNil", right == nil)
		return newError("infix operator operands evaluated to nil: op:%s l:%v r:%v", operator, left, right)
	}

	if left.Type() != right.Type() {
		e.logger.Warn("type mismatch for infix operator", "operator", operator, "leftType", left.Type(), "rightType", right.Type())
		return newError("type mismatch for infix operator: %s %s %s", left.Type(), operator, right.Type())
	}

	switch {
	case left.Type() == objects.TypeInteger && right.Type() == objects.TypeInteger:
		return e.evalExpressionInfixInteger(operator, left.(*objects.Integer), right.(*objects.Integer))
	case left.Type() == objects.TypeBoolean && right.Type() == objects.TypeBoolean:
		return e.evalExpressionInfixBoolean(operator, left.(*objects.Boolean), right.(*objects.Boolean))
	case left.Type() == objects.TypeString && right.Type() == objects.TypeString:
		return e.evalExpressionInfixString(operator, left.(*objects.String), right.(*objects.String))
	default:
		e.logger.Warn("unsupported operand types for infix operator", "operator", operator, "leftType", left.Type(), "rightType", right.Type())
		return newError("unsupported operand types for infix operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func (e *Evaluator) evalExpressionInfixInteger(operator string, left, right *objects.Integer) objects.Object {
	switch operator {
	case "+":
		return &objects.Integer{Value: left.Value + right.Value}
	case "-":
		return &objects.Integer{Value: left.Value - right.Value}
	case "*":
		return &objects.Integer{Value: left.Value * right.Value}
	case "/":
		if right.Value == 0 {
			e.logger.Warn("division by zero")
			return newError("division by zero")
		}

		return &objects.Integer{Value: left.Value / right.Value}
	case "%":
		if right.Value == 0 {
			e.logger.Warn("modulo by zero")
			return newError("modulo by zero")
		}

		return &objects.Integer{Value: left.Value % right.Value}
	case "^":
		if right.Value == 0 {
			return &objects.Integer{Value: 1}
		}

		result := 1
		for i := 0; i < right.Value; i++ {
			result *= left.Value
		}

		return &objects.Integer{Value: result}
	case "==":
		return &objects.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &objects.Boolean{Value: left.Value != right.Value}
	case "<":
		return &objects.Boolean{Value: left.Value < right.Value}
	case ">":
		return &objects.Boolean{Value: left.Value > right.Value}
	case "<=":
		return &objects.Boolean{Value: left.Value <= right.Value}
	case ">=":
		return &objects.Boolean{Value: left.Value >= right.Value}
	case "&":
		return &objects.Integer{Value: left.Value & right.Value}
	case "|":
		return &objects.Integer{Value: left.Value | right.Value}
	default:
		e.logger.Warn("unsupported infix operator for integers", "operator", operator)
		return newError("unsupported infix operator for integers: %s", operator)
	}
}

func (e *Evaluator) evalExpressionInfixBoolean(operator string, left, right *objects.Boolean) objects.Object {
	switch operator {
	case "==":
		return &objects.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &objects.Boolean{Value: left.Value != right.Value}
	case "&&":
		return &objects.Boolean{Value: left.Value && right.Value}
	case "||":
		return &objects.Boolean{Value: left.Value || right.Value}
	default:
		e.logger.Warn("unsupported infix operator for booleans", "operator", operator)
		return newError("unsupported infix operator for booleans: %s", operator)
	}
}

func (e *Evaluator) evalExpressionInfixString(operator string, left, right *objects.String) objects.Object {
	switch operator {
	case "+":
		return &objects.String{Value: left.Value + right.Value}
	case "-":
		// TODO: not overly convinced about this one
		return &objects.String{Value: strings.TrimSuffix(left.Value, right.Value)}
	case "==":
		return &objects.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &objects.Boolean{Value: left.Value != right.Value}
	default:
		e.logger.Warn("unsupported infix operator for strings", "operator", operator)
		return newError("unsupported infix operator for strings: %s", operator)
	}
}

func (e *Evaluator) evalExpressionIf(node *ast.ExpressionIf, env *objects.Environment) objects.Object {
	for _, branch := range node.Branches {
		condition := e.Eval(branch.Condition, env)
		if condition == nil {
			e.logger.Error("if condition evaluated to nil")
			return newError("if condition evaluated to nil")
		}

		if condition.Type() != objects.TypeBoolean {
			e.logger.Warn("if condition did not evaluate to a boolean", "type", condition.Type())
			return newError("if condition did not evaluate to a boolean: %s", condition.Type())
		}

		if condition.(*objects.Boolean).Value {
			return e.evalStatementBlock(branch.Consequence, env)
		}
	}

	if node.Alternative != nil {
		return e.evalStatementBlock(node.Alternative, env)
	}

	return Nowt
}

func (e *Evaluator) evalStatementBlock(block *ast.StatementBlock, env *objects.Environment) objects.Object {
	var result objects.Object
	for i, stmt := range block.Statements {
		e.logger.Debug("evaluating statement in block", "index", i, "statement", stmt.String())
		result = e.Eval(stmt, env)

		if returnValue, ok := result.(*objects.ReturnValue); ok {
			return returnValue
		}
	}

	return result
}
