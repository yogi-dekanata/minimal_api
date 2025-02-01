package handler

import "minimal_api/internal/usecase"

type Handlers struct {
	PenerimaanHandler  *PenerimaanHandler
	PengeluaranHandler *PengeluaranHandler
	StockHandler       *StokHandler
}

func NewHandlers(
	penerimaanUsecase usecase.PenerimaanUsecase,
	pengeluaranUsecase usecase.PengeluaranUsecase,
	stockUsecase usecase.StokUsecase,

) *Handlers {
	return &Handlers{
		PenerimaanHandler:  NewPenerimaanHandler(penerimaanUsecase),
		PengeluaranHandler: NewPengeluaranHandler(pengeluaranUsecase),
		StockHandler:       NewStokHandler(stockUsecase),
	}
}
