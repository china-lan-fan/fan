package parser

import (
	"strings"
	"testing"

	"fan/internal/ast"
)

func TestWhileStatement(t *testing.T) {
	src := strings.Join([]string{
		`当 计数 小于 10 循环`,
		`    计数 = 计数 加 1`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	stmt, ok := prog.Statements[0].(*ast.WhileStmt)
	if !ok {
		t.Fatalf("应为 当 循环语句，实际 %T", prog.Statements[0])
	}
	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("块内应有 1 条语句，实际 %d", len(stmt.Body.Statements))
	}
}

func TestWhileEmptyBody(t *testing.T) {
	prog := mustParse(t, "当 真 循环\n结束")
	stmt := prog.Statements[0].(*ast.WhileStmt)
	if len(stmt.Body.Statements) != 0 {
		t.Fatal("空块应无语句")
	}
}

func TestWhileMissingLoopKeyword(t *testing.T) {
	expectError(t, "当 真\n    甲 = 1\n结束")
}

func TestWhileMissingEnd(t *testing.T) {
	expectError(t, "当 真 循环\n    甲 = 1")
}

func TestRepeatSameLineEnd(t *testing.T) {
	src := strings.Join([]string{
		`重复`,
		`    计数 = 计数 加 1`,
		`直到 计数 等于 10 结束`,
	}, "\n")
	prog := mustParse(t, src)
	stmt, ok := prog.Statements[0].(*ast.RepeatStmt)
	if !ok {
		t.Fatalf("应为 重复 循环语句，实际 %T", prog.Statements[0])
	}
	if stmt.Until == nil {
		t.Fatal("应有 直到 条件")
	}
}

func TestRepeatNewlineEnd(t *testing.T) {
	src := strings.Join([]string{
		`重复`,
		`    计数 = 计数 加 1`,
		`直到 计数 等于 10`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	if _, ok := prog.Statements[0].(*ast.RepeatStmt); !ok {
		t.Fatalf("换行写 结束 应合法，实际 %T", prog.Statements[0])
	}
}

func TestRepeatMissingUntil(t *testing.T) {
	expectError(t, "重复\n    甲 = 1\n结束")
}

func TestRepeatMissingEnd(t *testing.T) {
	expectError(t, "重复\n    甲 = 1\n直到 真")
}

func TestBreakContinueParse(t *testing.T) {
	src := strings.Join([]string{
		`当 真 循环`,
		`    如果 甲 那么`,
		`        继续`,
		`    结束`,
		`    跳出`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	while := prog.Statements[0].(*ast.WhileStmt)
	if len(while.Body.Statements) != 2 {
		t.Fatalf("循环体应有 2 条语句，实际 %d", len(while.Body.Statements))
	}
	if _, ok := while.Body.Statements[1].(*ast.BreakStmt); !ok {
		t.Fatalf("第 2 条应为 跳出，实际 %T", while.Body.Statements[1])
	}
	ifStmt := while.Body.Statements[0].(*ast.IfStmt)
	if _, ok := ifStmt.Branches[0].Body.Statements[0].(*ast.ContinueStmt); !ok {
		t.Fatalf("应为 继续，实际 %T", ifStmt.Branches[0].Body.Statements[0])
	}
}

func TestNestedLoops(t *testing.T) {
	src := strings.Join([]string{
		`当 甲 循环`,
		`    重复`,
		`        乙 = 1`,
		`    直到 真 结束`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	while := prog.Statements[0].(*ast.WhileStmt)
	if _, ok := while.Body.Statements[0].(*ast.RepeatStmt); !ok {
		t.Fatalf("应为嵌套 重复 循环，实际 %T", while.Body.Statements[0])
	}
}

func TestTimesStatement(t *testing.T) {
	src := strings.Join([]string{
		`重复 3 次`,
		`    打印("你好")`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	stmt, ok := prog.Statements[0].(*ast.TimesStmt)
	if !ok {
		t.Fatalf("应为计数循环，实际 %T", prog.Statements[0])
	}
	if stmt.Count.String() != "3" {
		t.Fatalf("次数应为 3，实际 %s", stmt.Count.String())
	}
	if len(stmt.Body.Statements) != 1 {
		t.Fatalf("块内应有 1 条语句，实际 %d", len(stmt.Body.Statements))
	}
}

func TestTimesCountExpressions(t *testing.T) {
	cases := []struct {
		src  string
		want string
	}{
		{"重复 3 次\n结束", "3"},
		{"重复 甲 次\n结束", "甲"},
		{"重复 长度(名单) 次\n结束", "长度(名单)"},
		{"重复 甲 加 2 次\n结束", "(甲 加 2)"},
	}
	for _, c := range cases {
		prog := mustParse(t, c.src)
		stmt := prog.Statements[0].(*ast.TimesStmt)
		if got := stmt.Count.String(); got != c.want {
			t.Errorf("源码 %q：次数期望 %s，实际 %s", c.src, c.want, got)
		}
	}
}

func TestTimesEmptyBody(t *testing.T) {
	prog := mustParse(t, "重复 3 次\n结束")
	stmt := prog.Statements[0].(*ast.TimesStmt)
	if len(stmt.Body.Statements) != 0 {
		t.Fatal("空块应无语句")
	}
}

func TestTimesMissingTimesKeyword(t *testing.T) {
	expectError(t, "重复 3\n    甲 = 1\n结束")
}

func TestTimesMissingEnd(t *testing.T) {
	expectError(t, "重复 3 次\n    甲 = 1")
}

func TestTimesCannotUseUntil(t *testing.T) {
	expectError(t, "重复 3 次\n    甲 = 1\n直到 真 结束")
}

func TestRepeatUntilStillWorks(t *testing.T) {
	prog := mustParse(t, "重复\n    甲 = 1\n直到 真 结束")
	if _, ok := prog.Statements[0].(*ast.RepeatStmt); !ok {
		t.Fatalf("重复 后换行应为 直到 风格，实际 %T", prog.Statements[0])
	}
}

func TestTimesWithBreakContinue(t *testing.T) {
	src := strings.Join([]string{
		`重复 5 次`,
		`    如果 甲 那么`,
		`        继续`,
		`    结束`,
		`    跳出`,
		`结束`,
	}, "\n")
	prog := mustParse(t, src)
	stmt := prog.Statements[0].(*ast.TimesStmt)
	if len(stmt.Body.Statements) != 2 {
		t.Fatalf("块内应有 2 条语句，实际 %d", len(stmt.Body.Statements))
	}
}

func TestForEachBasic(t *testing.T) {
	src := "遍历 名单 中的 人\n    打印(人)\n结束"
	prog := mustParse(t, src)
	stmt, ok := prog.Statements[0].(*ast.ForEachStmt)
	if !ok {
		t.Fatalf("应为遍历语句，实际 %T", prog.Statements[0])
	}
	if stmt.FirstName != "人" {
		t.Fatalf("元素名应为 人，实际 %s", stmt.FirstName)
	}
	if stmt.SecondName != "" {
		t.Fatalf("无序号时应为空，实际 %s", stmt.SecondName)
	}
}

func TestForEachWithIndex(t *testing.T) {
	src := "遍历 名单 中的 人 序号\n    打印(序号, 人)\n结束"
	prog := mustParse(t, src)
	stmt := prog.Statements[0].(*ast.ForEachStmt)
	if stmt.FirstName != "人" || stmt.SecondName != "序号" {
		t.Fatalf("元素/序号不符：%s %s", stmt.FirstName, stmt.SecondName)
	}
}

func TestForEachMissingIn(t *testing.T) {
	expectError(t, "遍历 名单 人\n结束")
}

func TestForEachMissingValueName(t *testing.T) {
	expectError(t, "遍历 名单 中的\n结束")
}

func TestForEachMissingEnd(t *testing.T) {
	expectError(t, "遍历 名单 中的 人\n    打印(人)")
}

func TestForEachIterableExpression(t *testing.T) {
	src := "遍历 追加(名单, 1) 中的 项\n    打印(项)\n结束"
	prog := mustParse(t, src)
	stmt := prog.Statements[0].(*ast.ForEachStmt)
	if got := stmt.Iterable.String(); got != "追加(名单, 1)" {
		t.Fatalf("集合表达式解析不符：%s", got)
	}
}

func TestForEachStringIterable(t *testing.T) {
	src := `遍历 "你好" 中的 字` + "\n    打印(字)\n结束"
	prog := mustParse(t, src)
	if _, ok := prog.Statements[0].(*ast.ForEachStmt); !ok {
		t.Fatalf("字符串遍历应合法，实际 %T", prog.Statements[0])
	}
}
