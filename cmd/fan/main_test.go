package main

import "testing"

func TestRunNoArgsShowsUsage(t *testing.T) {
	if err := run(nil); err != nil {
		t.Fatalf("空参数应打印帮助而非报错，得到：%v", err)
	}
}

func TestRunVersion(t *testing.T) {
	if err := run([]string{"version"}); err != nil {
		t.Fatalf("version 子命令不应报错：%v", err)
	}
}

func TestRunUnknown(t *testing.T) {
	if err := run([]string{"未知"}); err == nil {
		t.Fatal("未知子命令应返回错误")
	}
}

func TestRunRequiresScript(t *testing.T) {
	if err := run([]string{"run"}); err == nil {
		t.Fatal("run 缺少脚本参数应返回错误")
	}
}
