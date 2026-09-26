package evaluator

import (
	"fmt"

	"fan/internal/ast"
	"fan/internal/object"
)

const KindFunction object.Kind = "函数"

type Function struct {
	Params      []ast.Parameter
	Body        *ast.BlockStmt
	Env         *Environment
	Name        string
	ReturnTypes []ast.DeclType
}

func (f *Function) Kind() object.Kind { return KindFunction }
func (f *Function) Inspect() string {
	if f.Name != "" {
		return fmt.Sprintf("<函数 %s>", f.Name)
	}
	return "<函数>"
}

type returnSignal struct {
	values []object.Object
}

func (r *returnSignal) Error() string {
	return "return signal"
}

func evalFunctionLiteral(node *ast.FunctionLiteral, env *Environment) (object.Object, error) {
	return &Function{
		Params:      node.Params,
		Body:        node.Body,
		Env:         env,
		Name:        node.Name,
		ReturnTypes: node.ReturnTypes,
	}, nil
}

func evalMemberExpression(node *ast.MemberExpr, env *Environment) (object.Object, error) {
	obj, err := Eval(node.Object, env)
	if err != nil {
		return nil, err
	}
	return evalFieldAccess(node, env, obj)
}

func applyFunction(fn *Function, args []object.Object, pos ast.Position) (object.Object, error) {
	if len(args) != len(fn.Params) {
		return nil, &EvalError{
			Pos:    pos,
			Reason: fmt.Sprintf("函数参数数量不符：需要 %d 个，实际 %d 个", len(fn.Params), len(args)),
		}
	}
	env := NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Params {
		if err := checkDeclType(param.Type, args[i]); err != nil {
			return nil, &EvalError{Pos: pos, Reason: fmt.Sprintf("参数 %s 类型不匹配：%s", param.Name, err.Error())}
		}
		if err := env.declare(param.Name, args[i], false, param.Type); err != nil {
			return nil, &EvalError{Pos: pos, Reason: err.Error()}
		}
	}
	result, err := evalBlock(fn.Body, env)
	var values []object.Object
	if err != nil {
		if sig, ok := err.(*returnSignal); ok {
			values = sig.values
		} else if cs, ok := err.(*checkSignal); ok {
			values, err = fn.checkReturn(cs.err)
			if err != nil {
				return nil, err
			}
			return coerceReturnValues(pos, fn.ReturnTypes, values)
		} else {
			return nil, err
		}
	} else {
		if result == nil {
			result = object.Null
		}
		values = []object.Object{result}
	}
	return coerceReturnValues(pos, fn.ReturnTypes, values)
}

func (f *Function) checkReturn(cause *object.Error) ([]object.Object, error) {
	if len(f.ReturnTypes) == 0 {
		return nil, cause
	}
	last := f.ReturnTypes[len(f.ReturnTypes)-1]
	if last != ast.TypeError {
		return nil, cause
	}
	values := make([]object.Object, len(f.ReturnTypes))
	for i, dt := range f.ReturnTypes[:len(f.ReturnTypes)-1] {
		values[i] = zeroValueForType(dt)
	}
	values[len(values)-1] = cause
	return values, nil
}

func coerceReturnValues(pos ast.Position, types []ast.DeclType, values []object.Object) (object.Object, error) {
	if len(types) == 0 {
		if len(values) <= 1 {
			return values[0], nil
		}
		return &object.Tuple{Values: values}, nil
	}
	if len(types) != len(values) {
		return nil, &EvalError{Pos: pos, Reason: fmt.Sprintf("返回值数量不符：声明 %d 个，实际 %d 个", len(types), len(values))}
	}
	for i, dt := range types {
		if values[i] == nil {
			values[i] = object.Null
		}
		if dt == ast.TypeError && values[i].Kind() == object.KindNil {
			continue
		}
		if err := checkDeclType(dt, values[i]); err != nil {
			return nil, &EvalError{Pos: pos, Reason: fmt.Sprintf("第 %d 个返回值%s", i+1, err.Error())}
		}
	}
	if len(values) == 1 {
		return values[0], nil
	}
	return &object.Tuple{Values: values}, nil
}
