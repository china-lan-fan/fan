package ast

import "fmt"

type ImportStmt struct {
	Position Position
	Path     string
	Name     string
}

func (n *ImportStmt) Pos() Position  { return n.Position }
func (n *ImportStmt) statementNode() {}
func (n *ImportStmt) String() string {
	return fmt.Sprintf("导入 %q 作为 %s", n.Path, n.Name)
}

type ExportStmt struct {
	Position Position
	Name     string
	Inner    Statement
}

func (n *ExportStmt) Pos() Position  { return n.Position }
func (n *ExportStmt) statementNode() {}
func (n *ExportStmt) String() string {
	return fmt.Sprintf("导出 %s", n.Inner.String())
}
