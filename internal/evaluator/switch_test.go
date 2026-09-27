package evaluator

import "testing"

func TestSwitchValue(t *testing.T) {
	src := `变量 x = 2
变量 结果 = ""
判断 x
    为 1
        结果 = "一"
    为 2、3
        结果 = "多"
    其他
        结果 = "无"
结束
结果 等于 "多"`
	mustBool(t, src, true)
}

func TestSwitchGuard(t *testing.T) {
	src := `变量 分数 = 85
变量 结果 = ""
判断
    当 分数 >= 90
        结果 = "优秀"
    当 分数 >= 60
        结果 = "及格"
    其他
        结果 = "不及格"
结束
结果 等于 "及格"`
	mustBool(t, src, true)
}

func TestSwitchMixed(t *testing.T) {
	src := `变量 x = 500
变量 结果 = ""
判断 x
    为 200
        结果 = "成功"
    当 x >= 500
        结果 = "错误"
结束
结果 等于 "错误"`
	mustBool(t, src, true)
}

func TestSwitchGuardRequiresBool(t *testing.T) {
	mustError(t, "判断\n    当 1\n结束", "布尔")
}

func TestSwitchNoFallthrough(t *testing.T) {
	src := `变量 x = 1
变量 n = 0
判断 x
    为 1
        n++
    为 1
        n++
结束
n 等于 1`
	mustBool(t, src, true)
}
