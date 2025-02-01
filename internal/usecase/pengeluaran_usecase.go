package usecase

import (
	"context"
	"errors"
	"fmt"
	"minimal_api/internal/domain"
	"minimal_api/internal/handler/dto"
	"time"
)

// PengeluaranUsecase ...
type PengeluaranUsecase interface {
	GetPengeluaranByID(ctx context.Context, trxOutNo string) (*domain.PengeluaranBarangHeader, error)
	CreatePengeluaran(ctx context.Context, request *dto.PengeluaranRequest) (int, string, error)
}

type pengeluaranUsecase struct {
	pengeluaranRepo domain.PengeluaranRepository
}

func NewPengeluaranUsecase(repo domain.PengeluaranRepository) PengeluaranUsecase {
	return &pengeluaranUsecase{pengeluaranRepo: repo}
}

func (u *pengeluaranUsecase) GetPengeluaranByID(ctx context.Context, trxOutNo string) (*domain.PengeluaranBarangHeader, error) {
	pengeluaran, err := u.pengeluaranRepo.GetPengeluaranByID(ctx, trxOutNo)
	if err != nil {
		return nil, err
	}
	if pengeluaran == nil {
		return nil, errors.New("pengeluaran not found")
	}
	return pengeluaran, nil
}

func (u *pengeluaranUsecase) CreatePengeluaran(ctx context.Context, request *dto.PengeluaranRequest) (int, string, error) {
	parsedDate, err := time.Parse(time.RFC3339, request.TrxOutDate)
	if err != nil {
		return 0, "", fmt.Errorf("invalid date format: %w", err)
	}

	formattedDate := parsedDate.Format("2006-01-02 15:04:05")

	pengeluaran := domain.PengeluaranBarangHeader{
		WhsIdf:        request.WhsIdf,
		TrxOutDate:    formattedDate, // Menggunakan string dengan format yang benar
		TrxOutSuppIdf: request.TrxOutSuppIdf,
		TrxOutNotes:   request.TrxOutNotes,
	}

	for _, detail := range request.Details {
		pengeluaran.Details = append(pengeluaran.Details, domain.PengeluaranBarangDetail{
			TrxOutDProductIdf: detail.TrxOutDProductIdf,
			TrxOutDQtyDus:     detail.TrxOutDQtyDus,
			TrxOutDQtyPcs:     detail.TrxOutDQtyPcs,
		})
	}

	pengeluaran.TrxOutNo = fmt.Sprintf("TRXOUT-%d", time.Now().UnixNano())
	if err := u.pengeluaranRepo.ValidateForeignKeys(ctx, request.WhsIdf, request.TrxOutSuppIdf); err != nil {
		return 0, "", err
	}

	var productIDs []int
	for _, detail := range request.Details {
		productIDs = append(productIDs, detail.TrxOutDProductIdf)
	}

	if err := u.pengeluaranRepo.ValidateProductIDs(ctx, productIDs); err != nil {
		return 0, "", err
	}

	// Simpan ke database via repository
	trxID, err := u.pengeluaranRepo.CreatePengeluaran(ctx, &pengeluaran)
	if err != nil {
		return 0, "", err
	}

	return trxID, pengeluaran.TrxOutNo, nil
}
