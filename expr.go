package orm

import (
	"fmt"
	"strings"
)

type ExprIfc interface {
	Expr() (string, []any)
}

var _ ExprIfc = (*Cond)(nil)
var _ ExprIfc = (*groupExpr)(nil)
var _ ExprIfc = (*Order)(nil)
var _ ExprIfc = (*orderBy)(nil)
var _ ExprIfc = (*limit)(nil)
var _ ExprIfc = (*offset)(nil)
var _ ExprIfc = (*selectExpr)(nil)
var _ ExprIfc = (*updateExpr)(nil)
var _ ExprIfc = (*groupBy)(nil)
var _ ExprIfc = (*where)(nil)

type Cond struct {
	left       FieldIfc
	rightVal   any
	rightField FieldIfc
	Op         string
}

func (c *Cond) Expr() (expr string, args []any) {
	if c.rightField != nil {
		expr = FieldWrapper(c.left.ColName()) + c.Op + c.rightField.ColName()
	} else {
		expr = FieldWrapper(c.left.ColName()) + c.Op + "?"
		args = []any{c.rightVal}
	}
	return
}

type groupExpr struct {
	op    string
	conds []Cond
}

func (a *groupExpr) Expr() (expr string, args []any) {
	condExpr := []string{}
	for _, cond := range a.conds {
		e, a := cond.Expr()
		condExpr = append(condExpr, e)
		args = append(args, a...)
	}
	expr = "(" + strings.Join(condExpr, " "+a.op+" ") + ")"
	return
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

func OrderBy(order ...Order) ExprIfc {
	return orderBy(order)
}

type orderBy []Order

func (a orderBy) Expr() (expr string, args []any) {
	if len(a) == 0 {
		return
	}
	orders := []string{}
	for _, o := range a {
		e, a := o.Expr()
		orders = append(orders, e)
		args = append(args, a...)
	}
	expr = "order by " + strings.Join(orders, ", ")
	return
}

type Order struct {
	Field FieldIfc
	Desc  bool
}

func (a *Order) Expr() (expr string, args []any) {
	expr = FieldWrapper(a.Field.ColName())
	if a.Desc {
		expr += " DESC"
	}
	return
}

type limit int64
type offset int64

func (a *limit) Expr() (expr string, args []any) {
	expr = "LIMIT ?"
	args = []any{*a}
	return
}

func (a *offset) Expr() (expr string, args []any) {
	if *a == 0 {
		return
	}
	expr = "OFFSET ?"
	args = []any{*a}
	return
}

type selectExpr struct {
	fields []FieldIfc
	schema Schema
}

func (a selectExpr) Expr() (expr string, args []any) {
	fields := []string{}
	for _, field := range a.fields {
		fields = append(fields, FieldWrapper(field.ColName()))
	}
	expr = fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), FieldWrapper(a.schema.TableName()))
	return
}

func Select(table Schema, fields ...FieldIfc) ExprIfc {
	return selectExpr{
		fields: fields,
		schema: table,
	}
}

type updateExpr struct {
	sets   []Cond
	schema Schema
}

func (e *updateExpr) Expr() (expr string, args []any) {
	set := []string{}
	for _, field := range e.sets {
		s, a := field.Expr()
		set = append(set, s)
		args = append(args, a...)
	}
	expr = fmt.Sprintf("UPDATE %s SET %s", FieldWrapper(e.schema.TableName()), strings.Join(set, ", "))
	return
}

type groupBy []FieldIfc

func (gb groupBy) Expr() (expr string, args []any) {
	fields := []string{}
	for _, field := range gb {
		fields = append(fields, FieldWrapper(field.ColName()))
	}
	expr = "GROUP BY " + strings.Join(fields, ", ")
	return
}

type where []Cond

func Where(cond ...Cond) ExprIfc {
	return where(cond)
}

func (a where) Expr() (expr string, args []any) {
	if len(a) == 0 {
		return
	}
	var cond ExprIfc = &a[0]
	if len(a) > 1 {
		cond = And(a...)
	}
	e, ar := cond.Expr()
	expr = "WHERE " + e
	args = ar
	return
}

type ExprSlice []ExprIfc

func (a ExprSlice) Expr() (expr string, args []any) {
	sb := strings.Builder{}
	for _, e := range a {
		e, a := e.Expr()
		if e != "" {
			sb.WriteString(e)
			sb.WriteString(" ")
		}
		args = append(args, a...)
	}
	return sb.String(), args
}

type deleteExpr struct {
	schema Schema
}

func (e *deleteExpr) Expr() (expr string, args []any) {
	expr = fmt.Sprintf("DELETE FROM %s", FieldWrapper(e.schema.TableName()))
	return
}
