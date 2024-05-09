package orm

import "errors"

type Action struct {
	session *Session
	schema  Schema
}

func (o *Action) Select(field ...FieldIfc) *Stmt {
	stm := &Stmt{
		session:     o.session,
		schema:      o.schema,
		selectField: field,
	}
	stm.completeFn = stm.completeSelect
	return stm
}

func (o *Action) UpdatePayload(payload PayloadIfc) *Stmt {
	stm := &Stmt{
		session: o.session,
		schema:  o.schema,
	}
	stm.completeFn = stm.completeUpdate
	payload.Bind()
	fields := payload.Fields()
	for i := range fields {
		if fields[i].Dirty() {
			stm.sets = append(stm.sets, Cond{
				left:     fields[i],
				Op:       "=",
				rightVal: fields[i].Val(),
			})
		}
	}
	return stm
}

func (o *Action) Update(cond ...Cond) *Stmt {
	stm := &Stmt{
		session: o.session,
		schema:  o.schema,
		sets:    cond,
	}
	stm.completeFn = stm.completeUpdate
	return stm
}

func (o *Action) Delete() *Stmt {
	stm := &Stmt{
		session: o.session,
		schema:  o.schema,
	}
	stm.completeFn = stm.completeDelete
	return stm
}

func (o *Action) InsertPayload(row ...PayloadIfc) *Stmt {
	if len(row) == 0 {
		return &Stmt{err: errors.New("no payload")}
	}
	values := [][]FieldIfc{}
	autoIncrementFields := []FieldIfc{}
	for i := range row {
		row[i].Bind()
		fields := row[i].Fields()
		ignoredFields := []FieldIfc{}
		for j := range fields {
			if fields[j].AutoIncrement() {
				autoIncrementFields = append(autoIncrementFields, fields[j])
				continue
			}
			ignoredFields = append(ignoredFields, fields[j])
		}
		values = append(values, ignoredFields)
	}
	stm := &Stmt{
		session:     o.session,
		schema:      o.schema,
		selectField: values[0],
		values:      values,
	}
	stm.completeFn = stm.completeInsert
	stm.afterExecFn = func(row, id int64) {
		for i := range autoIncrementFields {
			f := autoIncrementFields[i]
			f.set(id + int64(i))
		}
	}
	return stm
}
