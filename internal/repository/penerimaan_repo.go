package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"time"

	"github.com/jmoiron/sqlx"
	"minimal_api/internal/domain"
)

type penerimaanRepository struct {
	db *sqlx.DB
}

func generateTrxInNo() string {
	timestamp := time.Now().Format("20060102150405") // Format YYYYMMDDHHMMSS
	return "TRXIN-" + timestamp
}

func NewPenerimaanRepository(db *sqlx.DB) domain.PenerimaanRepository {
	return &penerimaanRepository{db: db}
}

func (r *penerimaanRepository) GetPenerimaanByID(ctx context.Context, trxInNo string) (*domain.PenerimaanBarangHeader, error) {
	penerimaan := domain.PenerimaanBarangHeader{}
	err := r.db.Get(&penerimaan, `
		SELECT trx_in_pk, trx_in_no, whs_idf, trx_in_date, trx_in_supp_idf, trx_in_notes
		FROM penerimaan_barang_header
		WHERE trx_in_no = ?`, trxInNo)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	var details []domain.PenerimaanBarangDetail
	err = r.db.Select(&details, `
		SELECT trx_in_dpk, trx_in_idf, trx_in_d_product_idf, trx_in_d_qty_dus, trx_in_d_qty_pcs
		FROM penerimaan_barang_detail
		WHERE trx_in_idf = ?`, penerimaan.TrxInPK)

	if err != nil {
		return nil, err
	}

	penerimaan.Details = details

	return &penerimaan, nil
}

func (r *penerimaanRepository) IsTrxInNoExists(trxInNo string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM penerimaan_barang_header WHERE trx_in_no = ?", trxInNo)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *penerimaanRepository) CreatePenerimaan(ctx context.Context, penerimaan *domain.PenerimaanBarangHeader) (int, error) {
	tx, err := r.db.Beginx() // Mulai transaksi
	if err != nil {
		return 0, err
	}

	// 🔹 Query untuk insert header penerimaan barang
	stmt, err := tx.Preparex("INSERT INTO penerimaan_barang_header (trx_in_no, whs_idf, trx_in_date, trx_in_supp_idf, trx_in_notes) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	penerimaan.TrxInNo = fmt.Sprintf("TRXIN-%d", time.Now().UnixNano())

	res, err := stmt.Exec(penerimaan.TrxInNo, penerimaan.WhsIdf, penerimaan.TrxInDate, penerimaan.TrxInSuppIdf, penerimaan.TrxInNotes)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			tx.Rollback()
			return 0, errors.New("duplicate entry: trx_in_no already exists")
		}
		tx.Rollback()
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmtDetail, err := tx.Preparex("INSERT INTO penerimaan_barang_detail (trx_in_idf, trx_in_d_product_idf, trx_in_d_qty_dus, trx_in_d_qty_pcs) VALUES (?, ?, ?, ?)")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmtDetail.Close()

	for _, detail := range penerimaan.Details {
		_, err := stmtDetail.Exec(lastID, detail.TrxInDProductIdf, detail.TrxInDQtyDus, detail.TrxInDQtyPcs)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func (r *penerimaanRepository) ValidateForeignKeys(ctx context.Context, whsIdf int, suppIdf int) error {
	var warehouseExists bool
	err := r.db.GetContext(ctx, &warehouseExists, "SELECT EXISTS(SELECT 1 FROM warehouse WHERE whs_pk = ?)", whsIdf)
	if err != nil {
		return err
	}
	if !warehouseExists {
		return errors.New("warehouse not found")
	}

	var supplierExists bool
	err = r.db.GetContext(ctx, &supplierExists, "SELECT EXISTS(SELECT 1 FROM supplier WHERE supplier_pk = ?)", suppIdf)
	if err != nil {
		return err
	}
	if !supplierExists {
		return errors.New("supplier not found")
	}

	return nil
}

func (r *penerimaanRepository) ValidateProductIDs(ctx context.Context, productIDs []int) error {
	if len(productIDs) == 0 {
		return errors.New("no products provided")
	}

	// Query untuk memastikan semua product IDs ada di database
	query := "SELECT COUNT(*) FROM product WHERE product_pk IN (?)"
	query, args, err := sqlx.In(query, productIDs)
	if err != nil {
		return err
	}
	query = r.db.Rebind(query)

	var count int
	err = r.db.GetContext(ctx, &count, query, args...)
	if err != nil {
		return err
	}

	// Jika jumlah produk yang ditemukan tidak sesuai dengan jumlah produk yang dikirim, maka ada yang tidak valid
	if count != len(productIDs) {
		return errors.New("one or more products not found")
	}

	return nil
}
