package evaluator

import (
	"strings"
	"testing"
)

func TestIfTrueBranch(t *testing.T) {
	src := strings.Join([]string{
		`变量 等级 = ""`,
		`如果 80 大于等于 60 那么`,
		`    等级 = "及格"`,
		`结束`,
		`等级`,
	}, "\n")
	res, err := runSource(t, src)
	if err != nil {
		t.Fatalf("运行错误：%v", err)
	}
	if res.Inspect() != "及格" {
		t.Fatalf("期望 及格，实际 %s", res.Inspect())
	}
}

func TestIfFalseNoElse(t *testing.T) {
	src := strings.Join([]string{
		`变量 等级 = "初始"`,
		`如果 50 大于等于 60 那么`,
		`    等级 = "及格"`,
		`结束`,
		`等级`,
	}, "\n")
	res, err := runSource(t, src)
	if err != nil {
		t.Fatalf("运行错误：%v", err)
	}
	if res.Inspect() != "初始" {
		t.Fatalf("期望 初始，实际 %s", res.Inspect())
	}
}

func TestIfElseIfChain(t *testing.T) {
	build := func(score string) string {
		return strings.Join([]string{
			`变量 分数 = ` + score,
			`变量 等级 = ""`,
			`如果 分数 大于等于 90 那么`,
			`    等级 = "优"`,
			`否则如果 分数 大于等于 60 那么`,
			`    等级 = "及格"`,
			`否则`,
			`    等级 = "不及格"`,
			`结束`,
			`等级`,
		}, "\n")
	}
	cases := map[string]string{
		"95": "优",
		"70": "及格",
		"30": "不及格",
	}
	for score, want := range cases {
		res, err := runSource(t, build(score))
		if err != nil {
			t.Fatalf("分数 %s 运行错误：%v", score, err)
		}
		if res.Inspect() != want {
			t.Errorf("分数 %s 期望 %s，实际 %s", score, want, res.Inspect())
		}
	}
}

func TestIfElseSpaceForm(t *testing.T) {
	src := strings.Join([]string{
		`变量 结果 = 0`,
		`如果 假 那么`,
		`    结果 = 1`,
		`否则 如果 真 那么`,
		`    结果 = 2`,
		`结束`,
		`结果`,
	}, "\n")
	mustIntSrc(t, src, 2)
}

func TestIfNested(t *testing.T) {
	src := strings.Join([]string{
		`变量 甲 = 1`,
		`变量 乙 = 2`,
		`变量 象限 = 0`,
		`如果 甲 大于 0 那么`,
		`    如果 乙 大于 0 那么`,
		`        象限 = 1`,
		`    结束`,
		`结束`,
		`象限`,
	}, "\n")
	mustIntSrc(t, src, 1)
}

func TestIfConditionMustBeBool(t *testing.T) {
	mustError(t, "如果 1 那么\n结束", "必须是布尔")
	mustError(t, "如果 \"是\" 那么\n结束", "必须是布尔")
}

func TestIfBlockScope(t *testing.T) {
	src := strings.Join([]string{
		`如果 真 那么`,
		`    变量 内部 = 1`,
		`结束`,
		`内部`,
	}, "\n")
	mustError(t, src, "未声明")
}

func TestIfBlockCanModifyOuter(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 1`,
		`如果 真 那么`,
		`    计数 = 计数 加 10`,
		`结束`,
		`计数`,
	}, "\n")
	mustIntSrc(t, src, 11)
}

func TestIfBlockShadowsOuterDeclaration(t *testing.T) {
	src := strings.Join([]string{
		`变量 甲 = 1`,
		`如果 真 那么`,
		`    变量 甲 = 99`,
		`结束`,
		`甲`,
	}, "\n")
	mustIntSrc(t, src, 1)
}

func TestIfConstInBlock(t *testing.T) {
	src := strings.Join([]string{
		`如果 真 那么`,
		`    常量 甲 = 1`,
		`    甲 = 2`,
		`结束`,
	}, "\n")
	mustError(t, src, "常量")
}

func TestIfWithLogicCondition(t *testing.T) {
	src := strings.Join([]string{
		`变量 甲 = 5`,
		`变量 结果 = 0`,
		`如果 甲 大于 0 与 甲 小于 10 那么`,
		`    结果 = 1`,
		`结束`,
		`结果`,
	}, "\n")
	mustIntSrc(t, src, 1)
}

func TestIfEmptyBlockReturnsNil(t *testing.T) {
	res, err := runSource(t, "如果 真 那么\n结束")
	if err != nil {
		t.Fatalf("运行错误：%v", err)
	}
	if res.Inspect() != "空" {
		t.Fatalf("空块应返回 空，实际 %s", res.Inspect())
	}
}

func mustIntSrc(t *testing.T, src string, want int64) {
	t.Helper()
	mustInt(t, src, want)
}
