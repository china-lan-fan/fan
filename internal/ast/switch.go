package ast

import "fmt"

type SwitchBranch struct {
	Position Position
	Guard    bool
	Values   []Expression
	Body     *BlockStmt
	Default  bool
}

type SwitchStmt struct {
	Position Position
	Subject  Expression
	Branches []SwitchBranch
}

func (n *SwitchStmt) Pos() Position  { return n.Position }
func (n *SwitchStmt) statementNode() {}
func (n *SwitchStmt) String() string {
	return fmt.Sprintf("判断 %s ... 结束", n.Subject)
}
