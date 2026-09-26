package object

import (
	"sort"
	"strings"
)

const KindModule Kind = "模块"

type Module struct {
	Name    string
	Exports map[string]Object
}

func NewModule(name string) *Module {
	return &Module{Name: name, Exports: map[string]Object{}}
}

func (m *Module) Kind() Kind { return KindModule }

func (m *Module) Inspect() string {
	names := make([]string, 0, len(m.Exports))
	for k := range m.Exports {
		names = append(names, k)
	}
	sort.Strings(names)
	return "<模块 " + m.Name + "：" + strings.Join(names, "、") + ">"
}
