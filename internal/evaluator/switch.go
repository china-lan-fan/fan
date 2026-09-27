package evaluator

import (
	"fmt"

	"fan/internal/ast"
	"fan/internal/object"
)

func evalSwitchStmt(stmt *ast.SwitchStmt, env *Environment) (object.Object, error) {
	var subject object.Object
	if stmt.Subject != nil {
		v, err := Eval(stmt.Subject, env)
		if err != nil {
			return nil, err
		}
		subject = v
	}
	for _, branch := range stmt.Branches {
		matched := branch.Default
		if !matched {
			val, err := Eval(branch.Values[0], env)
			if err != nil {
				return nil, err
			}
			if subject == nil || branch.Guard {
				b, ok := val.(*object.Bool)
				if !ok {
					return nil, &EvalError{Pos: branch.Position, Reason: fmt.Sprintf("当 的条件必须是布尔，实际是 %s", val.Kind())}
				}
				matched = b.Value
			} else {
				if ok, err := safeEqual(subject, val); err == nil && ok {
					matched = true
				}
			}
		}
		if !matched && subject != nil && !branch.Guard {
			for _, valueExpr := range branch.Values {
				value, err := Eval(valueExpr, env)
				if err != nil {
					return nil, err
				}
				ok, err := safeEqual(subject, value)
				if err != nil {
					return nil, err
				}
				if ok {
					matched = true
					break
				}
			}
		}
		if !matched {
			continue
		}
		if _, err := evalBlock(branch.Body, NewEnclosedEnvironment(env)); err != nil {
			if sig, ok := err.(*controlSignal); ok && sig.kind == controlBreak {
				return object.Null, nil
			}
			return nil, err
		}
		return object.Null, nil
	}
	return object.Null, nil
}

func safeEqual(left, right object.Object) (bool, error) {
	result, err := evalEq(ast.Position{}, left, right)
	if err != nil {
		return false, nil
	}
	return result.(*object.Bool).Value, nil
}
