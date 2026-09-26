package ast

import (
	"fmt"
	"strings"
)

type ArrayLiteral struct {
	Position Position
	Elements []Expression
}

func (n *ArrayLiteral) Pos() Position   { return n.Position }
func (n *ArrayLiteral) expressionNode() {}
func (n *ArrayLiteral) String() string {
	parts := make([]string, 0, len(n.Elements))
	for _, e := range n.Elements {
		parts = append(parts, e.String())
	}
	return "[" + strings.Join(parts, ", ") + "]"
}

type IndexExpr struct {
	Position Position
	Target   Expression
	Index    Expression
}

func (n *IndexExpr) Pos() Position   { return n.Position }
func (n *IndexExpr) expressionNode() {}
func (n *IndexExpr) String() string {
	return fmt.Sprintf("%s[%s]", n.Target.String(), n.Index.String())
}

type CallExpr struct {
	Position Position
	Callee   Expression
	Args     []Expression
}

func (n *CallExpr) Pos() Position   { return n.Position }
func (n *CallExpr) expressionNode() {}
func (n *CallExpr) String() string {
	parts := make([]string, 0, len(n.Args))
	for _, a := range n.Args {
		parts = append(parts, a.String())
	}
	return fmt.Sprintf("%s(%s)", n.Callee.String(), strings.Join(parts, ", "))
}

type IndexAssignStmt struct {
	Position Position
	Target   Expression
	Index    Expression
	Value    Expression
}

func (n *IndexAssignStmt) Pos() Position  { return n.Position }
func (n *IndexAssignStmt) statementNode() {}
func (n *IndexAssignStmt) String() string {
	return fmt.Sprintf("%s[%s] = %s", n.Target.String(), n.Index.String(), n.Value.String())
}

type DictPair struct {
	Key   Expression
	Value Expression
}

type DictLiteral struct {
	Position Position
	Pairs    []DictPair
}

func (n *DictLiteral) Pos() Position   { return n.Position }
func (n *DictLiteral) expressionNode() {}
func (n *DictLiteral) String() string {
	parts := make([]string, 0, len(n.Pairs))
	for _, p := range n.Pairs {
		parts = append(parts, p.Key.String()+": "+p.Value.String())
	}
	return "{" + strings.Join(parts, ", ") + "}"
}
