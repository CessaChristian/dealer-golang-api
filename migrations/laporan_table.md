## ERD Title: Dealer Otomotif Management System

## ERD Database

![image](/migrations/ERD.png)

## Entities and Their Attributes

A. Entity: Users

Attributes:
 - user_id (PK, bigint, AI)

 - name (varchar(150))
 - email (varchar(150), UNIQUE)
 - password (text)
 - role (varchar(50), default = "customer")
 - created_at (timestamp)
 - updated_at (timestamp)

B. Entity: Vehicle_Types

Attributes:

 - type_id (PK, bigint, AI)

 - type_name (varchar(100), UNIQUE)

C. Entity: Brands

Attributes:

 - brand_id (PK, bigint, AI)

 - brand_name (varchar(100), UNIQUE)

D. Entity: Vehicles

Attributes:

 - vehicle_id (PK, bigint, AI)

 - type_id (FK → vehicle_types.type_id)
 - brand_id (FK → brands.brand_id)
 - name (varchar(150))
 - fuel_type (varchar(50))
 - transmission (varchar(50))
 - price (numeric(12,2))
 - stock (int, default 0)
 - created_at (timestamp)
 - updated_at (timestamp)

E. Entity: Favorites

Attributes:

 - user_id (PK, FK → users.user_id)

 - vehicle_id (PK, FK → vehicles.vehicle_id)
 - created_at (timestamp)

F. Entity: Transactions

Attributes:

 - transaction_id (PK, bigint, AI)

 - user_id (FK → users.user_id)
 - total_amount (numeric(12,2))
 - created_at (timestamp)
 - updated_at (timestamp)

G. Entity: Transaction_Items

Attributes:

 - detail_id (PK, bigint, AI)

 - transaction_id (FK → transactions.transaction_id)
 - vehicle_id (FK → vehicles.vehicle_id)
 - quantity (int)
 - price (numeric(12,2))

H. Entity: Payment_Detail

Attributes:

 - payment_id (PK, bigint, AI)

 - transaction_id (FK → transactions.transaction_id, UNIQUE)

 - payment_method (varchar(50))

 - status (varchar(50))

 - paid_at (timestamp, nullable)

 - note (text, nullable)

## Relationships
A. Users → Transactions

- Type: One-to-Many
- Description:
Satu user dapat membuat banyak transaksi.
Namun setiap transaksi hanya dimiliki oleh satu user.

B. Transactions → Transaction_Items

- Type: One-to-Many
- Description:
Satu transaksi dapat memiliki beberapa item (setiap kendaraan yang dibeli).
Setiap transaction_item hanya terkait dengan satu transaksi.

C. Vehicles → Transaction_Items

- Type: One-to-Many
- Description:
Satu kendaraan dapat muncul di banyak transaksi.
Setiap transaction_item hanya terkait dengan satu kendaraan.

D. Transactions → Payment_Detail

- Type: One-to-One
- Description:
Setiap transaksi memiliki satu detail pembayaran.
Satu payment_detail hanya berlaku untuk satu transaksi.

E. Vehicle_Types → Vehicles

- Type: One-to-Many
- Description:
Satu tipe kendaraan (SUV, MPV, motor sport, dll) dapat dimiliki oleh banyak kendaraan.
Satu kendaraan hanya mempunyai satu tipe.

F. Brands → Vehicles

- Type: One-to-Many
- Description:
Satu brand dapat memiliki banyak kendaraan.
Satu kendaraan hanya menggunakan satu brand.

G. Users ↔ Vehicles (using Favorites)

- Type: Many-to-Many
- Description: Satu user dapat mem-favorite banyak kendaraan Satu kendaraan dapat difavoritkan oleh banyak user.
Relasi many-to-many ini di-handle oleh tabel favorites.

## Integrity Constraints

- email pada users harus unik.

- price pada vehicles dan transaction_items harus bernilai positif.

- stock kendaraan tidak boleh negatif.

- foreign key pada semua tabel wajib konsisten (ON DELETE CASCADE atau RESTRICT sesuai fungsi).

- Tidak boleh ada duplikasi favorit karena PK pada favorites adalah kombinasi (user_id, vehicle_id).

- Tabel payment_detail memiliki UNIQUE(transaction_id) agar satu transaksi hanya memiliki satu pembayaran.

## Additional Notes

- Tabel favorites berfungsi menangani relasi many-to-many dan dipakai untuk fitur “Most Favorited Vehicles”.

- Model database dibuat terpisah per entitas agar lebih terstruktur dan mengikuti prinsip normalisasi.

- Dengan adanya tabel transaction_items, sistem mendukung multi-item per transaksi.

- Tabel types dan brands digunakan sebagai master data agar kendaraan lebih konsisten dan tidak terjadi duplikasi nama.

- Struktur ini mendukung ekspansi fitur seperti:
    - laporan penjualan,

    - rekomendasi kendaraan berdasarkan favorit,

    - filter kendaraan berdasarkan brand/type.