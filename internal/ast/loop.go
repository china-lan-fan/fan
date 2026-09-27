package ast

import "fmt"

type WhileStmt struct {
	Position  Position
	Condition Expression
	Body      *BlockStmt
}

func (n *WhileStmt) Pos() Position  { return n.Position }
func (n *WhileStmt) statementNode() {}
func (n *WhileStmt) String() string {
	return fmt.Sprintf("当 %s 循环 { %s } 结束", n.Condition.String(), n.Body.String())
}

type RepeatStmt struct {
	Position Position
	Body     *BlockStmt
	Until    Expression
}

func (n *RepeatStmt) Pos() Position  { return n.Position }
func (n *RepeatStmt) statementNode() {}
func (n *RepeatStmt) String() string {
	return fmt.Sprintf("重复 { %s } 直到 %s 结束", n.Body.String(), n.Until.String())
}

type BreakStmt struct {
	Position Position
}

func (n *BreakStmt) Pos() Position  { return n.Position }
func (n *BreakStmt) statementNode() {}
func (n *BreakStmt) String() string { return "跳出" }

type ContinueStmt struct {
	Position Position
}

func (n *ContinueStmt) Pos() Position  { return n.Position }
func (n *ContinueStmt) statementNode() {}
func (n *ContinueStmt) String() string { return "继续" }

type TimesStmt struct {
	Position Position
	Count    Expression
	Body     *BlockStmt
}

func (n *TimesStmt) Pos() Position  { return n.Position }
func (n *TimesStmt) statementNode() {}
func (n *TimesStmt) String() string {
	return fmt.Sprintf("重复 %s 次 { %s } 结束", n.Count.String(), n.Body.String())
}

type ForEachStmt struct {
	Position   Position
	Iterable   Expression
	FirstName  string
	SecondName string
	Body       *BlockStmt
}

func (n *ForEachStmt) Pos() Position  { return n.Position }
func (n *ForEachStmt) statementNode() {}
func (n *ForEachStmt) String() string {
	if n.SecondName == "" {
		return fmt.Sprintf("遍历 %s 中的 %s { %s } 结束", n.Iterable.String(), n.FirstName, n.Body.String())
	}
	return fmt.Sprintf("遍历 %s 中的 %s %s { %s } 结束", n.Iterable.String(), n.FirstName, n.SecondName, n.Body.String())
}
