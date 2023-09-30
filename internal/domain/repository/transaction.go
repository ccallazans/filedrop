package repository

import (
	"context"
	"database/sql"
)

func HasTransaction(ctx context.Context, db *sql.DB) *sql.Tx {
	tx, ok := ctx.Value("tx").(*sql.Tx)
	if ok {
		return tx
	}

	return nil
}
