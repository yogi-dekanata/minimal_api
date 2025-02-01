package domain

import "context"

// StokReport mewakili data laporan stok dari view stok_barang.
type StokReport struct {
	Gudang string `db:"gudang"`
	Produk string `db:"produk"`
	QtyDus int    `db:"qty_dus"`
	QtyPcs int    `db:"qty_pcs"`
}

// StokRepository adalah kontrak untuk mengambil data stok.
type StokRepository interface {
	// GetStok mengambil data stok dari view stok_barang.
	GetStok(ctx context.Context) ([]StokReport, error)
}
