package parser

import (
	"strings"
	"testing"

	"fan/internal/ast"
)

func TestIfSimple(t *testing.T) {
	src := "如果 分数 大于等于 60 那么\n    等级 = \"及格\"\n结束"
	prog := mustParse(t, src)
	if len(prog.Statements) != 1 {
		t.Fatalf("应为 1 条语句，实际 %d", len(prog.Statements))
	}
	stmt, ok := prog.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("应为条件语句，实际 %T", prog.Statements[0])
	}
	if len(stmt.Branches) != 1 {
		t.Fatalf("应有 1 个分支，实际 %d", len(stmt.Branches))
	}
	if stmt.Else != nil {
		t.Fatal("不应有 否则 分支")
	}
	if len(stmt.Branches[0].Body.Statements) != 1 {
		t.Fatalf("块内应有 1 条语句，实际 %d", len(stmt.Branches[0].Body.Statements))
	}
}

func TestIfElseIfElse(t *testing.T) {
	src := strings.Join([]string{
		`如果 分数 大于等于 90 那么`,
		`    等级 = "优"`,
		`否则如果 分数 大于等于 60 那么`,
		`    等级 = "及格"`,
		`否则`,
		`    等级 = "不及格"`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	stmt := prog.Statements[0].(*ast.IfStmt)
	if len(stmt.Branches) != 2 {
		t.Fatalf("应有 2 个条件分支，实际 %d", len(stmt.Branches))
	}
	if stmt.Else == nil {
		t.Fatal("应有 否则 分支")
	}
}

func TestIfElseSpaceForm(t *testing.T) {
	src := strings.Join([]string{
		`如果 甲 那么`,
		`    乙 = 1`,
		`否则 如果 丙 那么`,
		`    乙 = 2`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	stmt := prog.Statements[0].(*ast.IfStmt)
	if len(stmt.Branches) != 2 {
		t.Fatalf("否则 如果 分开写应识别为 2 个分支，实际 %d", len(stmt.Branches))
	}
}

func TestIfNested(t *testing.T) {
	src := strings.Join([]string{
		`如果 甲 大于 0 那么`,
		`    如果 乙 大于 0 那么`,
		`        象限 = 1`,
		`    结束`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	outer := prog.Statements[0].(*ast.IfStmt)
	inner, ok := outer.Branches[0].Body.Statements[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("应为嵌套条件语句，实际 %T", outer.Branches[0].Body.Statements[0])
	}
	if len(inner.Branches) != 1 {
		t.Fatalf("内层应有 1 个分支，实际 %d", len(inner.Branches))
	}
}

func TestIfEmptyBlock(t *testing.T) {
	prog := mustParse(t, "如果 甲 那么\n结束")
	stmt := prog.Statements[0].(*ast.IfStmt)
	if len(stmt.Branches[0].Body.Statements) != 0 {
		t.Fatal("空块应无语句")
	}
}

func TestIfMissingThen(t *testing.T) {
	expectError(t, "如果 甲\n    乙 = 1\n结束")
}

func TestIfMissingEnd(t *testing.T) {
	expectError(t, "如果 甲 那么\n    乙 = 1")
}

func TestIfElseWithThenIsError(t *testing.T) {
	expectError(t, "如果 甲 那么\n乙 = 1\n否则 那么\n乙 = 2\n结束")
}

func TestIfMultipleStatementsInBlock(t *testing.T) {
	src := strings.Join([]string{
		`如果 真 那么`,
		`    变量 甲 = 1`,
		`    变量 乙 = 2`,
		`    丙 = 甲 加 乙`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	stmt := prog.Statements[0].(*ast.IfStmt)
	if len(stmt.Branches[0].Body.Statements) != 3 {
		t.Fatalf("块内应有 3 条语句，实际 %d", len(stmt.Branches[0].Body.Statements))
	}
}

func TestIfFollowedByOtherStatement(t *testing.T) {
	src := strings.Join([]string{
		`变量 甲 = 1`,
		`如果 真 那么`,
		`    甲 = 2`,
		`结束`,
		`甲 加 1`,
	}, "\n")
	prog := mustParse(t, src)
	if len(prog.Statements) != 3 {
		t.Fatalf("应有 3 条语句，实际 %d", len(prog.Statements))
	}
}
