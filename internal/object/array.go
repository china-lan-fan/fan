package object

import "strings"

const KindArray Kind = "数组"

type Array struct {
	Elements []Object
}

func (a *Array) Kind() Kind { return KindArray }

func (a *Array) Inspect() string {
	parts := make([]string, 0, len(a.Elements))
	for _, e := range a.Elements {
		parts = append(parts, Format(e))
	}
	return "[" + strings.Join(parts, ", ") + "]"
}
