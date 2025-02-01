package handler

import (
	"github.com/gin-gonic/gin"
	"log"
	"minimal_api/internal/handler/dto"
	"minimal_api/internal/usecase"
	"net/http"
)

type PengeluaranHandler struct {
	pengeluaranUsecase usecase.PengeluaranUsecase
}

func NewPengeluaranHandler(u usecase.PengeluaranUsecase) *PengeluaranHandler {
	return &PengeluaranHandler{pengeluaranUsecase: u}
}

func (h *PengeluaranHandler) GetPengeluaranByID(c *gin.Context) {
	// Ambil parameter "id" dari URL (misal: /pengeluaran/:id)
	trxOutNo := c.Param("id")
	if trxOutNo == "" {
		log.Println("[ERROR] Invalid ID parameter: empty string")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID parameter"})
		return
	}

	// Panggil usecase dengan parameter string
	pengeluaran, err := h.pengeluaranUsecase.GetPengeluaranByID(c.Request.Context(), trxOutNo)
	if err != nil {
		if err.Error() == "pengeluaran not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pengeluaran not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get pengeluaran", "details": err.Error()})
		return
	}

	var details []dto.PengeluaranDetailResponse
	for _, d := range pengeluaran.Details {
		details = append(details, dto.PengeluaranDetailResponse{
			TrxOutDProductIdf: d.TrxOutDProductIdf,
			TrxOutDQtyDus:     d.TrxOutDQtyDus,
			TrxOutDQtyPcs:     d.TrxOutDQtyPcs,
		})
	}

	response := dto.PengeluaranResponse{
		TrxOutNo:      pengeluaran.TrxOutNo,
		WhsIdf:        pengeluaran.WhsIdf,
		TrxOutDate:    pengeluaran.TrxOutDate,
		TrxOutSuppIdf: pengeluaran.TrxOutSuppIdf,
		TrxOutNotes:   pengeluaran.TrxOutNotes,
		Details:       details,
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Data pengeluaran barang berhasil diambil",
		"pengeluaran": response,
	})
}

func (h *PengeluaranHandler) CreatePengeluaran(c *gin.Context) {
	var request dto.PengeluaranRequest

	// Validasi JSON request
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Println("[ERROR] Invalid request payload:", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request payload",
		})
		return
	}

	if err := request.ConvertDate(); err != nil {
		log.Println("[ERROR] Invalid date format:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format", "details": err.Error()})
		return
	}

	_, trxOutNo, err := h.pengeluaranUsecase.CreatePengeluaran(c.Request.Context(), &request)
	if err != nil {
		if err.Error() == "warehouse not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse", "details": "Warehouse ID not found"})
			return
		}
		if err.Error() == "customer not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer", "details": "Customer ID not found"})
			return
		}
		if err.Error() == "one or more products not found" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product", "details": "One or more product IDs not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create pengeluaran",
			"details": err.Error(),
		})
		return
	}

	// Mapping response
	response := dto.PengeluaranResponse{
		TrxOutNo:      trxOutNo,
		WhsIdf:        request.WhsIdf,
		TrxOutDate:    request.TrxOutDate,
		TrxOutSuppIdf: request.TrxOutSuppIdf,
		TrxOutNotes:   request.TrxOutNotes,
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Pengeluaran barang berhasil ditambahkan",
		"TrxOutNo": response.TrxOutNo,
	})
}
