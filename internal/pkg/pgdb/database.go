package pgdb

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txCtxKey struct{}

func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

type Database struct {
	pool *pgxpool.Pool
}

type TransactionCallback func(ctx context.Context) error

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

func (d *Database) Select(ctx context.Context, query sq.SelectBuilder) (pgx.Rows, error) {
	tx, withTransaction := transactionFromContext(ctx)

	querySql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	if withTransaction {
		return tx.Query(ctx, querySql, args...)
	}
	return d.pool.Query(ctx, querySql, args...)
}

func (d *Database) Update(ctx context.Context, query sq.UpdateBuilder) (pgx.Rows, error) {
	tx, withTransaction := transactionFromContext(ctx)

	querySql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	if withTransaction {
		return tx.Query(ctx, querySql, args...)
	}
	return d.pool.Query(ctx, querySql, args...)
}

func (d *Database) Insert(ctx context.Context, query sq.InsertBuilder) (pgx.Rows, error) {
	tx, withTransaction := transactionFromContext(ctx)

	querySql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	if withTransaction {
		return tx.Query(ctx, querySql, args...)
	}
	return d.pool.Query(ctx, querySql, args...)
}

func (d *Database) Delete(ctx context.Context, query sq.DeleteBuilder) (pgx.Rows, error) {
	tx, withTransaction := transactionFromContext(ctx)

	querySql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
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
