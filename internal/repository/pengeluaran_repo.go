package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"minimal_api/internal/domain"
)

type pengeluaranRepository struct {
	db *sqlx.DB
}

func NewPengeluaranRepository(db *sqlx.DB) domain.PengeluaranRepository {
	return &pengeluaranRepository{db: db}
}

func (r *pengeluaranRepository) validateCustomerExists(ctx context.Context, customerID int) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, "SELECT COUNT(*) FROM customer WHERE customer_pk = ?", customerID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *pengeluaranRepository) GetPengeluaranByID(ctx context.Context, trxOutNo string) (*domain.PengeluaranBarangHeader, error) {
	var header domain.PengeluaranBarangHeader

	queryHeader := `
        SELECT trx_out_pk, trx_out_no, whs_idf, trx_out_date, trx_out_supp_idf, trx_out_notes
        FROM pengeluaran_barang_header
        WHERE trx_out_no = ?`

	err := r.db.GetContext(ctx, &header, queryHeader, trxOutNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	queryDetails := `
        SELECT trx_out_d_product_idf, trx_out_d_qty_dus, trx_out_d_qty_pcs
        FROM pengeluaran_barang_detail
        WHERE trx_out_idf = ?`
	var details []domain.PengeluaranBarangDetail
	if err := r.db.SelectContext(ctx, &details, queryDetails, header.TrxOutPK); err != nil {
		return nil, err
	}

	header.Details = details
	return &header, nil
}

func (r *pengeluaranRepository) IsTrxOutNoExists(trxOutNo string) (bool, error) {
	var count int
	err := r.db.Get(&count, "SELECT COUNT(*) FROM pengeluaran_barang_header WHERE trx_out_no = ?", trxOutNo)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *pengeluaranRepository) CreatePengeluaran(ctx context.Context, pengeluaran *domain.PengeluaranBarangHeader) (int, error) {
	// Lakukan validasi terhadap foreign key trx_out_supp_idf (customer)
	exists, err := r.validateCustomerExists(ctx, pengeluaran.TrxOutSuppIdf)
	if err != nil {
		return 0, err
	}
	if !exists {

		return 0, fmt.Errorf("customer with id %d does not exist", pengeluaran.TrxOutSuppIdf)
	}

	tx, err := r.db.Beginx() // Mulai transaksi
	if err != nil {
		return 0, err
	}

	// Query untuk insert header pengeluaran barang
	stmt, err := tx.Preparex(`
		INSERT INTO pengeluaran_barang_header 
		(trx_out_no, whs_idf, trx_out_date, trx_out_supp_idf, trx_out_notes) 
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	pengeluaran.TrxOutNo = fmt.Sprintf("TRXOUT-%d", time.Now().UnixNano())

	res, err := stmt.Exec(pengeluaran.TrxOutNo, pengeluaran.WhsIdf, pengeluaran.TrxOutDate, pengeluaran.TrxOutSuppIdf, pengeluaran.TrxOutNotes)
	if err != nil {
		// Cek jika terjadi error duplicate entry
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			tx.Rollback()
			return 0, errors.New("duplicate entry: trx_out_no already exists")
		}
		tx.Rollback()
		return 0, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	stmtDetail, err := tx.Preparex(`
		INSERT INTO pengeluaran_barang_detail 
		(trx_out_idf, trx_out_d_product_idf, trx_out_d_qty_dus, trx_out_d_qty_pcs)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmtDetail.Close()

	for _, detail := range pengeluaran.Details {
		if _, err := stmtDetail.Exec(lastID, detail.TrxOutDProductIdf, detail.TrxOutDQtyDus, detail.TrxOutDQtyPcs); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return int(lastID), nil
}

func (r *pengeluaranRepository) ValidateForeignKeys(ctx context.Context, whsIdf int, suppIdf int) error {
	var warehouseExists bool
	err := r.db.GetContext(ctx, &warehouseExists, "SELECT EXISTS(SELECT 1 FROM warehouse WHERE whs_pk = ?)", whsIdf)
	if err != nil {
		return err
	}
	if !warehouseExists {
		return errors.New("warehouse not found")
	}

	return nil
}

func (r *pengeluaranRepository) ValidateProductIDs(ctx context.Context, productIDs []int) error {
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
