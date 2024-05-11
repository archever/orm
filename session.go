package orm

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
)

type Session struct {
	db ExecutorIfc
}

func (s *Session) Table(schema Schema) *Action {
	return &Action{
		session: s,
		schema:  schema,
	}
}

var payloadIfcType = reflect.TypeOf((*PayloadIfc)(nil)).Elem()

func (s *Session) queryPayload(ctx context.Context, stmt *Stmt, payloadRef PayloadIfc, nestedPayloadRef ...any) error {
	// TODO: 自动识别 payload 嵌套, 或者使用 nestPayloadRef 指定
	bindFields := boundFields(payloadRef)
	for _, item := range nestedPayloadRef {
		itemV := reflect.ValueOf(item)
		if itemV.Type().Kind() != reflect.Ptr {
			return fmt.Errorf("payload must be pointer")
		}
		if itemV.Type().Implements(payloadIfcType) {
			p := item.(PayloadIfc)
			bindFields = append(bindFields, boundFields(p)...)
		} else if itemV.Type().Elem().Kind() == reflect.Ptr &&
			itemV.Type().Elem().Implements(payloadIfcType) {
			if itemV.Elem().IsNil() && itemV.Elem().CanSet() {
				newItem := reflect.New(itemV.Type().Elem().Elem())
				itemV.Elem().Set(newItem)
			}
			itemDef := itemV.Elem().Interface()
			p := itemDef.(PayloadIfc)
			bindFields = append(bindFields, boundFields(p)...)
		}
	}
	fields := []FieldIfc{}
	for _, field := range bindFields {
		fields = append(fields, field.field)
	}
	stmt.selectField = fields
	expr, err := stmt.completeSelect()
	if err != nil {
		return err
	}
	sqlRaw, argsRaw := expr.Expr()
	rows, err := s.db.QueryContext(ctx, sqlRaw, argsRaw...)
	if err != nil {
		return err
	}
	values := make([]any, 0, len(bindFields))
	for _, field := range bindFields {
		values = append(values, field.RefVal())
	}
	for rows.Next() {
		if err := rows.Scan(values...); err != nil {
			return err
		}
	}
	for _, field := range bindFields {
		field.setPreVal(field.Val())
	}
	return nil
}

func (s *Session) queryPayloadSlice(ctx context.Context, stmt *Stmt, payloadSliceRef any) error {
	rv := reflect.ValueOf(payloadSliceRef)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("must be ptr, find :%T", rv.Interface())
	}
	rvElem := rv.Elem()
	if rvElem.Kind() != reflect.Slice {
		return fmt.Errorf("must be slice, find :%T", rvElem.Interface())
	}
	if rvElem.IsNil() {
		rv.Set(reflect.MakeSlice(rv.Type(), 0, 0).Addr())
	}
	newPayload := reflect.New(rvElem.Type().Elem().Elem())
	p, ok := newPayload.Interface().(PayloadIfc)
	if !ok {
		return fmt.Errorf("must be PayloadIfc, find :%T", newPayload.Interface())
	}
	bindFields := boundFields(p)
	fields := []FieldIfc{}
	for _, field := range bindFields {
		fields = append(fields, field.field)
	}
	stmt.selectField = fields
	expr, err := stmt.completeSelect()
	if err != nil {
		return err
	}
	sqlRaw, argsRaw := expr.Expr()
	rows, err := s.db.QueryContext(ctx, sqlRaw, argsRaw...)
	if err != nil {
		return err
	}
	for rows.Next() {
		rvPayload := reflect.New(rvElem.Type().Elem().Elem())
		p, ok := rvPayload.Interface().(PayloadIfc)
		if !ok {
			return fmt.Errorf("must be PayloadIfc, find :%T", rvPayload.Interface())
		}
		bindFields := boundFields(p)
		values := make([]any, 0, len(bindFields))
		for _, field := range bindFields {
			values = append(values, field.RefVal())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		for _, field := range bindFields {
			field.setPreVal(field.Val())
		}
		rvElem.Set(reflect.Append(rvElem, rvPayload))
	}
	return nil
}

func (s *Session) exec(ctx context.Context, stmt *Stmt) (sql.Result, error) {
	expr, err := stmt.complete()
	if err != nil {
		return nil, err
	}
	sqlRaw, argsRaw := expr.Expr()
	return s.db.ExecContext(ctx, sqlRaw, argsRaw...)
}
