package orm

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

// func (o *Action) Insert(row interface{}) *Stmt {
// 	sql, args, err := sqlInsert(o.table, row)
// 	if err != nil {
// 		return o.errStmt(err)
// 	}
// 	return o.passStmt(sql, args...)
// }

// func (o *Action) InsertMany(rows interface{}) *Stmt {
// 	sql, args, err := sqlInsertMany(o.table, rows)
// 	if err != nil {
// 		return o.errStmt(err)
// 	}
// 	return o.passStmt(sql, args...)
// }

// func (o *Action) Replace(row interface{}) *Stmt {
// 	sql, args, err := sqlReplace(o.table, row)
// 	if err != nil {
// 		return o.errStmt(err)
// 	}
// 	return o.passStmt(sql, args...)
// }

// func (o *Action) ReplaceMany(rows interface{}) *Stmt {
// 	sql, args, err := sqlReplaceMany(o.table, rows)
// 	if err != nil {
// 		return o.errStmt(err)
// 	}
// 	return o.passStmt(sql, args...)
// }

func (o *Action) Delete() *Stmt {
	stm := &Stmt{
		session: o.session,
		schema:  o.schema,
	}
	stm.completeFn = stm.completeDelete
	return stm
}
