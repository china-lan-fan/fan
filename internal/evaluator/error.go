package evaluator

import (
	"fan/internal/ast"
	"fan/internal/object"
)

type checkSignal struct {
	err *object.Error
}

func (c *checkSignal) Error() string { return c.err.Inspect() }

func evalCheckExpr(expr *ast.CheckExpr, env *Environment) (object.Object, error) {
	val, err := Eval(expr.Call, env)
	if err != nil {
		if sig, ok := err.(*checkSignal); ok {
			return nil, sig
		}
		return nil, err
	}
	if e, ok := errorFromValue(val); ok {
		return nil, &checkSignal{err: e}
	}
	return stripErrorTail(val), nil
}

func errorFromValue(val object.Object) (*object.Error, bool) {
	if val == nil {
		return nil, false
	}
	if e, ok := val.(*object.Error); ok {
		return e, true
	}
	if t, ok := val.(*object.Tuple); ok && len(t.Values) > 0 {
		last := t.Values[len(t.Values)-1]
		if e, ok := last.(*object.Error); ok {
			return e, true
		}
	}
	return nil, false
}

func stripErrorTail(val object.Object) object.Object {
	if t, ok := val.(*object.Tuple); ok && len(t.Values) > 1 {
		last := t.Values[len(t.Values)-1]
		if last.Kind() == object.KindError || last.Kind() == object.KindNil {
			values := t.Values[:len(t.Values)-1]
			if len(values) == 1 {
				return values[0]
			}
			return &object.Tuple{Values: values}
		}
	}
	return val
}

func evalTryStmt(stmt *ast.TryStmt, env *Environment) (object.Object, error) {
	tryEnv := NewEnclosedEnvironment(env)
	res, err := evalBlock(stmt.Body, tryEnv)
	if err == nil {
		return res, nil
	}
	if sig, ok := err.(*checkSignal); ok {
		return evalCatch(stmt, env, sig.err)
	}
	if e, ok := err.(*object.Error); ok {
		return evalCatch(stmt, env, e)
	}
	return nil, err
}

func evalCatch(stmt *ast.TryStmt, env *Environment, errObj *object.Error) (object.Object, error) {
	if stmt.Catch == nil {
		return object.Null, nil
	}
	catchEnv := NewEnclosedEnvironment(env)
	name := stmt.CatchName
	if name == "" {
		name = "错误"
	}
	if err := catchEnv.declare(name, errObj, false, ast.TypeAny); err != nil {
		return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
	}
	return evalBlock(stmt.Catch, catchEnv)
}

func unwrapCheckSignal(err error) (*object.Error, bool) {
	if sig, ok := err.(*checkSignal); ok {
		return sig.err, true
	}
	return nil, false
}

func zeroValueForType(dt ast.DeclType) object.Object {
	switch dt {
	case ast.TypeInt:
		return &object.Integer{}
	case ast.TypeFloat:
		return &object.Float{}
	case ast.TypeString:
		return &object.String{}
	case ast.TypeBool:
		return object.False
	case ast.TypeArray:
		return &object.Array{}
	case ast.TypeError:
		return object.Null
	}
	return object.Null
}
