package evaluator

import (
	"fmt"
	"slices"
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
	stdLib builtinLib
	logger logs.Logger
}

// New creates a new Evaluator instance with the provided logger. If the logger is nil, a null logger will be used.
func New(l logs.Logger) *Evaluator {
	if l == nil {
		l = logs.NewNullLogger()
	}

	return &Evaluator{
		stdLib: builtins,
		logger: l.With("component", "evaluator")}
}

// Eval evaluates the given AST node and returns the resulting object.
func (e *Evaluator) Eval(node ast.Node, env *objects.Environment) objects.Object {
	switch nt := node.(type) {
	case *ast.Program:
		return e.evalStatements(nt, env)
	case *ast.StatementExpression:
		return e.Eval(nt.Expression, env)
	case *ast.StatementBlock:
		return e.evalStatementBlock(nt, env)
	case *ast.StatementBind:
		value := e.Eval(nt.Value, env)
		if value == nil {
			e.logger.Error("bind statement value evaluated to nil", "name", nt.Name.Value)
			return newError("bind statement value evaluated to nil: name:%s v:%v", nt.Name.Value, value)
		}

		return env.Bind(nt.Name.Value, value)
	case *ast.StatementIndexBind:
		return e.evalStatementIndexBind(nt, env)
	case *ast.StatementFunctionBind:
		return env.Bind(nt.Name.Value, &objects.Function{
			Parameters: nt.Value.Parameters,
			Body:       nt.Value.Body,
			Env:        env,
		})
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
	case *ast.StatementWhile:
		return e.evalStatementWhile(nt, env)
	case *ast.StatementFor:
		return e.evalStatementFor(nt, env)
	case *ast.ExpressionLiteralInteger, *ast.ExpressionLiteralBoolean, *ast.ExpressionLiteralString,
		*ast.ExpressionLiteralFunction, *ast.ExpressionLiteralList, *ast.ExpressionLiteralMap:
		return e.evalLiteral(nt, env)
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
		return e.evalExpressionIdentifier(nt, env)
	case *ast.ExpressionIndex:
		return e.evalExpressionIndex(nt, env)
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
	case *ast.ExpressionKeyword:
		return e.evalExpressionKeyword(nt, env)
	case *ast.ExpressionCall:
		functionObj := e.Eval(nt.Function, env)
		if functionObj == nil {
			e.logger.Error("call expression function evaluated to nil", "function", nt.Function.String())
			return newError("call expression function evaluated to nil: fn:%s", nt.Function.String())
		}

		args := make([]objects.Object, len(nt.Arguments))
		for i, arg := range nt.Arguments {
			evaluatedArg := e.Eval(arg, env)
			if evaluatedArg == nil {
				e.logger.Error("call expression argument evaluated to nil", "index", i, "argument", arg.String())
				return newError("call expression argument evaluated to nil: index:%d arg:%s", i, arg.String())
			}

			args[i] = evaluatedArg
		}

		return e.applyFunction(functionObj, args)
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

func (e *Evaluator) evalLiteral(node ast.Node, env *objects.Environment) objects.Object {
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
		return &objects.Function{
			Parameters: nt.Parameters,
			Body:       nt.Body,
			Env:        env,
		}
	case *ast.ExpressionLiteralList:
		elems := make([]objects.Object, len(nt.Elements))
		for i, elem := range nt.Elements {
			evaluatedElem := e.Eval(elem, env)
			if evaluatedElem == nil {
				e.logger.Error("list element evaluated to nil", "index", i, "element", elem.String())
				return newError("list element evaluated to nil: index:%d elem:%s", i, elem.String())
			}

			elems[i] = evaluatedElem
		}

		return &objects.List{Elements: elems}
	case *ast.ExpressionLiteralMap:
		pairs := make(map[objects.HashKey]objects.HashPair, len(nt.Pairs))
		for i, pair := range nt.Pairs {
			evaluatedKey := e.Eval(pair.Key, env)
			if evaluatedKey == nil {
				e.logger.Error("map key evaluated to nil", "index", i, "key", pair.Key.String())
				return newError("map key evaluated to nil: index:%d key:%s", i, pair.Key.String())
			}

			hashKey, ok := evaluatedKey.(objects.ObjectKey)
			if !ok {
				e.logger.Warn("map key is not a valid map key", "index", i, "key", pair.Key.String(), "type", evaluatedKey.Type())
				return newError("map key is not a valid map key: index:%d key:%s type:%s", i, pair.Key.String(), evaluatedKey.Type())
			}

			evaluatedValue := e.Eval(pair.Value, env)
			if evaluatedValue == nil {
				e.logger.Error("map value evaluated to nil", "index", i, "value", pair.Value.String())
				return newError("map value evaluated to nil: index:%d value:%s", i, pair.Value.String())
			}

			pairs[hashKey.HashKey()] = objects.HashPair{Key: evaluatedKey, Value: evaluatedValue}
		}

		return &objects.Map{Pairs: pairs}
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
	case tokens.TypePlus.String():
		return &objects.Integer{Value: left.Value + right.Value}
	case tokens.TypeMinus.String():
		return &objects.Integer{Value: left.Value - right.Value}
	case tokens.TypeAsterisk.String():
		return &objects.Integer{Value: left.Value * right.Value}
	case tokens.TypeForwardSlash.String():
		if right.Value == 0 {
			e.logger.Warn("division by zero")
			return newError("division by zero")
		}

		return &objects.Integer{Value: left.Value / right.Value}
	case tokens.TypePercent.String():
		if right.Value == 0 {
			e.logger.Warn("modulo by zero")
			return newError("modulo by zero")
		}

		return &objects.Integer{Value: left.Value % right.Value}
	case tokens.TypeExponent.String():
		if right.Value == 0 {
			return &objects.Integer{Value: 1}
		}

		result := 1
		for i := 0; i < right.Value; i++ {
			result *= left.Value
		}

		return &objects.Integer{Value: result}
	case tokens.TypeEQ.String():
		return &objects.Boolean{Value: left.Value == right.Value}
	case tokens.TypeNotEQ.String():
		return &objects.Boolean{Value: left.Value != right.Value}
	case tokens.TypeLT.String():
		return &objects.Boolean{Value: left.Value < right.Value}
	case tokens.TypeGT.String():
		return &objects.Boolean{Value: left.Value > right.Value}
	case tokens.TypeLTEQ.String():
		return &objects.Boolean{Value: left.Value <= right.Value}
	case tokens.TypeGTEQ.String():
		return &objects.Boolean{Value: left.Value >= right.Value}
	case tokens.TypeBitwiseAnd.String():
		return &objects.Integer{Value: left.Value & right.Value}
	case tokens.TypeBitwiseOr.String():
		return &objects.Integer{Value: left.Value | right.Value}
	case tokens.TypeBitwiseXor.String():
		return &objects.Integer{Value: left.Value ^ right.Value}
	case tokens.TypeBitwiseShiftLeft.String():
		return &objects.Integer{Value: left.Value << right.Value}
	case tokens.TypeBitwiseShiftRight.String():
		return &objects.Integer{Value: left.Value >> right.Value}
	default:
		e.logger.Warn("unsupported infix operator for integers", "operator", operator)
		return newError("unsupported infix operator for integers: %s", operator)
	}
}

func (e *Evaluator) evalExpressionInfixBoolean(operator string, left, right *objects.Boolean) objects.Object {
	switch operator {
	case tokens.TypeEQ.String():
		return &objects.Boolean{Value: left.Value == right.Value}
	case tokens.TypeNotEQ.String():
		return &objects.Boolean{Value: left.Value != right.Value}
	case tokens.TypeLogicalAnd.String():
		return &objects.Boolean{Value: left.Value && right.Value}
	case tokens.TypeLogicalOr.String():
		return &objects.Boolean{Value: left.Value || right.Value}
	default:
		e.logger.Warn("unsupported infix operator for booleans", "operator", operator)
		return newError("unsupported infix operator for booleans: %s", operator)
	}
}

func (e *Evaluator) evalExpressionInfixString(operator string, left, right *objects.String) objects.Object {
	switch operator {
	case tokens.TypePlus.String():
		return &objects.String{Value: left.Value + right.Value}
	case tokens.TypeMinus.String():
		// TODO: not overly convinced about this one
		return &objects.String{Value: strings.TrimSuffix(left.Value, right.Value)}
	case tokens.TypeEQ.String():
		return &objects.Boolean{Value: left.Value == right.Value}
	case tokens.TypeNotEQ.String():
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
	var blockEnv = objects.NewEnclosedEnvironment(env)

	var result objects.Object
	for i, stmt := range block.Statements {
		e.logger.Debug("evaluating statement in block", "index", i, "statement", stmt.String())
		result = e.Eval(stmt, blockEnv)

		if returnValue, ok := result.(*objects.ReturnValue); ok {
			return returnValue
		}

		if _, ok := result.(*objects.Continue); ok {
			return result
		}

		if _, ok := result.(*objects.Break); ok {
			return result
		}
	}

	return result
}

func (e *Evaluator) evalStatementWhile(node *ast.StatementWhile, env *objects.Environment) objects.Object {
	for {
		// while loops with no condition are treated as infinite loops, so we only evaluate the condition if it is present
		if node.Condition != nil {
			condition := e.Eval(node.Condition, env)
			if condition == nil {
				e.logger.Error("while loop condition evaluated to nil")
				return newError("while loop condition evaluated to nil")
			}

			if condition.Type() != objects.TypeBoolean {
				e.logger.Warn("while loop condition did not evaluate to a boolean", "type", condition.Type())
				return newError("while loop condition did not evaluate to a boolean: %s", condition.Type())
			}

			if !condition.(*objects.Boolean).Value {
				break
			}
		}

		result := e.evalStatementBlock(node.Body, env)
		if _, ok := result.(*objects.ReturnValue); ok {
			// bubble up return values from inside the loop body so that they can be handled by the caller
			return result
		}
	}

	// while loops do not produce a value, so we return Nowt to indicate the absence of a value
	return Nowt
}

func (e *Evaluator) evalStatementFor(node *ast.StatementFor, env *objects.Environment) objects.Object {
	// create a new environment for the loop body that is enclosed by the current environment
	// so that variables declared in the loop body do not leak out into the surrounding code
	loopEnv := objects.NewEnclosedEnvironment(env)

	// evaluate the initializer statement
	if err := e.EvalInitializer(node, loopEnv); err != nil {
		e.logger.Warn("failed to evaluate for loop initializer", "error", err)
		return newError("failed to evaluate for loop initializer: %s", err)
	}

	for {
		// evaluate the condition expression if it is present, otherwise treat the loop as infinite
		if truthy, _ := e.EvalCondition(node, loopEnv); !truthy {
			break
		}

		// evaluate the loop body
		result := e.evalStatementBlock(node.Body, loopEnv)
		if _, ok := result.(*objects.ReturnValue); ok {
			// bubble up return values from inside the loop body so that they can be handled by the caller
			return result
		}

		// if _, ok := result.(*objects.Continue); ok {
		// 	// continue statements skip the rest of the loop body and proceed to the next iteration
		// 	continue
		// }

		if _, ok := result.(*objects.Break); ok {
			// break statements exit the loop immediately
			break
		}

		// evaluate the step statement if it is present
		if node.Step != nil {
			if result := e.Eval(node.Step, loopEnv); result != nil {
				if _, ok := result.(*objects.ErrorValue); ok {
					return result
				}
			}
		}
	}

	// for loops do not produce a value, so we return Nowt to indicate the absence of a value
	return Nowt
}

func (e *Evaluator) evalExpressionKeyword(node *ast.ExpressionKeyword, env *objects.Environment) objects.Object {
	switch node.Keyword {
	case "continue":
		return &objects.Continue{}
	case "break":
		return &objects.Break{}
	default:
		e.logger.Warn("unsupported keyword", "keyword", node.Keyword)
		return newError("unsupported keyword: %s", node.Keyword)
	}
}

func (e *Evaluator) EvalInitializer(node *ast.StatementFor, env *objects.Environment) error {
	if node.Initializer == nil {
		return ErrInvalidForLoopInitializer
	}

	result := e.Eval(node.Initializer, env)
	if result == nil {
		return ErrInvalidForLoopInitializer
	}

	return nil
}

func (e *Evaluator) EvalCondition(node *ast.StatementFor, env *objects.Environment) (bool, error) {
	if node.Condition == nil {
		return true, ErrInvalidLoopCondition
	}

	conditionValue := e.Eval(node.Condition, env)
	conditionBool, ok := conditionValue.(*objects.Boolean)
	if !ok {
		return false, ErrInvalidLoopConditionType
	}

	return conditionBool.Value, nil
}

func (e *Evaluator) applyFunction(function objects.Object, args []objects.Object) objects.Object {
	switch fn := function.(type) {
	case *objects.BuiltinFunction:
		return fn.Implementation(args...)
	case *objects.Function:
		// create a new environment for the function execution that is enclosed by the function's defining environment
		env := objects.NewEnclosedEnvironment(fn.Env)

		// bind the function parameters to the argument values in the new environment
		for i, param := range fn.Parameters {
			if i >= len(args) {
				e.logger.Warn("not enough arguments provided for function call", "expected", len(fn.Parameters), "provided", len(args))
				return newError("not enough arguments provided for function call: expected %d, got %d", len(fn.Parameters), len(args))
			}

			env.Bind(param.Value, args[i])
		}

		// evaluate the function body in the new environment
		result := e.evalStatementBlock(fn.Body, env)
		if returnValue, ok := result.(*objects.ReturnValue); ok {
			return returnValue.Value
		}
	}

	// if the function body does not contain a return statement, we return Nowt to indicate the absence of a value
	return Nowt
}

func (e *Evaluator) evalExpressionIdentifier(node *ast.ExpressionIdentifier, env *objects.Environment) objects.Object {
	obj, ok := env.Get(node.Value)
	if ok {
		return obj
	}

	if builtin, ok := e.stdLib[node.Value]; ok {
		return builtin
	}

	e.logger.Warn("identifier not found in environment", "name", node.Value)

	return newError("identifier not found: %s", node.Value)
}

func (e *Evaluator) evalExpressionIndex(node *ast.ExpressionIndex, env *objects.Environment) objects.Object {
	left := e.Eval(node.Left, env)
	if left == nil {
		e.logger.Error("index expression left-hand side evaluated to nil", "expression", node.Left.String())
		return newError("index expression left-hand side evaluated to nil: expr:%s", node.Left.String())
	}

	switch left.Type() {
	case objects.TypeList:
		return e.evalExpressionListIndex(node, env)
	case objects.TypeMap:
		return e.evalExpressionMapIndex(node, env)
	default:
		e.logger.Warn("index expression left-hand side is not indexable", "type", left.Type())
		return newError("index expression left-hand side is not indexable: %s", left.Type())
	}
}

func (e *Evaluator) evalExpressionListIndex(node *ast.ExpressionIndex, env *objects.Environment) objects.Object {
	left := e.Eval(node.Left, env)
	if left == nil {
		e.logger.Error("index expression left-hand side evaluated to nil", "expression", node.Left.String())
		return newError("index expression left-hand side evaluated to nil: expr:%s", node.Left.String())
	}

	// left must be a list for indexing to be valid
	if left.Type() != objects.TypeList {
		e.logger.Warn("index expression left-hand side is not a list", "type", left.Type())
		return newError("index expression left-hand side is not a list: %s", left.Type())
	}

	index := e.Eval(node.Index, env)
	if index == nil {
		e.logger.Error("index expression index evaluated to nil", "expression", node.Index.String())
		return newError("index expression index evaluated to nil: expr:%s", node.Index.String())
	}

	if index.Type() != objects.TypeInteger {
		e.logger.Warn("index expression index is not an integer", "type", index.Type())
		return newError("index expression index is not an integer: %s", index.Type())
	}

	list := left.(*objects.List)
	idx := index.(*objects.Integer).Value

	if idx >= int(len(list.Elements)) {
		e.logger.Warn("index expression index out of bounds", "index", idx, "listLength", len(list.Elements))
		return newError("index expression index out of bounds: index:%d listLength:%d", idx, len(list.Elements))
	}

	// lets allow negative indices to count from the end of the list, so -1 is the last element, -2 is the second to last, and so on
	// if the negative index is out of bounds, we will return an error just like we do for positive indices that are out of bounds
	if idx < 0 {
		elemLen := int(len(list.Elements))
		if -idx > elemLen {
			e.logger.Warn("index expression negative index out of bounds", "index", idx, "listLength", elemLen)
			return newError("index expression negative index out of bounds: index:%d listLength:%d", idx, elemLen)
		}

		idx = int(len(list.Elements)) + idx
	}

	return list.Elements[idx]
}

func (e *Evaluator) evalExpressionMapIndex(node *ast.ExpressionIndex, env *objects.Environment) objects.Object {
	left := e.Eval(node.Left, env)
	if left == nil {
		e.logger.Error("index expression left-hand side evaluated to nil", "expression", node.Left.String())
		return newError("index expression left-hand side evaluated to nil: expr:%s", node.Left.String())
	}

	mapObj, ok := left.(*objects.Map)
	if !ok {
		e.logger.Warn("index expression left-hand side is not a map", "type", left.Type())
		return newError("index expression left-hand side is not a map: %s", left.Type())
	}

	index := e.Eval(node.Index, env)
	if index == nil {
		e.logger.Error("index expression index evaluated to nil", "expression", node.Index.String())
		return newError("index expression index evaluated to nil: expr:%s", node.Index.String())
	}

	indexKey, ok := index.(objects.ObjectKey)
	if !ok {
		e.logger.Warn("index expression index is not a valid map key", "type", index.Type())
		return newError("index expression index is not a valid map key: %s", index.Type())
	}

	valueObj, ok := mapObj.Get(indexKey)
	if !ok {
		e.logger.Warn("key not found in map", "key", index.Inspect())
		return newError("key not found in map: %s", index.Inspect())
	}

	return valueObj
}

func (e *Evaluator) getIndexableObjectAndIndices(node *ast.ExpressionIndex, env *objects.Environment) (objects.Object, []objects.Object, objects.Object) {
	var (
		idxExpr     = node
		idxExprLeft = node.Left
		obj         objects.Object
		idxs        []objects.Object
		ok          bool
	)

	for {
		switch l := idxExprLeft.(type) {
		case *ast.ExpressionIdentifier:
			obj, ok = env.Get(l.Value)
			if !ok {
				e.logger.Warn("identifier not found in environment", "name", l.Value)
				return nil, nil, newError("identifier not found: %s", l.Value)
			}
			if obj == nil {
				e.logger.Warn("identifier not found in environment", "name", l.Value)
				return nil, nil, newError("identifier not found: %s", l.Value)
			}

			idx := e.Eval(idxExpr.Index, env)
			idxs = append(idxs, idx)
		case *ast.ExpressionLiteralList:
			obj = e.Eval(l, env)
			if obj == nil {
				e.logger.Error("list literal evaluated to nil", "literal", l.String())
				return nil, nil, newError("list literal evaluated to nil: literal:%s", l.String())
			}

			idx := e.Eval(idxExpr.Index, env)
			idxs = append(idxs, idx)
		case *ast.ExpressionLiteralMap:
			obj = e.Eval(l, env)
			if obj == nil {
				e.logger.Error("map literal evaluated to nil", "literal", l.String())
				return nil, nil, newError("map literal evaluated to nil: literal:%s", l.String())
			}

			idx := e.Eval(idxExpr.Index, env)
			idxs = append(idxs, idx)
		case *ast.ExpressionIndex:
			idx := e.Eval(idxExpr.Index, env)
			idxs = append(idxs, idx)

			idxExprLeft = l.Left
			idxExpr = l
			continue
		default:
			e.logger.Warn("unsupported left expression type for index bind", "type", fmt.Sprintf("%T", l))
			return nil, nil, newError("unsupported left expression type for index bind: %T", l)
		}

		break
	}

	// check the type of the object to ensure it is indexable
	if obj.Type() != objects.TypeList && obj.Type() != objects.TypeMap {
		e.logger.Warn("object is not indexable", "type", obj.Type())
		return nil, nil, newError("object is not indexable: %s", obj.Type())
	}

	slices.Reverse(idxs)

	return obj, idxs, nil
}

func (e *Evaluator) evalStatementIndexBind(node *ast.StatementIndexBind, env *objects.Environment) objects.Object {
	obj, idxs, err := e.getIndexableObjectAndIndices(node.Left, env)
	if err != nil {
		return err
	}

	// now we have the object and the indices, we can set the value at the specified index
	right := e.Eval(node.Right, env)
	if right == nil {
		e.logger.Error("index bind right-hand side evaluated to nil", "expression", node.Right.String())
		return newError("index bind right-hand side evaluated to nil: expr:%s", node.Right.String())
	}

	objIndexable, ok := obj.(objects.Indexable)
	if !ok {
		e.logger.Warn("object is not indexable", "type", obj.Type())
		return newError("object is not indexable: %s", obj.Type())
	}

	return e.setValueAtIndex(objIndexable, idxs, right)
}

func (e *Evaluator) setValueAtIndex(obj objects.Indexable, idxs []objects.Object, value objects.Object) objects.Object {
	if err := obj.SetMultiDimensional(idxs, value); err != nil {
		e.logger.Warn("failed to set value at index", "error", err)
		return newError("failed to set value at index: %s", err)
	}

	return obj
}
