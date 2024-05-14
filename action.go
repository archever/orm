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
	bindFields := boundFields(payload)
	for i := range bindFields {
		item := bindFields[i]
		if item.Dirty() {
			stm.sets = append(stm.sets, Cond{
				left:  item.field,
				Op:    "=",
				right: anyVal{item.Val()},
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

func (o *Action) InsertPayload(rows ...PayloadIfc) *Stmt {
	if len(rows) == 0 {
		return &Stmt{err: errors.New("no payload")}
	}
	values := [][]*fieldBind{}
	autoIncrementFields := []*fieldBind{}
	for i := range rows {
		row := rows[i]
		bindFields := boundFields(row)
		notIgnoredFields := []*fieldBind{}
		for j := range bindFields {
			if bindFields[j].field.IsAutoIncrement() {
				autoIncrementFields = append(autoIncrementFields, bindFields[j])
				continue
			}
			notIgnoredFields = append(notIgnoredFields, bindFields[j])
		}
		values = append(values, notIgnoredFields)
	}
	fields := []FieldIfc{}
	for i := range values[0] {
		fields = append(fields, values[0][i].field)
	}
	stm := &Stmt{
		session:     o.session,
		schema:      o.schema,
		selectField: fields,
		values:      values,
	}
	stm.completeFn = stm.completeInsert
	stm.afterExecFn = func(row, id int64) {
		for i := range autoIncrementFields {
			f := autoIncrementFields[i]
			f.Set(id + int64(i))
		}
	}
	return stm
}
