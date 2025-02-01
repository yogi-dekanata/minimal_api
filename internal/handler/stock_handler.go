package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"minimal_api/internal/usecase"
)

// StokHandler ...
type StokHandler struct {
	stokUsecase usecase.StokUsecase
}

// NewStokHandler ...
func NewStokHandler(u usecase.StokUsecase) *StokHandler {
	return &StokHandler{stokUsecase: u}
}

// GetStok ...
func (h *StokHandler) GetStok(c *gin.Context) {
	// Panggil usecase untuk mengambil laporan stok
	stokReports, err := h.stokUsecase.GetStok(c.Request.Context())
	if err != nil {
		log.Println("[ERROR] Failed to get stock report:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get stock report",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock report retrieved successfully",
		"data":    stokReports,
	})
}
