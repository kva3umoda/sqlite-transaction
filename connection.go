package sqlite_transaction

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var _ Connection = (*SQLiteConnection)(nil)

// SQLiteConnection base connection adapter
type SQLiteConnection struct {
	sqlite    *sql.DB
	db        *goqu.Database
	conf      Config
	sequences map[string]*Sequence
}

func NewSQLiteConnection(conf Config) *SQLiteConnection {
	return &SQLiteConnection{
		conf:      conf,
		sequences: make(map[string]*Sequence),
	}
}

func (con *SQLiteConnection) Open() error {
	var err error

	con.sqlite, err = sql.Open("sqlite3", con.conf.DBName)
	if err != nil {
		return errors.Wrapf(err, "opend database \"%s\"", con.conf.DBName)
	}

	con.db = goqu.New("sqlite3", con.sqlite)
	con.db.Logger(newGoquLogger(con.conf.Logger))

	if _, err = con.db.Exec(fmt.Sprintf("PRAGMA auto_vacuum=%s", con.conf.AutoVacuum)); err != nil {
		return errors.Wrapf(err, "pragma auto_vacuum=%s", con.conf.AutoVacuum)
	}
	if _, err = con.db.Exec(fmt.Sprintf("PRAGMA cache_size=%d", con.conf.CacheSizePages)); err != nil {
		return errors.Wrapf(err, "pragma cache_size=%d", con.conf.CacheSizePages)
	}
	if _, err = con.db.Exec(fmt.Sprintf("PRAGMA encoding=%s", "utf8")); err != nil {
		return errors.Wrapf(err, "pragma encoding=%s", "utf8")
	}
	if _, err = con.db.Exec(fmt.Sprintf("PRAGMA journal_mode=%s", con.conf.JournalMode)); err != nil {
		return errors.Wrapf(err, "pragma journal_mode=%s", con.conf.JournalMode)
	}
	if _, err = con.db.Exec(fmt.Sprintf("PRAGMA page_size=%d", con.conf.PageSizeBytes)); err != nil {
		return errors.Wrapf(err, "pragma page_size=%d", con.conf.PageSizeBytes)
	}
	if _, err = con.db.Exec(fmt.Sprintf("PRAGMA synchronous=%s", con.conf.Synchronous)); err != nil {
		return errors.Wrapf(err, "pragma synchronous=%s", con.conf.Synchronous)
	}
	if _, err = con.db.Exec(fmt.Sprintf("PRAGMA temp_store=%s", con.conf.TempStore)); err != nil {
		return errors.Wrapf(err, "pragma temp_store=%s", con.conf.TempStore)
	}

	return nil
}

// DB  return *goqu.Database
func (con *SQLiteConnection) DB() *sql.DB {
	return con.sqlite
}

// Begin starts a session for session dependent objects.
func (con *SQLiteConnection) Begin(ctx context.Context) (context.Context, Transaction, error) {
	return con.inject(ctx)
}

// Use - return database connection.
func (con *SQLiteConnection) Use(ctx context.Context) Database {
	tx, ok := con.extract(ctx)
	if !ok {
		return con.db
	}

	return tx.db
}

// RegSequence - register sequence
func (con *SQLiteConnection) RegSequence(table exp.IdentifierExpression) error {
	seq := NewSequence(con, table)

	err := seq.Refresh(context.Background())
	if err != nil {
		return errors.Wrapf(err, "refresh sequence table: %s", table.GetTable())
	}

	con.sequences[table.GetTable()] = seq

	return nil
}

// GetSequence - get registered sequence
func (con *SQLiteConnection) GetSequence(table exp.IdentifierExpression) (*Sequence, error) {
	seq, ok := con.sequences[table.GetTable()]
	if !ok {
		return nil, errors.Errorf("sequence table: %s not found", table.GetTable())
	}

	return seq, nil
}

// Wrap - automatically COMMIT or ROLLBACK.
func (con *SQLiteConnection) Wrap(ctx context.Context, tFunc func(ctx context.Context) error) error {
	ctxTx, tx, err := con.Begin(ctx)
	if err != nil {
		return err
	}

	err = tFunc(ctxTx)
	if err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			con.conf.Logger.Error("Rollback transaction Adapter.Transaction: %s", errRollback.Error())
		}

		return err
	}

	return tx.Commit()
}

// inject - inject transaction into context.
func (con *SQLiteConnection) inject(ctx context.Context) (context.Context, Transaction, error) {
	if _, ok := con.extract(ctx); ok {
		return ctx, newNoopTransaction(), nil
	}

	tx, err := con.db.Begin()
	if err != nil {
		return nil, nil, err
	}

	ctxTx := newTxAdapter(tx, con.conf.Logger)

	ctx = context.WithValue(ctx, transactionKey{}, ctxTx)

	return ctx, ctxTx, nil
}

// extract - extract transaction from context.
func (con *SQLiteConnection) extract(ctx context.Context) (tx *TxAdapter, ok bool) {
	tx, ok = ctx.Value(transactionKey{}).(*TxAdapter)

	return tx, ok
}
