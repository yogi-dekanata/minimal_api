package domain

import "context"

type PenerimaanBarangHeader struct {
	TrxInPK      int                      `json:"trx_in_pk" db:"trx_in_pk"`
	TrxInNo      string                   `json:"trx_in_no" db:"trx_in_no"`
	WhsIdf       int                      `json:"whs_idf" binding:"required" db:"whs_idf"`
	TrxInDate    string                   `json:"trx_in_date" binding:"required" db:"trx_in_date"`
	TrxInSuppIdf int                      `json:"trx_in_supp_idf" binding:"required" db:"trx_in_supp_idf"`
	TrxInNotes   string                   `json:"trx_in_notes" db:"trx_in_notes"`
	Details      []PenerimaanBarangDetail `json:"details" db:"details"`
}

type PenerimaanBarangDetail struct {
	TrxInDPK         int `json:"trx_in_dpk" db:"trx_in_dpk"`
	TrxInIDF         int `json:"trx_in_idf" db:"trx_in_idf"`
	TrxInDProductIdf int `json:"trx_in_d_product_idf" db:"trx_in_d_product_idf"`
	TrxInDQtyDus     int `json:"trx_in_d_qty_dus" db:"trx_in_d_qty_dus"`
	TrxInDQtyPcs     int `json:"trx_in_d_qty_pcs" db:"trx_in_d_qty_pcs"`
}

type PenerimaanRepository interface {
	GetPenerimaanByID(ctx context.Context, trxInNo string) (*PenerimaanBarangHeader, error)
	CreatePenerimaan(ctx context.Context, penerimaan *PenerimaanBarangHeader) (int, error)
	ValidateForeignKeys(ctx context.Context, whsIdf int, suppIdf int) error
	ValidateProductIDs(ctx context.Context, productIDs []int) error
}
