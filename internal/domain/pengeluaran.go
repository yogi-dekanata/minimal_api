package domain

import "context"

type PengeluaranBarangHeader struct {
	TrxOutPK      int                       `json:"trx_out_pk" db:"trx_out_pk"`
	TrxOutNo      string                    `json:"trx_out_no" db:"trx_out_no"`
	WhsIdf        int                       `json:"whs_idf" binding:"required" db:"whs_idf"`
	TrxOutDate    string                    `json:"trx_out_date" binding:"required" db:"trx_out_date"`
	TrxOutSuppIdf int                       `json:"trx_out_supp_idf" binding:"required" db:"trx_out_supp_idf"`
	TrxOutNotes   string                    `json:"trx_out_notes" db:"trx_out_notes"`
	Details       []PengeluaranBarangDetail `json:"details" db:"details"`
}

type PengeluaranBarangDetail struct {
	TrxOutDPK         int `json:"trx_out_dpk" db:"trx_out_dpk"`
	TrxOutIDF         int `json:"trx_out_idf" db:"trx_out_idf"`
	TrxOutDProductIdf int `json:"trx_out_d_product_idf" db:"trx_out_d_product_idf"`
	TrxOutDQtyDus     int `json:"trx_out_d_qty_dus" db:"trx_out_d_qty_dus"`
	TrxOutDQtyPcs     int `json:"trx_out_d_qty_pcs" db:"trx_out_d_qty_pcs"`
}

type PengeluaranRepository interface {
	GetPengeluaranByID(ctx context.Context, trxOutNo string) (*PengeluaranBarangHeader, error)
	ValidateProductIDs(ctx context.Context, productIDs []int) error
	CreatePengeluaran(ctx context.Context, pengeluaran *PengeluaranBarangHeader) (int, error)
	ValidateForeignKeys(ctx context.Context, whsIdf int, suppIdf int) error
}
