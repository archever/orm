package orm

type Schema interface {
	TableName() string
	IDField() FieldIfc
}
