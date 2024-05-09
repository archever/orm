package orm

import "context"

type Stmt struct {
	err        error
	session    *Session
	schema     Schema
	completeFn func() (ExprIfc, error)

	joins       []joinExpr
	values      [][]FieldIfc
	conds       []Cond
	orderBy     []Order
	groupBy     []FieldIfc
	sets        []Cond
	selectField []FieldIfc
	limit       *limit
	offset      *offset

	afterExecFn func(row, id int64)
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
	// TODO 去重
	a.sets = append(a.sets, cond...)
	return a
}

func (a *Stmt) Select(field ...FieldIfc) *Stmt {
	// TODO 去重
	a.selectField = append(a.selectField, field...)
	return a
}

func (a *Stmt) Join(s Schema, on ...Cond) *Stmt {
	a.joins = append(a.joins, joinExpr{
		schema: s,
		on:     on,
	})
	return a
}

func (a *Stmt) completeSelect() (ExprIfc, error) {
	action := &selectExpr{
		fields: a.selectField,
		schema: a.schema,
	}
	if len(a.joins) > 0 {
		action.withTableName = true
	}
	exprs := []ExprIfc{action}
	for _, join := range a.joins {
		exprs = append(exprs, &join)
	}
	if len(a.conds) > 0 {
		exprs = append(exprs, Where(a.conds...))
	}
	if len(a.groupBy) > 0 {
		exprs = append(exprs, groupBy(a.groupBy))
	}
	if len(a.orderBy) > 0 {
		exprs = append(exprs, orderBy(a.orderBy))
	}
	if a.limit != nil {
		exprs = append(exprs, a.limit)
	}
	if a.offset != nil {
		exprs = append(exprs, a.offset)
	}
	return ExprSlice(exprs), a.err
}

func (a *Stmt) completeDelete() (ExprIfc, error) {
	action := &deleteExpr{
		schema: a.schema,
	}
	exprs := []ExprIfc{action}
	if len(a.conds) > 0 {
		exprs = append(exprs, Where(a.conds...))
	}
	if a.limit != nil {
		exprs = append(exprs, a.limit)
	}
	if a.offset != nil {
		exprs = append(exprs, a.offset)
	}
	return ExprSlice(exprs), a.err
}

func (a *Stmt) completeInsert() (ExprIfc, error) {
	action := &insertExpr{
		schema: a.schema,
		vales:  a.values,
		fields: a.selectField,
	}
	exprs := []ExprIfc{action}
	return ExprSlice(exprs), a.err
}

func (a *Stmt) completeUpdate() (ExprIfc, error) {
	// TODO: 根据 select 过滤 set
	action := &updateExpr{
		sets:   a.sets,
		schema: a.schema,
	}
	exprs := []ExprIfc{action}
	if len(a.conds) > 0 {
		exprs = append(exprs, Where(a.conds...))
	}
	if len(a.groupBy) > 0 {
		exprs = append(exprs, groupBy(a.groupBy))
	}
	if len(a.orderBy) > 0 {
		exprs = append(exprs, orderBy(a.orderBy))
	}
	if a.limit != nil {
		exprs = append(exprs, a.limit)
	}
	if a.offset != nil {
		exprs = append(exprs, a.offset)
	}
	return ExprSlice(exprs), a.err
}

func (a *Stmt) complete() (ExprIfc, error) {
	if a.completeFn != nil {
		return a.completeFn()
	}
	if a.selectField != nil {
		return a.completeSelect()
	}
	return a.completeUpdate()
}

func (a *Stmt) SubQuery() ExprIfc {
	expr, _ := a.completeSelect()
	return expr
}

func (a *Stmt) TakePayload(ctx context.Context, payload PayloadIfc, nestedPayload ...any) error {
	a.limit = new(limit)
	*a.limit = 1
	return a.session.queryPayload(ctx, a, payload, nestedPayload...)
}

func (a *Stmt) FindPayload(ctx context.Context, payloadsRef any) error {
	return a.session.queryPayloadSlice(ctx, a, payloadsRef)
}

func (a *Stmt) Do(ctx context.Context) (rowCnt int64, err error) {
	ret, err := a.session.exec(ctx, a)
	if err != nil {
		return
	}
	id, err := ret.LastInsertId()
	if err != nil {
		return
	}
	rowCnt, err = ret.RowsAffected()
	if err != nil {
		return
	}
	if a.afterExecFn != nil {
		a.afterExecFn(rowCnt, id)
	}
	return
}
