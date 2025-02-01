package dto

import (
	"fmt"
	"time"
)

type PengeluaranRequest struct {
	WhsIdf        int                        `json:"whs_idf" binding:"required"`
	TrxOutDate    string                     `json:"trx_out_date" binding:"required"`
	TrxOutSuppIdf int                        `json:"trx_out_supp_idf" binding:"required"`
	TrxOutNotes   string                     `json:"trx_out_notes"`
	Details       []PengeluaranDetailRequest `json:"details" binding:"required,dive"`
}

type PengeluaranDetailRequest struct {
	TrxOutDProductIdf int `json:"trx_out_d_product_idf" binding:"required"`
	TrxOutDQtyDus     int `json:"trx_out_d_qty_dus" binding:"required"`
	TrxOutDQtyPcs     int `json:"trx_out_d_qty_pcs" binding:"required"`
}

type PengeluaranResponse struct {
	TrxOutNo      string                      `json:"trx_out_no"`
	WhsIdf        int                         `json:"whs_idf"`
	TrxOutDate    string                      `json:"trx_out_date"`
	TrxOutSuppIdf int                         `json:"trx_out_supp_idf"`
	TrxOutNotes   string                      `json:"trx_out_notes"`
	Details       []PengeluaranDetailResponse `json:"details"`
}

type PengeluaranDetailResponse struct {
	TrxOutDProductIdf int `json:"trx_out_d_product_idf"`
	TrxOutDQtyDus     int `json:"trx_out_d_qty_dus"`
	TrxOutDQtyPcs     int `json:"trx_out_d_qty_pcs"`
}

func (r *PengeluaranRequest) ConvertDate() error {
	parsedTime, err := time.Parse("2006-01-02", r.TrxOutDate)
	if err != nil {
		return fmt.Errorf("invalid date format: expected YYYY-MM-DD")
	}
	r.TrxOutDate = parsedTime.Format("2006-01-02T15:04:05Z") // Format yang sesuai
	return nil
}
