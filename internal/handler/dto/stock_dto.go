package dto

type StokResponse struct {
	Gudang string `json:"gudang"`
	Produk string `json:"produk"`
	QtyDus int    `json:"qty_dus"`
	QtyPcs int    `json:"qty_pcs"`
}
