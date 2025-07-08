package sqlite_transaction

import (
	"context"
	"sync/atomic"

	"github.com/doug-martin/goqu/v9"
	"github.com/doug-martin/goqu/v9/exp"
)

type Sequence struct {
	table exp.IdentifierExpression
	db    Connection
	value int64
}

// Sequence generate in values
func NewSequence(db Connection, table exp.IdentifierExpression) *Sequence {
	seq := &Sequence{
		value: 0,
		db:    db,
		table: table,
	}

	return seq
}

func (seq *Sequence) Refresh(ctx context.Context) error {
	scanner, err := seq.db.Use(ctx).
		Select(goqu.Func("ifnull", goqu.MAX("id"), 0)).
		From(seq.table).
		Executor().
		Scanner()
	if err != nil {
		return err
	}

	defer scanner.Close()

	if scanner.Next() {
		if err = scanner.ScanVal(&seq.value); err != nil {
			return err
		}
	}
	return nil
}

func (seq *Sequence) NextValue() int {
	return int(atomic.AddInt64(&seq.value, 1))
}

func (seq *Sequence) CurrValue() int {
	return int(atomic.LoadInt64(&seq.value))
}
