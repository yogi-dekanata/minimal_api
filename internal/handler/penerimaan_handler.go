package handler

import (
	"log"
	"minimal_api/internal/handler/dto"
	"net/http"

	"github.com/gin-gonic/gin"
	"minimal_api/internal/usecase"
)

type PenerimaanHandler struct {
	penerimaanUsecase usecase.PenerimaanUsecase
}

func NewPenerimaanHandler(u usecase.PenerimaanUsecase) *PenerimaanHandler {
	return &PenerimaanHandler{penerimaanUsecase: u}
}

func (h *PenerimaanHandler) CreatePenerimaan(c *gin.Context) {
	var request dto.PenerimaanRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("[ERROR] Invalid request payload:", err)
		c.JSON(http.StatusBadRequest,
			gin.H{"error": "Invalid request payload"})
		return
	}

	_, trxInNo, err := h.penerimaanUsecase.CreatePenerimaan(c.Request.Context(), &request)
	if err != nil {
		if err.Error() == "warehouse not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse", "details": "Warehouse ID not found"})
			return
		}
		if err.Error() == "supplier not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid supplier", "details": "Supplier ID not found"})
			return
		}
		if err.Error() == "one or more products not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product", "details": "One or more product IDs not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create penerimaan",
			"details": err.Error(),
		})
		return
	}

	response := dto.PenerimaanResponse{
		TrxInNo:      trxInNo,
		WhsIdf:       request.WhsIdf,
		TrxInDate:    request.TrxInDate,
		TrxInSuppIdf: request.TrxInSuppIdf,
		TrxInNotes:   request.TrxInNotes,
	}

	for _, detail := range request.Details {
		response.Details = append(response.Details, dto.PenerimaanDetailResponse{
			TrxInDProductIdf: detail.TrxInDProductIdf,
			TrxInDQtyDus:     detail.TrxInDQtyDus,
			TrxInDQtyPcs:     detail.TrxInDQtyPcs,
		})
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Penerimaan barang berhasil ditambahkan",
		"TrxInNo": response.TrxInNo,
	})
}

func (h *PenerimaanHandler) GetPenerimaan(c *gin.Context) {
	trxInNo := c.Param("id")
	if trxInNo == "" {
		log.Println("[ERROR] Invalid ID parameter: empty string")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Panggil usecase dengan parameter string
	penerimaan, err := h.penerimaanUsecase.GetPenerimaanByID(c.Request.Context(), trxInNo)
	if err != nil {
		if err.Error() == "penerimaan not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Penerimaan not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get penerimaan", "details": err.Error()})
		return
	}

	// Mapping data domain ke DTO response
	var details []dto.PenerimaanDetailResponse
	for _, d := range penerimaan.Details {
		details = append(details, dto.PenerimaanDetailResponse{
			TrxInDProductIdf: d.TrxInDProductIdf,
			TrxInDQtyDus:     d.TrxInDQtyDus,
			TrxInDQtyPcs:     d.TrxInDQtyPcs,
		})
	}

	response := dto.PenerimaanResponse{
		TrxInNo:      penerimaan.TrxInNo,
		WhsIdf:       penerimaan.WhsIdf,
		TrxInDate:    penerimaan.TrxInDate,
		TrxInSuppIdf: penerimaan.TrxInSuppIdf,
		TrxInNotes:   penerimaan.TrxInNotes,
		Details:      details,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Data penerimaan barang berhasil diambil",
		"penerimaan": response,
	})
}
