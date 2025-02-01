package usecase

import (
	"context"
	"errors"
	"fmt"
	"minimal_api/internal/domain"
	"minimal_api/internal/handler/dto"
	"time"
)

type PenerimaanUsecase interface {
	GetPenerimaanByID(ctx context.Context, trxInNo string) (*domain.PenerimaanBarangHeader, error)
	CreatePenerimaan(ctx context.Context, request *dto.PenerimaanRequest) (int, string, error)
}

type penerimaanUsecase struct {
	penerimaanRepo domain.PenerimaanRepository
}

func NewPenerimaanUsecase(repo domain.PenerimaanRepository) PenerimaanUsecase {
	return &penerimaanUsecase{penerimaanRepo: repo}
}

func (u *penerimaanUsecase) GetPenerimaanByID(ctx context.Context, trxInNo string) (*domain.PenerimaanBarangHeader, error) {
	penerimaan, err := u.penerimaanRepo.GetPenerimaanByID(ctx, trxInNo)
	if err != nil {
		return nil, err
	}

	if penerimaan == nil {
		return nil, errors.New("penerimaan not found")
	}

	return penerimaan, nil
}

func (u *penerimaanUsecase) CreatePenerimaan(ctx context.Context, request *dto.PenerimaanRequest) (int, string, error) {
	penerimaan := domain.PenerimaanBarangHeader{
		WhsIdf:       request.WhsIdf,
		TrxInDate:    request.TrxInDate,
		TrxInSuppIdf: request.TrxInSuppIdf,
		TrxInNotes:   request.TrxInNotes,
	}

	for _, detail := range request.Details {
		penerimaan.Details = append(penerimaan.Details, domain.PenerimaanBarangDetail{
			TrxInDProductIdf: detail.TrxInDProductIdf,
			TrxInDQtyDus:     detail.TrxInDQtyDus,
			TrxInDQtyPcs:     detail.TrxInDQtyPcs,
		})
	}
	if err := u.penerimaanRepo.ValidateForeignKeys(ctx, request.WhsIdf, request.TrxInSuppIdf); err != nil {
		return 0, "", err
	}

	var productIDs []int
	for _, detail := range request.Details {
		productIDs = append(productIDs, detail.TrxInDProductIdf)
	}

	if err := u.penerimaanRepo.ValidateProductIDs(ctx, productIDs); err != nil {
		return 0, "", err
	}

	penerimaan.TrxInNo = fmt.Sprintf("TRXIN-%d", time.Now().Unix())

	trxID, err := u.penerimaanRepo.CreatePenerimaan(ctx, &penerimaan)
	if err != nil {
		return 0, "", err
	}

	return trxID, penerimaan.TrxInNo, nil
}
