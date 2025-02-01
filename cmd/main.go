package main

import (
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	// Import internal packages
	"minimal_api/internal/handler"
	"minimal_api/internal/repository"
	"minimal_api/internal/router"
	"minimal_api/internal/usecase"

	// Import pkg
	"minimal_api/pkg/config"
)

func main() {
	// 1. Load Config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Setup DB
	dsn := cfg.Database.User + ":" + cfg.Database.Password + "@tcp(" +
		cfg.Database.Host + ":" + strconv.Itoa(cfg.Database.Port) + ")/" +
		cfg.Database.DBName + "?parseTime=true"

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect DB: %v", err)
	}
	defer db.Close()

	// 3. Inisialisasi Repository
	penerimaanRepo := repository.NewPenerimaanRepository(db)
	pengeluaranRepo := repository.NewPengeluaranRepository(db)
	stockRepo := repository.NewStockRepository(db)

	// 4. Inisialisasi Usecase
	penerimaanUsecase := usecase.NewPenerimaanUsecase(penerimaanRepo)
	pengeluaranUsecase := usecase.NewPengeluaranUsecase(pengeluaranRepo)
	stockUsecase := usecase.NewStokUsecase(stockRepo)

	// 5. Inisialisasi Handler (dipisahkan di init.go)
	handlers := handler.NewHandlers(penerimaanUsecase, pengeluaranUsecase,
		stockUsecase,
	)

	// 6. Setup Router (dipisah ke router.go)
	r := router.SetupRouter(handlers)

	// 7. Jalankan Server
	log.Printf("Server running on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
