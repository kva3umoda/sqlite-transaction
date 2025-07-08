package sqlite_transaction

import (
	"github.com/doug-martin/goqu/v9"
)

type transactionKey struct{}

// Transaction - manager transaction.
type Transaction interface {
	Rollback() error
	Commit() error
	End(err *error)
}

var _ Transaction = (*noopTransaction)(nil)

type noopTransaction struct {
}

func newNoopTransaction() *noopTransaction {
	return &noopTransaction{}
}

func (nt *noopTransaction) Rollback() error {
	return nil
}

func (nt *noopTransaction) Commit() error {
	return nil
}

func (nt *noopTransaction) End(err *error) {
}

var _ Transaction = (*TxAdapter)(nil)

type TxAdapter struct {
	logger Logger
	db     *goqu.TxDatabase
}

func newTxAdapter(tx *goqu.TxDatabase, logger Logger) *TxAdapter {
	return &TxAdapter{
		db:     tx,
		logger: logger,
	}
}

func (t *TxAdapter) Rollback() error {
	return t.db.Rollback()
}

func (t *TxAdapter) Commit() error {
	return t.db.Commit()
}

// End commit or rollback transaction.
func (t *TxAdapter) End(err *error) {
	if err != nil && *err != nil {
		errRollback := t.Rollback()
		if errRollback != nil {
			t.logger.Error("Rollback transaction Failed: %s", errRollback.Error())
		}

		return
	}

	errCommit := t.Commit()
	if errCommit != nil && err != nil {
		*err = errCommit
	}
}
