package usecase

import (
	"context"
	"minimal_api/internal/domain"
	"minimal_api/internal/handler/dto"
)

// StokUsecase ...
type StokUsecase interface {
	// GetStok mengambil laporan stok sebagai slice dari dto.StokResponse.
	GetStok(ctx context.Context) ([]dto.StokResponse, error)
}

// stokUsecase ...
type stokUsecase struct {
	stokRepo domain.StokRepository
}

// NewStokUsecase ...
func NewStokUsecase(repo domain.StokRepository) StokUsecase {
	return &stokUsecase{stokRepo: repo}
}

// GetStok ...
func (s *stokUsecase) GetStok(ctx context.Context) ([]dto.StokResponse, error) {
	// Ambil data stok dari repository (yang mengakses view stok_barang).
	stokReports, err := s.stokRepo.GetStok(ctx)
	if err != nil {
		return nil, err
	}

	var responses []dto.StokResponse
	for _, report := range stokReports {
		responses = append(responses, dto.StokResponse{
			Gudang: report.Gudang,
			Produk: report.Produk,
			QtyDus: report.QtyDus,
			QtyPcs: report.QtyPcs,
		})
	}

	return responses, nil
}
