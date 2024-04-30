package orm

import "strings"

type ExprIfc interface {
	Expr() string
	Args() []any
}

type Cond struct {
	left       FieldIfc
	rightVal   any
	rightField FieldIfc
	Op         string
}

func (c *Cond) Expr() string {
	if c.rightField != nil {
		return FieldWapper(c.left.ColName()) + c.Op + c.rightField.ColName()
	}
	return FieldWapper(c.left.ColName()) + c.Op + "?"
}

func (c *Cond) Args() []any {
	if c.rightField == nil {
		return []any{c.rightVal}
	}
	return []any{}
}

type groupExpr struct {
	op    string
	conds []Cond
}

func (a *groupExpr) Expr() string {
	condExpr := []string{}
	for _, cond := range a.conds {
		condExpr = append(condExpr, cond.Expr())
	}
	return "(" + strings.Join(condExpr, " "+a.op+" ") + ")"
}

func (a *groupExpr) Args() []any {
	args := []any{}
	for _, cond := range a.conds {
		args = append(args, cond.Args()...)
	}
	return args
}

func And(cond ...Cond) ExprIfc {
	return &groupExpr{
		op:    "AND",
		conds: cond,
	}
}

func Or(cond ...Cond) ExprIfc {
	return &groupExpr{
		op:    "OR",
		conds: cond,
	}
}
