package router

import (
	"github.com/gin-gonic/gin"
	"minimal_api/internal/handler"
	"minimal_api/pkg/middleware"
)

// SetupRouter ...
func SetupRouter(handlers *handler.Handlers) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	// Penerimaan routes
	r.GET("/penerimaan/:id", handlers.PenerimaanHandler.GetPenerimaan)
	r.POST("/penerimaan", handlers.PenerimaanHandler.CreatePenerimaan)

	// Pengeluaran routes
	r.GET("/pengeluaran/:id", handlers.PengeluaranHandler.GetPengeluaranByID)
	r.POST("/pengeluaran", handlers.PengeluaranHandler.CreatePengeluaran)

	r.GET("/stocks", handlers.StockHandler.GetStok)

	return r
}
