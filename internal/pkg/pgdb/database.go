package pgdb

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type txCtxKey struct{}

// NewDatabase returns new *Database.
func NewDatabase(pool *pgxpool.Pool) *Database {
	return &Database{pool: pool}
}

// Database rules db connection.
type Database struct {
	pool *pgxpool.Pool
}

// TransactionCallback represents function that will be executed withing single db transaction
type TransactionCallback func(ctx context.Context) error

// ReadonlyTx makes transaction reading without commit.
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

// Tx executes a callback within a single transaction.
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

// QuerySq executes query with squirrel.
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
