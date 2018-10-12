// interface summary
// orm -> table -> action -> do
// orm -> exec -> action -> do
// tx -> TableTx -> actiontx -> do
// tx -> TableTx -> actiontx -> do -> commit
// tx -> exec -> actiontx -> do -> rollback/commit

package orm

// OrmI interface to build sql
type OrmI interface {
	Table(t string) TabledI
	Exec(sql string, arg ...interface{}) ActionI
}

type TabledI interface {
	Select(field ...string) ActionI
	Update(data map[string]interface{}) ActionI
	Insert(data ...interface{}) ActionI
	Replace(data ...interface{}) ActionI
	Delete() ActionI
}

// OrmI interface to build sql
type TxI interface {
	Table(t string) TabledTxI
	Exec(sql string, arg ...interface{}) ActionTxI
	Begin() error
	Commit() error
	RollBack() error
}

type TabledTxI interface {
	Select(field ...string) ActionTxI
	Update(data map[string]interface{}) ActionTxI
	Insert(data ...interface{}) ActionTxI
	Replace(data ...interface{}) ActionTxI
	Delete() ActionTxI
}

// ActionI sql executor interface
type ActionI interface {
	Get(dest interface{}) error
	One(dest interface{}) error
	Do() (int64, int64, error)
	Filter(f ...*FilterItem) ActionI
	Where(string, ...interface{}) ActionI
	OrderBy(o ...string) ActionI
	GroupBy(o ...string) ActionI
	Limit(l int64) ActionI
	Offset(o int64) ActionI
	Page(page, psize int64) ActionI
	SQL() (string, []interface{}, error)
}

// ActionTxI sql executor interface
type ActionTxI interface {
	Get(dest interface{}) error
	One(dest interface{}) error
	Do() (int64, int64, error)
	Filter(f ...*FilterItem) ActionTxI
	Where(string, ...interface{}) ActionTxI
	OrderBy(o ...string) ActionTxI
	GroupBy(o ...string) ActionTxI
	Limit(l int64) ActionTxI
	Offset(o int64) ActionTxI
	Page(page, psize int64) ActionTxI
	SQL() (string, []interface{}, error)
}
