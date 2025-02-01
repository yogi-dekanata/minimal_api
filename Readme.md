# 📦 Minimal API - Warehouse Management

## 📖 Description
Minimal API is a warehouse management system for **Inbound and Outbound Goods**, developed using **Golang** and **React**. This application records goods received from suppliers and goods issued to customers with real-time stock reports.

---

## 🚀 Technologies Used
- **Backend**: Golang (Gin, SQLX, MySQL)
- **Frontend**: React (If available)
- **Database**: MySQL

---

## 📂 Directory Structure
```
minimal_api/
│-- internal/
│   ├── domain/             # Data models definition
│   ├── handler/            # API handlers
│   ├── repository/         # Database interactions
│   ├── router/             # API routing
│   ├── usecase/            # Business logic
│-- pkg/                    # Helper packages (middleware, validator, config)
│-- main.go                 # Application entry point
│-- go.mod                  # Module dependencies
│-- README.md               # This documentation
```

---

## 🛠️ Installation & Setup

### **1️⃣ Prerequisites**
Ensure **Go, MySQL, and Git** are installed.

### **2️⃣ Clone Repository**
```sh
git clone <repository-url>
cd minimal_api
```

### **3️⃣ Setup Database**
**Create the database and run the following SQL:**
```sql
CREATE DATABASE warehouse_db;
USE warehouse_db;
-- Run migration.sql if available
```

### **4️⃣ Configure Environment**
Create a `.env` file and add the following configurations:
```
DB_USER=root
DB_PASSWORD=password
DB_HOST=localhost
DB_PORT=3306
DB_NAME=warehouse_db
SERVER_PORT=8080
```

### **5️⃣ Install Dependencies**
```sh
go mod tidy
```

### **6️⃣ Run the Application**
```sh
go run main.go
```

The application will run at **http://localhost:8080**

---

## 📡 API Endpoints

### 📥 **Inbound Goods (Receiving)**
✅ **Create Inbound Goods**  
`POST /penerimaan`
```json
{
    "whs_idf": 1,
    "trx_in_date": "2024-02-10",
    "trx_in_supp_idf": 1,
    "trx_in_notes": "Receiving Goods A & B at Warehouse A",
    "details": [
        { "trx_in_d_product_idf": 1, "trx_in_d_qty_dus": 10, "trx_in_d_qty_pcs": 20 },
        { "trx_in_d_product_idf": 2, "trx_in_d_qty_dus": 5, "trx_in_d_qty_pcs": 10 }
    ]
}
```

✅ **Get Inbound Goods**  
`GET /penerimaan/{trx_in_no}`

---

### 📤 **Outbound Goods (Issuance)**
✅ **Create Outbound Goods**  
`POST /pengeluaran`
```json
{
    "whs_idf": 2,
    "trx_out_date": "2024-02-13",
    "trx_out_supp_idf": 3,
    "trx_out_notes": "Goods Issuance from Warehouse B",
    "details": [
        { "trx_out_d_product_idf": 3, "trx_out_d_qty_dus": 2, "trx_out_d_qty_pcs": 4 }
    ]
}
```

✅ **Get Outbound Goods**  
`GET /pengeluaran/{trx_out_no}`

---

### 📊 **Stock Report**
✅ **Get Stock Report**  
`GET /stok`

Sample Response:
```json
{
    "data": [
        {"gudang": "Warehouse A", "produk": "Goods A", "qty_dus": 10, "qty_pcs": 20},
        {"gudang": "Warehouse B", "produk": "Goods C", "qty_dus": 5, "qty_pcs": 10}
    ],
    "message": "Stock report retrieved successfully"
}
```

---

## 📞 Support & Contribution
- **Report Issues**: If you find any bugs or want to suggest features, create an **Issue** on GitHub.
- **Contribute**: Feel free to make a **Pull Request (PR)** if you want to contribute!

