package dto

type PenerimaanRequest struct {
	WhsIdf       int                       `json:"whs_idf" binding:"required"`
	TrxInDate    string                    `json:"trx_in_date" binding:"required"`
	TrxInSuppIdf int                       `json:"trx_in_supp_idf" binding:"required"`
	TrxInNotes   string                    `json:"trx_in_notes"`
	Details      []PenerimaanDetailRequest `json:"details" binding:"required,dive"`
}

type PenerimaanDetailRequest struct {
	TrxInDProductIdf int `json:"trx_in_d_product_idf" binding:"required"`
	TrxInDQtyDus     int `json:"trx_in_d_qty_dus" binding:"required"`
	TrxInDQtyPcs     int `json:"trx_in_d_qty_pcs" binding:"required"`
}

type PenerimaanResponse struct {
	TrxInNo      string                     `json:"trx_in_no"`
	WhsIdf       int                        `json:"whs_idf"`
	TrxInDate    string                     `json:"trx_in_date"`
	TrxInSuppIdf int                        `json:"trx_in_supp_idf"`
	TrxInNotes   string                     `json:"trx_in_notes"`
	Details      []PenerimaanDetailResponse `json:"details"`
}

type PenerimaanDetailResponse struct {
	TrxInDProductIdf int `json:"trx_in_d_product_idf"`
	TrxInDQtyDus     int `json:"trx_in_d_qty_dus"`
	TrxInDQtyPcs     int `json:"trx_in_d_qty_pcs"`
}
