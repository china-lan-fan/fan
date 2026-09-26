package evaluator

import (
	"fmt"

	"fan/internal/ast"
	"fan/internal/object"
)

type controlKind int

const (
	controlBreak controlKind = iota
	controlContinue
)

type controlSignal struct {
	kind controlKind
	pos  ast.Position
}

func (c *controlSignal) Error() string {
	if c.kind == controlBreak {
		return fmt.Sprintf("第%d行%d列：跳出 只能用于循环内", c.pos.Line, c.pos.Column)
	}
	return fmt.Sprintf("第%d行%d列：继续 只能用于循环内", c.pos.Line, c.pos.Column)
}

func evalWhileStmt(stmt *ast.WhileStmt, env *Environment) (object.Object, error) {
	for {
		cond, err := Eval(stmt.Condition, env)
		if err != nil {
			return nil, err
		}
		if cond.Kind() != object.KindBool {
			return nil, &EvalError{
				Pos:    stmt.Condition.Pos(),
				Reason: fmt.Sprintf("当 的条件必须是布尔，实际是 %s", cond.Kind()),
			}
		}
		if !cond.(*object.Bool).Value {
			return object.Null, nil
		}
		if _, err := evalBlock(stmt.Body, NewEnclosedEnvironment(env)); err != nil {
			sig, ok := err.(*controlSignal)
			if !ok {
				return nil, err
			}
			if sig.kind == controlBreak {
				return object.Null, nil
			}
		}
	}
}

func evalRepeatStmt(stmt *ast.RepeatStmt, env *Environment) (object.Object, error) {
	for {
		if _, err := evalBlock(stmt.Body, NewEnclosedEnvironment(env)); err != nil {
			sig, ok := err.(*controlSignal)
			if !ok {
				return nil, err
			}
			if sig.kind == controlBreak {
				return object.Null, nil
			}
		}
		cond, err := Eval(stmt.Until, env)
		if err != nil {
			return nil, err
		}
		if cond.Kind() != object.KindBool {
			return nil, &EvalError{
				Pos:    stmt.Until.Pos(),
				Reason: fmt.Sprintf("直到 的条件必须是布尔，实际是 %s", cond.Kind()),
			}
		}
		if cond.(*object.Bool).Value {
			return object.Null, nil
		}
	}
}

func evalTimesStmt(stmt *ast.TimesStmt, env *Environment) (object.Object, error) {
	countObj, err := Eval(stmt.Count, env)
	if err != nil {
		return nil, err
	}
	n, ok := countObj.(*object.Integer)
	if !ok {
		return nil, &EvalError{
			Pos:    stmt.Count.Pos(),
			Reason: fmt.Sprintf("次数必须是整数，实际是 %s", countObj.Kind()),
		}
	}
	for i := int64(0); i < n.Value; i++ {
		if _, err := evalBlock(stmt.Body, NewEnclosedEnvironment(env)); err != nil {
			sig, ok := err.(*controlSignal)
			if !ok {
				return nil, err
			}
			if sig.kind == controlBreak {
				return object.Null, nil
			}
		}
	}
	return object.Null, nil
}

func evalForEachStmt(stmt *ast.ForEachStmt, outer *Environment) (object.Object, error) {
	iterObj, err := Eval(stmt.Iterable, outer)
	if err != nil {
		return nil, err
	}
	switch iter := iterObj.(type) {
	case *object.Array:
		count := len(iter.Elements)
		for i := 0; i < count; i++ {
			if i >= len(iter.Elements) {
				break
			}
			env := NewEnclosedEnvironment(outer)
			if err := env.declare(stmt.ValueName, iter.Elements[i], false, ast.TypeAny); err != nil {
				return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
			}
			if stmt.IndexName != "" {
				if err := env.declare(stmt.IndexName, &object.Integer{Value: int64(i)}, false, ast.TypeAny); err != nil {
					return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
				}
			}
			if _, err := evalBlock(stmt.Body, env); err != nil {
				sig, ok := err.(*controlSignal)
				if !ok {
					return nil, err
				}
				if sig.kind == controlBreak {
					return object.Null, nil
				}
			}
		}
		return object.Null, nil
	case *object.String:
		runes := []rune(iter.Value)
		for i, r := range runes {
			env := NewEnclosedEnvironment(outer)
			if err := env.declare(stmt.ValueName, &object.String{Value: string(r)}, false, ast.TypeAny); err != nil {
				return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
			}
			if stmt.IndexName != "" {
				if err := env.declare(stmt.IndexName, &object.Integer{Value: int64(i)}, false, ast.TypeAny); err != nil {
					return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
				}
			}
			if _, err := evalBlock(stmt.Body, env); err != nil {
				sig, ok := err.(*controlSignal)
				if !ok {
					return nil, err
				}
				if sig.kind == controlBreak {
					return object.Null, nil
				}
			}
		}
		return object.Null, nil
	case *object.Dict:
		keys := append([]object.Object(nil), iter.Keys...)
		for _, key := range keys {
			env := NewEnclosedEnvironment(outer)
			if stmt.IndexName == "" {
				if err := env.declare(stmt.ValueName, key, false, ast.TypeAny); err != nil {
					return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
				}
			} else {
				val, _ := iter.Get(key)
				if err := env.declare(stmt.ValueName, val, false, ast.TypeAny); err != nil {
					return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
				}
				if err := env.declare(stmt.IndexName, key, false, ast.TypeAny); err != nil {
					return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
				}
			}
			if _, err := evalBlock(stmt.Body, env); err != nil {
				sig, ok := err.(*controlSignal)
				if !ok {
					return nil, err
				}
				if sig.kind == controlBreak {
					return object.Null, nil
				}
			}
		}
		return object.Null, nil
	default:
		return nil, &EvalError{
			Pos:    stmt.Position,
			Reason: fmt.Sprintf("%s 不可遍历", iterObj.Kind()),
		}
	}
}
