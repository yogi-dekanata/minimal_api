package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"minimal_api/internal/domain"
)

// stockRepository ...
type stockRepository struct {
	db *sqlx.DB
}

// NewStockRepository ...
func NewStockRepository(db *sqlx.DB) domain.StokRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) GetStok(ctx context.Context) ([]domain.StokReport, error) {
	var reports []domain.StokReport

	query := `
		SELECT gudang, produk, qty_dus, qty_pcs
		FROM stock_barang
	`
	if err := r.db.SelectContext(ctx, &reports, query); err != nil {
		return nil, err
	}

	return reports, nil
}
