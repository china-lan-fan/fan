package object

import "strings"

const KindTuple Kind = "多元值"

type Tuple struct {
	Values []Object
}

func (t *Tuple) Kind() Kind { return KindTuple }
func (t *Tuple) Inspect() string {
	parts := make([]string, len(t.Values))
	for i, v := range t.Values {
		if v == nil {
			parts[i] = "空"
			continue
		}
		parts[i] = v.Inspect()
	}
	return "(" + strings.Join(parts, ", ") + ")"
}
