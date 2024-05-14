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
var _ ExprIfc = (*groupBy)(nil)
var _ ExprIfc = (*where)(nil)
var _ ExprIfc = (*selectExpr)(nil)
var _ ExprIfc = (*updateExpr)(nil)
var _ ExprIfc = (*deleteExpr)(nil)
var _ ExprIfc = (*insertExpr)(nil)
var _ ExprIfc = (*ExprSlice)(nil)
var _ ExprIfc = (*anyVal)(nil)
var _ ExprIfc = (*anyValList)(nil)
var _ ExprIfc = (*joinExpr)(nil)
var _ ExprIfc = (*brackets)(nil)
var _ ExprIfc = (*fields)(nil)

type Cond struct {
	left  ExprIfc
	right ExprIfc
	Op    string
}

func (c *Cond) Expr() (expr string, args []any) {
	switch c.Op {
	case "":
		return c.left.Expr()
	case "IS NULL", "IS NOT NULL":
		leftE, leftA := c.left.Expr()
		expr = leftE + " " + c.Op
		args = append(args, leftA...)
	case "IN", "NOT IN":
		fallthrough
	default:
		leftE, leftA := c.left.Expr()
		rightE, rightA := c.right.Expr()
		expr = fmt.Sprintf("%s %s %s", leftE, c.Op, rightE)
		args = append(args, leftA...)
		args = append(args, rightA...)
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

func And(cond ...Cond) Cond {
	return Cond{
		left: &groupExpr{
			op:    "AND",
			conds: cond,
		},
	}
}

func Or(cond ...Cond) Cond {
	return Cond{
		left: &groupExpr{
			op:    "OR",
			conds: cond,
		},
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
	expr = a.Field.DBColName(true)
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
	fields        []FieldIfc
	schema        Schema
	withTableName bool
}

func (a selectExpr) Expr() (expr string, args []any) {
	fields := []string{}
	for _, field := range a.fields {
		if a.withTableName {
			fields = append(fields, field.DBColName(true))
		} else {
			fields = append(fields, field.ColName(true))
		}
	}
	expr = fmt.Sprintf("SELECT %s FROM %s", strings.Join(fields, ", "), FieldWrapper(a.schema.TableName()))
	return
}

type fields struct {
	fields        []FieldIfc
	withTableName bool
}

func (a fields) Expr() (expr string, args []any) {
	fields := []string{}
	for _, field := range a.fields {
		if a.withTableName {
			fields = append(fields, field.DBColName(true))
		} else {
			fields = append(fields, field.ColName(true))
		}
	}
	expr = strings.Join(fields, ", ")
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
		fields = append(fields, field.ColName(true))
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
		cond = &groupExpr{
			op:    "AND",
			conds: a,
		}
	}
	e, ar := cond.Expr()
	expr = "WHERE " + e
	args = ar
	return
}

type ExprSlice []ExprIfc

func (a ExprSlice) Expr() (expr string, args []any) {
	sb := strings.Builder{}
	length := len(a)
	for i, e := range a {
		e, a := e.Expr()
		if e != "" {
			sb.WriteString(e)
			if i < length-1 {
				sb.WriteString(" ")
			}
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

type insertExpr struct {
	vales  [][]*fieldBind
	fields []FieldIfc
	schema Schema
}

func (e *insertExpr) Expr() (expr string, args []any) {
	fields := []string{}
	for _, field := range e.fields {
		fields = append(fields, field.ColName(true))
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES", FieldWrapper(e.schema.TableName()), strings.Join(fields, ",")))
	rows := []string{}
	for _, row := range e.vales {
		values := []string{}
		for _, field := range row {
			values = append(values, "?")
			args = append(args, field.Val())
		}
		rows = append(rows, "("+strings.Join(values, ",")+")")
	}
	sb.WriteString(strings.Join(rows, ","))
	expr = sb.String()
	return
}

type joinExpr struct {
	tp     string
	on     []Cond
	schema Schema
}

func (a *joinExpr) Expr() (expr string, args []any) {
	joinStr := "JOIN"
	if a.tp != "" {
		joinStr = fmt.Sprintf("%s %s", a.tp, joinStr)
	}
	expr = fmt.Sprintf("%s %s ON", joinStr, FieldWrapper(a.schema.TableName()))
	for _, cond := range a.on {
		e, a := cond.Expr()
		expr += " " + e
		args = append(args, a...)
	}
	return
}

type anyVal struct {
	any
}

func (a anyVal) Expr() (expr string, args []any) {
	expr = "?"
	args = []any{a.any}
	return
}

type anyValList []any

func (a anyValList) Expr() (expr string, args []any) {
	expr = strings.Repeat("?,", len(a))
	expr = expr[:len(expr)-1]
	args = a
	return
}

type brackets struct {
	ExprIfc
}

func (b brackets) Expr() (expr string, args []any) {
	e, a := b.ExprIfc.Expr()
	expr = "(" + e + ")"
	args = a
	return
}
