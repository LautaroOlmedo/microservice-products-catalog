package my_sql

import (
	"context"

	"gorm.io/gorm"
)

type txKey struct{}

type TxManager struct {
	db *gorm.DB
}

func NewTxManager(db *gorm.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {

	return m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func GetTx(ctx context.Context) (*gorm.DB, bool) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	return tx, ok
}

/*type txKey struct{}

type TxManager struct {
	db *sql.DB
}

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTransaction(
	ctx context.Context,
	fn func(ctx context.Context) error,
) error {

	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Inyectamos la tx en el contexto
	txCtx := context.WithValue(ctx, txKey{}, tx)

	if err := fn(txCtx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

func GetTx(ctx context.Context, db *sql.DB) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}
*/
