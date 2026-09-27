package parser

import (
	"testing"

	"fan/internal/ast"
)

func TestSwitchValueBranches(t *testing.T) {
	src := `判断 x
    为 1、2
        打印(x)
    其他
        打印(0)
结束`
	prog := mustParse(t, src)
	stmt, ok := prog.Statements[0].(*ast.SwitchStmt)
	if !ok {
		t.Fatalf("应为判断语句，实际 %T", prog.Statements[0])
	}
	if len(stmt.Branches) != 2 {
		t.Fatalf("应有两个分支，实际 %d", len(stmt.Branches))
	}
	if len(stmt.Branches[0].Values) != 2 {
		t.Fatalf("首分支应有两个匹配值")
	}
	if !stmt.Branches[1].Default {
		t.Fatalf("第二分支应是其他分支")
	}
}

func TestSwitchGuardBranches(t *testing.T) {
	src := `判断
    当 x >= 1
        打印(x)
结束`
	prog := mustParse(t, src)
	stmt := prog.Statements[0].(*ast.SwitchStmt)
	if stmt.Subject != nil {
		t.Fatalf("条件判断不应有匹配对象")
	}
	if !stmt.Branches[0].Guard {
		t.Fatalf("首分支应是条件分支")
	}
}

func TestSwitchInvalidBranches(t *testing.T) {
	expectError(t, "判断\n    为 1\n结束")
	expectError(t, "判断 x\n    其他\n        打印(0)\n    为 1\n结束")
	expectError(t, "判断 x\n    为 1")
}
