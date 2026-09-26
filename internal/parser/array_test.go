package parser

import (
	"testing"

	"fan/internal/ast"
)

func TestArrayLiteral(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{"[]", "[]"},
		{"[1, 2, 3]", "[1, 2, 3]"},
		{"[1, 2, 3,]", "[1, 2, 3]"},
		{`[1, "二", 真]`, `[1, "二", 真]`},
		{"[[1, 2], [3, 4]]", "[[1, 2], [3, 4]]"},
	}
	for _, c := range cases {
		prog := mustParse(t, c.src)
		stmt := prog.Statements[0].(*ast.ExpressionStmt)
		if got := stmt.Expression.String(); got != c.want {
			t.Errorf("源码 %q：期望 %s，实际 %s", c.src, c.want, got)
		}
	}
}

func TestArrayLiteralMultiline(t *testing.T) {
	src := "变量 名单 = [\n    \"小明\",\n    \"小红\",\n]"
	prog := mustParse(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("跨行数组应为 1 条语句，实际 %d", len(prog.Statements))
	}
	decl := prog.Statements[0].(*ast.VarDecl)
	arr := decl.Value.(*ast.ArrayLiteral)
	if len(arr.Elements) != 2 {
		t.Fatalf("应有 2 个元素，实际 %d", len(arr.Elements))
	}
}

func TestArrayTypeDecl(t *testing.T) {
	prog := mustParse(t, `变量 数组 名单 = ["小明"]`)
	decl := prog.Statements[0].(*ast.VarDecl)
	if decl.DeclType != ast.TypeArray {
		t.Fatalf("类型应为数组，实际 %q", decl.DeclType)
	}
}

func TestIndexExpr(t *testing.T) {
	prog := mustParse(t, "数字[0]")
	stmt := prog.Statements[0].(*ast.ExpressionStmt)
	if got := stmt.Expression.String(); got != "数字[0]" {
		t.Fatalf("期望 数字[0]，实际 %s", got)
	}
}

func TestNestedIndexExpr(t *testing.T) {
	prog := mustParse(t, "嵌套[0][1]")
	stmt := prog.Statements[0].(*ast.ExpressionStmt)
	if got := stmt.Expression.String(); got != "嵌套[0][1]" {
		t.Fatalf("期望 嵌套[0][1]，实际 %s", got)
	}
}

func TestIndexAssign(t *testing.T) {
	for _, src := range []string{"数字[0] = 99", "数字[0] 为 99"} {
		prog := mustParse(t, src)
		if _, ok := prog.Statements[0].(*ast.IndexAssignStmt); !ok {
			t.Fatalf("源码 %q 应为下标赋值，实际 %T", src, prog.Statements[0])
		}
	}
}

func TestCallExpr(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{"长度([1,2])", "长度([1, 2])"},
		{"打印()", "打印()"},
		{"追加(甲, 1, 2)", "追加(甲, 1, 2)"},
		{"打印(1, 2,)", "打印(1, 2)"},
	}
	for _, c := range cases {
		prog := mustParse(t, c.src)
		stmt := prog.Statements[0].(*ast.ExpressionStmt)
		if got := stmt.Expression.String(); got != c.want {
			t.Errorf("源码 %q：期望 %s，实际 %s", c.src, c.want, got)
		}
	}
}

func TestArrayConcatExpr(t *testing.T) {
	prog := mustParse(t, "[1, 2] 加 [3]")
	stmt := prog.Statements[0].(*ast.ExpressionStmt)
	if got := stmt.Expression.String(); got != "([1, 2] 加 [3])" {
		t.Fatalf("期望 ([1, 2] 加 [3])，实际 %s", got)
	}
}

func TestArrayErrors(t *testing.T) {
	expectError(t, "[1, 2")
	expectError(t, "数字[0")
	expectError(t, "长度(1")
}

func TestDictLiteral(t *testing.T) {
	src := `{"a": 1, "b": 2}`
	prog := mustParse(t, src)
	dict, ok := prog.Statements[0].(*ast.ExpressionStmt).Expression.(*ast.DictLiteral)
	if !ok {
		t.Fatalf("应为字典字面量，实际 %T", prog.Statements[0])
	}
	if len(dict.Pairs) != 2 {
		t.Fatalf("应有 2 对，实际 %d", len(dict.Pairs))
	}
}

func TestDictLiteralMultiline(t *testing.T) {
	src := "变量 d = {\n    \"x\": 1,\n    \"y\": 2,\n}\nd[\"x\"] = 3"
	prog := mustParse(t, src)
	if len(prog.Statements) != 2 {
		t.Fatalf("应解析为 2 条语句，实际 %d", len(prog.Statements))
	}
	if _, ok := prog.Statements[1].(*ast.IndexAssignStmt); !ok {
		t.Fatalf("第二条应为下标赋值，实际 %T", prog.Statements[1])
	}
}

func TestDictType(t *testing.T) {
	src := "变量 字典 d = {}"
	prog := mustParse(t, src)
	decl := prog.Statements[0].(*ast.VarDecl)
	if decl.DeclType != ast.TypeDict {
		t.Fatalf("声明类型应为字典，实际 %s", decl.DeclType)
	}
}
