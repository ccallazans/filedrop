package repository

import (
	"context"

	"gorm.io/gorm"
)

func HasTransaction(ctx context.Context, tx *gorm.DB) *gorm.DB {
	
	tr := tx
	hasTransaction := ctx.Value("tx")
	if hasTransaction != nil {
		tr = hasTransaction.(*gorm.DB)
	}

	return tr
}
