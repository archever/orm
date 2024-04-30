package orm

import (
	"context"
	"errors"
	"strings"
)

type Stmt struct {
	session *Session
	*selectExpr
	*updateExpr
	schema Schema

	conds   []Cond
	orderBy orderBy
	groupBy groupBy
	limit   *limit
	offset  *offset

	err error
}

func (a *Stmt) complete() (rawSQL string, argList []any, err error) {
	if a.err != nil {
		err = a.err
		return
	}
	parts := []ExprIfc{}
	if a.selectExpr != nil {
		parts = append(parts, a.selectExpr)
	}
	if a.updateExpr != nil {
		parts = append(parts, a.updateExpr)
	}

	if len(a.conds) > 0 {
		parts = append(parts, Where(a.conds...))
	}

	if len(a.groupBy) > 0 {
		parts = append(parts, a.groupBy)
	}
	if len(a.orderBy) > 0 {
		parts = append(parts, a.orderBy)
	}
	if a.offset != nil {
		parts = append(parts, a.offset)
	}
	if a.limit != nil {
		parts = append(parts, a.limit)
	}
	sb := strings.Builder{}
	for i := range parts {
		part := parts[i]
		expr, arg := part.Expr()
		sb.WriteString(expr)
		sb.WriteString(" ")
		argList = append(argList, arg...)
	}
	rawSQL = sb.String()
	return
}

// SQL return the sql info for testing
func (a *Stmt) SQL() (string, []any, error) {
	return a.complete()
}

// Where generate where condition
func (a *Stmt) Where(cond ...Cond) *Stmt {
	a.conds = append(a.conds, cond...)
	return a
}

// OrderBy set sql order by
func (a *Stmt) OrderBy(order ...Order) *Stmt {
	a.orderBy = append(a.orderBy, order...)
	return a
}

// GroupBy set sql group by
func (a *Stmt) GroupBy(field ...FieldIfc) *Stmt {
	a.groupBy = append(a.groupBy, field...)
	return a
}

// Limit set sql limit
func (a *Stmt) Limit(l int64) *Stmt {
	a.limit = (*limit)(&l)
	return a
}

// Offset set sql offset
func (a *Stmt) Offset(o int64) *Stmt {
	a.offset = (*offset)(&o)
	return a
}

// Page set sql limit and offset
func (a *Stmt) Page(page, size int64) *Stmt {
	if page < 1 {
		page = 1
	}
	_offset := (page - 1) * size
	a.limit = (*limit)(&size)
	a.offset = (*offset)(&_offset)
	return a
}

func (a *Stmt) Set(cond ...Cond) *Stmt {
	if a.updateExpr == nil {
		a.err = errors.New("set is only for update")
		return a
	}
	// TODO 去重
	a.updateExpr.sets = append(a.updateExpr.sets, cond...)
	return a
}

func (a *Stmt) TakePayload(ctx context.Context, payload PayloadIfc) error {
	a.limit = new(limit)
	*a.limit = 1
	return a.session.queryPayload(ctx, a, payload)
}

func (a *Stmt) Do(ctx context.Context) error {
	_, err := a.session.exec(ctx, a)
	return err
}
