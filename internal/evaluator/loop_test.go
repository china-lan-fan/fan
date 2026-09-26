package evaluator

import (
	"strings"
	"testing"
)

func TestWhileBasic(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`当 计数 小于 10 循环`,
		`    计数 = 计数 加 1`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 10)
}

func TestWhileZeroIterations(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`当 假 循环`,
		`    计数 = 计数 加 1`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 0)
}

func TestWhileConditionMustBeBool(t *testing.T) {
	mustError(t, "当 1 循环\n结束", "必须是布尔")
	mustError(t, `当 "是" 循环`+"\n结束", "必须是布尔")
}

func TestWhileAccumulate(t *testing.T) {
	src := strings.Join([]string{
		`变量 甲 = 1`,
		`变量 合计 = 0`,
		`当 甲 小于等于 10 循环`,
		`    合计 = 合计 加 甲`,
		`    甲 = 甲 加 1`,
		`结束`,
		`合计`,
	}, "\n")
	mustInt(t, src, 55)
}

func TestWhileBreak(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`当 真 循环`,
		`    计数 = 计数 加 1`,
		`    如果 计数 等于 5 那么`,
		`        跳出`,
		`    结束`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 5)
}

func TestWhileContinue(t *testing.T) {
	src := strings.Join([]string{
		`变量 甲 = 0`,
		`变量 合计 = 0`,
		`当 甲 小于 5 循环`,
		`    甲 = 甲 加 1`,
		`    如果 甲 等于 3 那么`,
		`        继续`,
		`    结束`,
		`    合计 = 合计 加 甲`,
		`结束`,
		`合计`,
	}, "\n")
	mustInt(t, src, 12)
}

func TestRepeatRunsAtLeastOnce(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`重复`,
		`    计数 = 计数 加 1`,
		`直到 真 结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 1)
}

func TestRepeatUntilCondition(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`重复`,
		`    计数 = 计数 加 1`,
		`直到 计数 等于 10 结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 10)
}

func TestRepeatNewlineEnd(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`重复`,
		`    计数 = 计数 加 1`,
		`直到 计数 等于 3`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 3)
}

func TestRepeatConditionMustBeBool(t *testing.T) {
	mustError(t, "重复\n直到 1 结束", "必须是布尔")
}

func TestRepeatBreak(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`重复`,
		`    计数 = 计数 加 1`,
		`    如果 计数 等于 4 那么`,
		`        跳出`,
		`    结束`,
		`直到 假 结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 4)
}

func TestRepeatContinueStillChecksUntil(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`变量 合计 = 0`,
		`重复`,
		`    计数 = 计数 加 1`,
		`    如果 计数 等于 2 那么`,
		`        继续`,
		`    结束`,
		`    合计 = 合计 加 计数`,
		`直到 计数 等于 4 结束`,
		`合计`,
	}, "\n")
	mustInt(t, src, 8)
}

func TestLoopBlockScope(t *testing.T) {
	src := strings.Join([]string{
		`变量 甲 = 0`,
		`当 甲 小于 3 循环`,
		`    变量 内部 = 1`,
		`    甲 = 甲 加 1`,
		`结束`,
		`内部`,
	}, "\n")
	mustError(t, src, "未声明")
}

func TestLoopScopePerIteration(t *testing.T) {
	src := strings.Join([]string{
		`变量 甲 = 0`,
		`当 甲 小于 3 循环`,
		`    变量 内部 = 甲`,
		`    甲 = 甲 加 1`,
		`结束`,
		`甲`,
	}, "\n")
	mustInt(t, src, 3)
}

func TestNestedLoopBreakInnerOnly(t *testing.T) {
	src := strings.Join([]string{
		`变量 外 = 0`,
		`变量 次数 = 0`,
		`当 外 小于 3 循环`,
		`    外 = 外 加 1`,
		`    变量 内 = 0`,
		`    当 真 循环`,
		`        内 = 内 加 1`,
		`        次数 = 次数 加 1`,
		`        如果 内 等于 2 那么`,
		`            跳出`,
		`        结束`,
		`    结束`,
		`结束`,
		`次数`,
	}, "\n")
	mustInt(t, src, 6)
}

func TestBreakOutsideLoopError(t *testing.T) {
	mustError(t, "跳出", "只能用于循环内")
}

func TestContinueOutsideLoopError(t *testing.T) {
	mustError(t, "继续", "只能用于循环内")
}

func TestBreakInIfOutsideLoopError(t *testing.T) {
	mustError(t, "如果 真 那么\n    跳出\n结束", "只能用于循环内")
}

func TestLoopWithArray(t *testing.T) {
	src := strings.Join([]string{
		`变量 名单 = [10, 20, 30]`,
		`变量 下标 = 0`,
		`变量 合计 = 0`,
		`当 下标 小于 长度(名单) 循环`,
		`    合计 = 合计 加 名单[下标]`,
		`    下标 = 下标 加 1`,
		`结束`,
		`合计`,
	}, "\n")
	mustInt(t, src, 60)
}

func TestLoopBuildArray(t *testing.T) {
	src := strings.Join([]string{
		`变量 结果 = []`,
		`变量 甲 = 1`,
		`当 甲 小于等于 3 循环`,
		`    结果 = 追加(结果, 甲)`,
		`    甲 = 甲 加 1`,
		`结束`,
		`结果`,
	}, "\n")
	mustInspect(t, src, "[1, 2, 3]")
}

func TestTimesBasic(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`重复 3 次`,
		`    计数 = 计数 加 1`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 3)
}

func TestTimesZeroAndNegative(t *testing.T) {
	for _, count := range []string{"0", "-1", "0 减 5"} {
		src := strings.Join([]string{
			`变量 计数 = 0`,
			`重复 ` + count + ` 次`,
			`    计数 = 计数 加 1`,
			`结束`,
			`计数`,
		}, "\n")
		mustInt(t, src, 0)
	}
}

func TestTimesCountFromVariable(t *testing.T) {
	src := strings.Join([]string{
		`变量 次数 = 4`,
		`变量 计数 = 0`,
		`重复 次数 次`,
		`    计数 = 计数 加 1`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 4)
}

func TestTimesCountFromExpression(t *testing.T) {
	src := strings.Join([]string{
		`变量 名单 = [1, 2, 3]`,
		`变量 计数 = 0`,
		`重复 长度(名单) 次`,
		`    计数 = 计数 加 1`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 3)
}

func TestTimesCountMustBeInt(t *testing.T) {
	mustError(t, "重复 2.5 次\n结束", "次数必须是整数")
	mustError(t, `重复 "三" 次`+"\n结束", "次数必须是整数")
	mustError(t, "重复 真 次\n结束", "次数必须是整数")
}

func TestTimesCountSnapshot(t *testing.T) {
	src := strings.Join([]string{
		`变量 次数 = 3`,
		`变量 计数 = 0`,
		`重复 次数 次`,
		`    计数 = 计数 加 1`,
		`    次数 = 100`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 3)
}

func TestTimesBreak(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`重复 10 次`,
		`    计数 = 计数 加 1`,
		`    如果 计数 等于 4 那么`,
		`        跳出`,
		`    结束`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 4)
}

func TestTimesContinueConsumesIteration(t *testing.T) {
	src := strings.Join([]string{
		`变量 轮 = 0`,
		`变量 合计 = 0`,
		`重复 5 次`,
		`    轮 = 轮 加 1`,
		`    如果 轮 等于 2 那么`,
		`        继续`,
		`    结束`,
		`    合计 = 合计 加 轮`,
		`结束`,
		`轮`,
	}, "\n")
	mustInt(t, src, 5)
}

func TestTimesContinueSkipsBody(t *testing.T) {
	src := strings.Join([]string{
		`变量 轮 = 0`,
		`变量 合计 = 0`,
		`重复 5 次`,
		`    轮 = 轮 加 1`,
		`    如果 轮 等于 2 那么`,
		`        继续`,
		`    结束`,
		`    合计 = 合计 加 轮`,
		`结束`,
		`合计`,
	}, "\n")
	mustInt(t, src, 13)
}

func TestTimesBlockScope(t *testing.T) {
	src := strings.Join([]string{
		`重复 2 次`,
		`    变量 内部 = 1`,
		`结束`,
		`内部`,
	}, "\n")
	mustError(t, src, "未声明")
}

func TestTimesNested(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`重复 3 次`,
		`    重复 4 次`,
		`        计数 = 计数 加 1`,
		`    结束`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 12)
}

func TestTimesBuildArray(t *testing.T) {
	src := strings.Join([]string{
		`变量 结果 = []`,
		`变量 甲 = 0`,
		`重复 3 次`,
		`    甲 = 甲 加 1`,
		`    结果 = 追加(结果, 甲)`,
		`结束`,
		`结果`,
	}, "\n")
	mustInspect(t, src, "[1, 2, 3]")
}

func TestRepeatUntilStillWorksEval(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`重复`,
		`    计数 = 计数 加 1`,
		`直到 计数 等于 2 结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 2)
}

func TestForEachArray(t *testing.T) {
	src := strings.Join([]string{
		`变量 名单 = ["小明", "小红", "小刚"]`,
		`变量 结果 = []`,
		`遍历 名单 中的 人`,
		`    结果 = 追加(结果, 人)`,
		`结束`,
		`结果`,
	}, "\n")
	mustInspect(t, src, "[小明, 小红, 小刚]")
}

func TestForEachWithIndex(t *testing.T) {
	src := strings.Join([]string{
		`变量 名单 = ["甲", "乙", "丙"]`,
		`变量 合计 = 0`,
		`遍历 名单 中的 人 序号`,
		`    合计 = 合计 加 序号`,
		`结束`,
		`合计`,
	}, "\n")
	mustInt(t, src, 3)
}

func TestForEachIndexStartsAtZero(t *testing.T) {
	src := strings.Join([]string{
		`变量 名单 = ["甲", "乙"]`,
		`变量 第一个序号 = -1`,
		`遍历 名单 中的 人 序号`,
		`    如果 序号 等于 0 那么`,
		`        第一个序号 = 0`,
		`    结束`,
		`结束`,
		`第一个序号`,
	}, "\n")
	mustInt(t, src, 0)
}

func TestForEachString(t *testing.T) {
	src := strings.Join([]string{
		`变量 结果 = []`,
		`遍历 "你好" 中的 字`,
		`    结果 = 追加(结果, 字)`,
		`结束`,
		`结果`,
	}, "\n")
	mustInspect(t, src, "[你, 好]")
}

func TestForEachStringIndex(t *testing.T) {
	src := strings.Join([]string{
		`变量 长度 = 0`,
		`遍历 "你好" 中的 字 序号`,
		`    长度 = 序号 加 1`,
		`结束`,
		`长度`,
	}, "\n")
	mustInt(t, src, 2)
}

func TestForEachEmptyCollection(t *testing.T) {
	src := strings.Join([]string{
		`变量 计数 = 0`,
		`遍历 [] 中的 项`,
		`    计数 = 计数 加 1`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 0)
}

func TestForEachNotIterable(t *testing.T) {
	mustError(t, "遍历 1 中的 项\n结束", "不可遍历")
	mustError(t, "遍历 真 中的 项\n结束", "不可遍历")
}

func TestForEachBreak(t *testing.T) {
	src := strings.Join([]string{
		`变量 结果 = []`,
		`遍历 [1, 2, 3, 4, 5] 中的 项`,
		`    如果 项 大于 3 那么`,
		`        跳出`,
		`    结束`,
		`    结果 = 追加(结果, 项)`,
		`结束`,
		`结果`,
	}, "\n")
	mustInspect(t, src, "[1, 2, 3]")
}

func TestForEachContinue(t *testing.T) {
	src := strings.Join([]string{
		`变量 结果 = []`,
		`遍历 [1, 2, 3, 4] 中的 项`,
		`    如果 项 等于 2 那么`,
		`        继续`,
		`    结束`,
		`    结果 = 追加(结果, 项)`,
		`结束`,
		`结果`,
	}, "\n")
	mustInspect(t, src, "[1, 3, 4]")
}

func TestForEachSnapshotLength(t *testing.T) {
	src := strings.Join([]string{
		`变量 名单 = [1, 2, 3]`,
		`变量 计数 = 0`,
		`遍历 名单 中的 项`,
		`    名单 = 追加(名单, 99)`,
		`    计数 = 计数 加 1`,
		`结束`,
		`计数`,
	}, "\n")
	mustInt(t, src, 3)
}

func TestForEachSeesUpdatedElements(t *testing.T) {
	src := strings.Join([]string{
		`变量 名单 = [10, 20, 30]`,
		`遍历 名单 中的 项 序号`,
		`    如果 序号 等于 1 那么`,
		`        名单[1] = 99`,
		`    结束`,
		`结束`,
		`名单`,
	}, "\n")
	mustInspect(t, src, "[10, 99, 30]")
}

func TestForEachVariableScope(t *testing.T) {
	src := strings.Join([]string{
		`遍历 [1] 中的 项`,
		`    变量 内部 = 1`,
		`结束`,
		`内部`,
	}, "\n")
	mustError(t, src, "未声明")
}

func TestForEachShadowsOuter(t *testing.T) {
	src := strings.Join([]string{
		`变量 项 = "外层"`,
		`遍历 [1] 中的 项`,
		`    变量 内层 = 项`,
		`结束`,
		`项`,
	}, "\n")
	mustInspect(t, src, "外层")
}

func TestForEachNested(t *testing.T) {
	src := strings.Join([]string{
		`变量 结果 = []`,
		`遍历 [[1, 2], [3, 4]] 中的 行`,
		`    遍历 行 中的 单元`,
		`        结果 = 追加(结果, 单元)`,
		`    结束`,
		`结束`,
		`结果`,
	}, "\n")
	mustInspect(t, src, "[1, 2, 3, 4]")
}
