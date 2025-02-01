-- Buat database jika belum ada
CREATE DATABASE IF NOT EXISTS warehouse_db;
USE warehouse_db;

-- Table Master Supplier
CREATE TABLE IF NOT EXISTS supplier
(
    supplier_pk   INT PRIMARY KEY AUTO_INCREMENT,
    supplier_name VARCHAR(255) NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE INDEX idx_supplier_name ON supplier (supplier_name);

-- Table Master Customer
CREATE TABLE IF NOT EXISTS customer
(
    customer_pk   INT PRIMARY KEY AUTO_INCREMENT,
    customer_name VARCHAR(255) NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE INDEX idx_customer_name ON customer (customer_name);

-- Table Master Product
CREATE TABLE IF NOT EXISTS product
(
    product_pk   INT PRIMARY KEY AUTO_INCREMENT,
    product_name VARCHAR(255) NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE INDEX idx_product_name ON product (product_name);

-- Table Master Warehouse
CREATE TABLE IF NOT EXISTS warehouse
(
    whs_pk   INT PRIMARY KEY AUTO_INCREMENT,
    whs_name VARCHAR(255) NOT NULL
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE INDEX idx_warehouse_name ON warehouse (whs_name);

-- Table Transaksi Penerimaan Barang (Header)
CREATE TABLE IF NOT EXISTS penerimaan_barang_header
(
    trx_in_pk       INT PRIMARY KEY AUTO_INCREMENT,
    trx_in_no       VARCHAR(50) UNIQUE NOT NULL,
    whs_idf         INT                NOT NULL,
    trx_in_date     DATE               NOT NULL,
    trx_in_supp_idf INT                NOT NULL,
    trx_in_notes    TEXT,
    FOREIGN KEY (whs_idf) REFERENCES warehouse (whs_pk) ON DELETE CASCADE,
    FOREIGN KEY (trx_in_supp_idf) REFERENCES supplier (supplier_pk) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE INDEX idx_penerimaan_whs ON penerimaan_barang_header (whs_idf);
CREATE INDEX idx_penerimaan_supp ON penerimaan_barang_header (trx_in_supp_idf);
CREATE INDEX idx_penerimaan_date ON penerimaan_barang_header (trx_in_date);
CREATE UNIQUE INDEX idx_penerimaan_no ON penerimaan_barang_header (trx_in_no);

-- Table Transaksi Penerimaan Barang (Detail)
CREATE TABLE IF NOT EXISTS penerimaan_barang_detail
(
    trx_in_dpk           INT PRIMARY KEY AUTO_INCREMENT,
    trx_in_idf           INT NOT NULL,
    trx_in_d_product_idf INT NOT NULL,
    trx_in_d_qty_dus     INT NOT NULL DEFAULT 0,
    trx_in_d_qty_pcs     INT NOT NULL DEFAULT 0,
    FOREIGN KEY (trx_in_idf) REFERENCES penerimaan_barang_header (trx_in_pk) ON DELETE CASCADE,
    FOREIGN KEY (trx_in_d_product_idf) REFERENCES product (product_pk) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE INDEX idx_penerimaan_detail_product ON penerimaan_barang_detail (trx_in_d_product_idf);
CREATE INDEX idx_penerimaan_detail_header ON penerimaan_barang_detail (trx_in_idf);

-- Table Transaksi Pengeluaran Barang (Header)
CREATE TABLE IF NOT EXISTS pengeluaran_barang_header
(
    trx_out_pk       INT PRIMARY KEY AUTO_INCREMENT,
    trx_out_no       VARCHAR(50) UNIQUE NOT NULL,
    whs_idf          INT                NOT NULL,
    trx_out_date     DATE               NOT NULL,
    trx_out_supp_idf INT                NOT NULL,
    trx_out_notes    TEXT,
    FOREIGN KEY (whs_idf) REFERENCES warehouse (whs_pk) ON DELETE CASCADE,
    FOREIGN KEY (trx_out_supp_idf) REFERENCES customer (customer_pk) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE INDEX idx_pengeluaran_whs ON pengeluaran_barang_header (whs_idf);
CREATE INDEX idx_pengeluaran_supp ON pengeluaran_barang_header (trx_out_supp_idf);
CREATE INDEX idx_pengeluaran_date ON pengeluaran_barang_header (trx_out_date);
CREATE UNIQUE INDEX idx_pengeluaran_no ON pengeluaran_barang_header (trx_out_no);

-- Table Transaksi Pengeluaran Barang (Detail)
CREATE TABLE IF NOT EXISTS pengeluaran_barang_detail
(
    trx_out_dpk           INT PRIMARY KEY AUTO_INCREMENT,
    trx_out_idf           INT NOT NULL,
    trx_out_d_product_idf INT NOT NULL,
    trx_out_d_qty_dus     INT NOT NULL DEFAULT 0,
    trx_out_d_qty_pcs     INT NOT NULL DEFAULT 0,
    FOREIGN KEY (trx_out_idf) REFERENCES pengeluaran_barang_header (trx_out_pk) ON DELETE CASCADE,
    FOREIGN KEY (trx_out_d_product_idf) REFERENCES product (product_pk) ON DELETE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;

CREATE INDEX idx_pengeluaran_detail_product ON pengeluaran_barang_detail (trx_out_d_product_idf);
CREATE INDEX idx_pengeluaran_detail_header ON pengeluaran_barang_detail (trx_out_idf);

-- Optimasi Query pada View StokBarang
CREATE INDEX idx_stok_barang_gudang ON warehouse (whs_pk);
CREATE INDEX idx_stok_barang_product ON product (product_pk);
CREATE INDEX idx_stok_penerimaan ON penerimaan_barang_detail (trx_in_d_product_idf);
CREATE INDEX idx_stok_pengeluaran ON pengeluaran_barang_detail (trx_out_d_product_idf);

-- View Stok Barang (Menggunakan Index yang Ditambahkan)

