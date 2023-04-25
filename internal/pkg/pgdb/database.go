package pgdb

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type txCtxKey struct{}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

type Database struct {
	pool *pgxpool.Pool
}

type TransactionCallback func(ctx context.Context) error

func (d *Database) ReadonlyTx(ctx context.Context, callback TransactionCallback) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctx = context.WithValue(ctx, txCtxKey{}, tx)

	if err := callback(ctx); err != nil {
		return err
	}

	return nil
}

func (d *Database) Tx(ctx context.Context, callback TransactionCallback) error {
	tx, err := d.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	ctx = context.WithValue(ctx, txCtxKey{}, tx)

	if err := callback(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (d *Database) QuerySq(ctx context.Context, query sq.Sqlizer) (pgx.Rows, error) {
	tx, withTransaction := transactionFromContext(ctx)

	querySql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if withTransaction {
		return tx.Query(ctx, querySql, args...)
	}
	return d.pool.Query(ctx, querySql, args...)
}

func transactionFromContext(ctx context.Context) (pgx.Tx, bool) {
	if tx := ctx.Value(txCtxKey{}); tx != nil {
		return tx.(pgx.Tx), true
	}
	return nil, false
}
