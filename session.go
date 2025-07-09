package sqlite_transaction

import (
	"context"
	"database/sql"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type Database interface {
	From(cols ...interface{}) *goqu.SelectDataset
	Select(cols ...interface{}) *goqu.SelectDataset
	Update(table interface{}) *goqu.UpdateDataset
	Insert(table interface{}) *goqu.InsertDataset
	Delete(table interface{}) *goqu.DeleteDataset
	Truncate(table ...interface{}) *goqu.TruncateDataset
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ScanStructs(i interface{}, query string, args ...interface{}) error
	ScanStructsContext(ctx context.Context, i interface{}, query string, args ...interface{}) error
	ScanStruct(i interface{}, query string, args ...interface{}) (bool, error)
	ScanStructContext(ctx context.Context, i interface{}, query string, args ...interface{}) (bool, error)
	ScanVals(i interface{}, query string, args ...interface{}) error
	ScanValsContext(ctx context.Context, i interface{}, query string, args ...interface{}) error
	ScanVal(i interface{}, query string, args ...interface{}) (bool, error)
	ScanValContext(ctx context.Context, i interface{}, query string, args ...interface{}) (bool, error)
}

// Manager transaction manager for application.
type Manager interface {
	// Begin - create transaction, if exist transaction return error
	Begin(inCtx context.Context) (outCtx context.Context, tx Transaction, err error)
	// Wrap -  automatically COMMIT or ROLLBACK
	Wrap(ctx context.Context, tFunc func(ctx context.Context) error) error
}

// Session database manager for repository.
type Session interface {
	Use(ctx context.Context) Database
}

type Connection interface {
	Manager
	Session
	RegSequence(table exp.IdentifierExpression) error
	GetSequence(table exp.IdentifierExpression) (*Sequence, error)
	DB() *sql.DB
}
