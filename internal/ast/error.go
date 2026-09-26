package ast

import "fmt"

type CheckExpr struct {
	Position Position
	Call     Expression
}

func (n *CheckExpr) Pos() Position   { return n.Position }
func (n *CheckExpr) expressionNode() {}
func (n *CheckExpr) String() string {
	return fmt.Sprintf("检查 %s", n.Call.String())
}

type TryStmt struct {
	Position  Position
	Body      *BlockStmt
	CatchName string
	Catch     *BlockStmt
}

func (n *TryStmt) Pos() Position  { return n.Position }
func (n *TryStmt) statementNode() {}
func (n *TryStmt) String() string {
	name := n.CatchName
	if name == "" {
		name = "_"
	}
	return fmt.Sprintf("尝试 { %s } 捕获 %s { %s } 结束", n.Body.String(), name, n.Catch.String())
}
