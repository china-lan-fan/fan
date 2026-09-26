package evaluator

import (
	"strings"
	"testing"
)

const rectModel = `
定义 模型 矩形
    变量 整数 宽
    变量 整数 高
结束

定义 模型 的 函数 初始化 矩形 r、整数 宽、整数 高
    r.宽 = 宽
    r.高 = 高
结束

定义 模型 的 函数 面积 矩形 r : 整数
    返回 r.宽 乘 r.高
结束

定义 模型 的 函数 放大 矩形 r、整数 倍数
    r.宽 = r.宽 乘 倍数
    r.高 = r.高 乘 倍数
结束
`

func TestClassInstantiationAndMethod(t *testing.T) {
	src := rectModel + "\n变量 r = 矩形(3, 4)\nr 的 面积()"
	mustInt(t, src, 12)
}

func TestClassMethodMutation(t *testing.T) {
	src := rectModel + "\n变量 r = 矩形(3, 4)\nr.放大(2)\nr 的 面积()"
	mustInt(t, src, 48)
}

func TestClassFieldAccess(t *testing.T) {
	src := rectModel + "\n变量 r = 矩形(3, 4)\nr.宽"
	mustInt(t, src, 3)
}

func TestClassFieldAssignOutside(t *testing.T) {
	src := rectModel + "\n变量 r = 矩形(1, 1)\nr.宽 = 5\nr.高 = 6\nr 的 面积()"
	mustInt(t, src, 30)
}

func TestClassNoInit(t *testing.T) {
	src := strings.Join([]string{
		`模型 点`,
		`    变量 x`,
		`    变量 y`,
		`结束`,
		`变量 p = 点()`,
		`p.x`,
	}, "\n")
	res, err := runSource(t, src)
	if err != nil {
		t.Fatalf("运行错误：%v", err)
	}
	if res.Inspect() != "空" {
		t.Fatalf("未初始化字段应为 空，实际 %s", res.Inspect())
	}
}

func TestClassNoInitWithArgsError(t *testing.T) {
	src := "模型 点\n    变量 x\n结束\n点(1)"
	mustError(t, src, "不应传参")
}

func TestClassUnknownField(t *testing.T) {
	src := "模型 点\n    变量 x\n结束\n变量 p = 点()\np.颜色"
	mustError(t, src, "没有字段或方法")
}

func TestClassCannotAddField(t *testing.T) {
	src := "模型 点\n    变量 x\n结束\n变量 p = 点()\np.颜色 = \"红\""
	mustError(t, src, "没有字段")
}

func TestClassUnknownMethod(t *testing.T) {
	src := "模型 点\n    变量 x\n结束\n变量 p = 点()\np.跑()"
	mustError(t, src, "没有方法")
}

func TestClassAsValue(t *testing.T) {
	src := strings.Join([]string{
		`模型 点`,
		`    变量 x`,
		`结束`,
		`模型 的 函数 初始化 点 p、整数 x`,
		`    p.x = x`,
		`结束`,
		`变量 造 = 点`,
		`变量 p = 造(5)`,
		`p.x`,
	}, "\n")
	mustInt(t, src, 5)
}

func TestClassEmbed(t *testing.T) {
	src := strings.Join([]string{
		`模型 人`,
		`    变量 字符串 姓名`,
		`结束`,
		`模型 的 函数 问好 人 r`,
		`    返回 r.姓名`,
		`结束`,
		`模型 学生`,
		`    嵌入 人`,
		`    变量 整数 学号`,
		`结束`,
		`变量 s = 学生()`,
		`s.姓名 = "小明"`,
		`s.问好()`,
	}, "\n")
	res, err := runSource(t, src)
	if err != nil {
		t.Fatalf("运行错误：%v", err)
	}
	if res.Inspect() != "小明" {
		t.Fatalf("嵌入方法提升应返回 小明，实际 %s", res.Inspect())
	}
}

func TestClassDuckTyping(t *testing.T) {
	src := strings.Join([]string{
		`模型 狗`,
		`    变量 x`,
		`结束`,
		`模型 的 函数 说话 狗 d`,
		`    返回 "汪汪"`,
		`结束`,
		`模型 猫`,
		`    变量 x`,
		`结束`,
		`模型 的 函数 说话 猫 c`,
		`    返回 "喵喵"`,
		`结束`,
		`函数 让它说(动物)`,
		`    返回 动物.说话()`,
		`结束`,
		`让它说(狗())`,
	}, "\n")
	res, err := runSource(t, src)
	if err != nil {
		t.Fatalf("运行错误：%v", err)
	}
	if res.Inspect() != "汪汪" {
		t.Fatalf("鸭子类型调用失败：%s", res.Inspect())
	}
}

func TestClassMethodReturnsInstance(t *testing.T) {
	src := strings.Join([]string{
		`模型 计数器`,
		`    变量 值`,
		`结束`,
		`模型 的 函数 初始化 计数器 c`,
		`    c.值 = 0`,
		`结束`,
		`模型 的 函数 递增 计数器 c`,
		`    c.值 = c.值 加 1`,
		`    返回 c`,
		`结束`,
		`变量 c = 计数器()`,
		`c.递增().递增()`,
		`c.值`,
	}, "\n")
	mustInt(t, src, 2)
}

func TestClassEmbedUnknownType(t *testing.T) {
	src := "模型 学生\n    嵌入 不存在的模型\n结束"
	mustError(t, src, "未定义")
}

func TestMethodExtensionOverride(t *testing.T) {
	src := strings.Join([]string{
		`模型 点`,
		`    变量 x`,
		`结束`,
		`模型 的 函数 描述 点 p`,
		`    返回 "旧"`,
		`结束`,
		`模型 的 函数 描述 点 p`,
		`    返回 "新"`,
		`结束`,
		`点().描述()`,
	}, "\n")
	res, err := runSource(t, src)
	if err != nil {
		t.Fatalf("运行错误：%v", err)
	}
	if res.Inspect() != "新" {
		t.Fatalf("后定义的方法应覆盖，实际 %s", res.Inspect())
	}
}

func TestMethodExtensionAnywhere(t *testing.T) {
	src := strings.Join([]string{
		`模型 点`,
		`    变量 x`,
		`结束`,
		`变量 p = 点()`,
		`模型 的 函数 获取 点 d`,
		`    返回 d.x`,
		`结束`,
		`p.获取()`,
	}, "\n")
	res, err := runSource(t, src)
	if err != nil {
		t.Fatalf("运行错误：%v", err)
	}
	if res.Inspect() != "空" {
		t.Fatalf("未赋值字段应为空，实际 %s", res.Inspect())
	}
}

func TestModelTypedField(t *testing.T) {
	src := strings.Join([]string{
		`模型 矩形`,
		`    变量 整数 宽`,
		`    变量 整数 高`,
		`结束`,
		`模型 的 函数 初始化 矩形 r、整数 宽、整数 高`,
		`    r.宽 = 宽`,
		`    r.高 = 高`,
		`结束`,
		`变量 r = 矩形(3, 4)`,
		`r.宽`,
	}, "\n")
	mustInt(t, src, 3)
}

func TestModelKeywordAlias(t *testing.T) {
	src := "模型 点\n    变量 x\n结束\n点().x"
	res, err := runSource(t, src)
	if err != nil {
		t.Fatalf("运行错误：%v", err)
	}
	if res.Inspect() != "空" {
		t.Fatalf("模型字段默认为空，实际 %s", res.Inspect())
	}
}
