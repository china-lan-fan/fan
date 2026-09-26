package evaluator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fan/internal/ast"
	"fan/internal/object"
	"fan/internal/parser"
)

type Loader struct {
	cache      map[string]*object.Module
	loading    map[string]bool
	modulesDir string
}

func NewLoader(modulesDir string) *Loader {
	return &Loader{
		cache:      map[string]*object.Module{},
		loading:    map[string]bool{},
		modulesDir: modulesDir,
	}
}

func (l *Loader) resolvePath(base string, path string) (string, error) {
	if strings.HasSuffix(path, ".fan") || filepath.IsAbs(path) {
		return filepath.Abs(path)
	}
	if base != "" {
		candidate := filepath.Join(base, path+".fan")
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Abs(candidate)
		}
	}
	if l.modulesDir != "" {
		candidate := filepath.Join(l.modulesDir, path+".fan")
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Abs(candidate)
		}
	}
	if !strings.HasSuffix(path, ".fan") {
		candidate := path + ".fan"
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Abs(candidate)
		}
	}
	return "", fmt.Errorf("找不到模块：%s", path)
}

func (l *Loader) Load(importerBase string, path string) (*object.Module, error) {
	full, err := l.resolvePath(importerBase, path)
	if err != nil {
		return nil, err
	}
	if mod, ok := l.cache[full]; ok {
		return mod, nil
	}
	if l.loading[full] {
		return nil, fmt.Errorf("检测到循环导入：%s", path)
	}
	l.loading[full] = true
	defer delete(l.loading, full)

	data, err := os.ReadFile(full)
	if err != nil {
		return nil, err
	}
	prog, errs := parser.ParseProgram(string(data))
	if len(errs) > 0 {
		return nil, fmt.Errorf("模块 %s 解析失败：%s", path, errs[0].Error())
	}

	base := filepath.Dir(full)
	env := NewEnvironment()
	env.BaseDir = base
	env.Loader = l
	env.moduleExport = true

	if _, err := Eval(prog, env); err != nil {
		return nil, err
	}

	mod := object.NewModule(moduleNameFromPath(full))
	for name, b := range env.store {
		if b.isExport {
			mod.Exports[name] = b.value
		}
	}
	l.cache[full] = mod
	return mod, nil
}

func moduleNameFromPath(full string) string {
	base := filepath.Base(full)
	if i := strings.LastIndex(base, "."); i > 0 {
		base = base[:i]
	}
	return base
}

func evalImportStmt(stmt *ast.ImportStmt, env *Environment) (object.Object, error) {
	if env.Loader == nil {
		env.Loader = NewLoader(env.BaseDir)
	}
	mod, err := env.Loader.Load(env.BaseDir, stmt.Path)
	if err != nil {
		return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
	}
	if _, exists := env.find(stmt.Name); exists {
		return nil, &EvalError{Pos: stmt.Position, Reason: fmt.Sprintf("模块名 %s 已存在", stmt.Name)}
	}
	if err := env.declare(stmt.Name, mod, true, ast.TypeAny); err != nil {
		return nil, &EvalError{Pos: stmt.Position, Reason: err.Error()}
	}
	return mod, nil
}

func evalExportStmt(stmt *ast.ExportStmt, env *Environment) (object.Object, error) {
	res, err := Eval(stmt.Inner, env)
	if err != nil {
		return nil, err
	}
	if b, ok := env.find(stmt.Name); ok {
		b.isExport = true
	}
	return res, nil
}
