package vehicle

import "time"

type Vehicle struct {
    ID           int       `db:"vehicle_id"`
    TypeID       int       `db:"type_id"`
    BrandID      int       `db:"brand_id"`
    Name         string    `db:"name"`
    FuelType     string    `db:"fuel_type"`
    Transmission string    `db:"transmission"`
    Price        float64   `db:"price"`
    Stock        int       `db:"stock"`
    CreatedAt    time.Time `db:"created_at"`
    UpdatedAt    time.Time `db:"updated_at"`
}

type VehicleLowStock struct {
    ID    int     `db:"vehicle_id"`
    Name  string  `db:"name"`
    Brand string  `db:"brand_name"`
    Type  string  `db:"type_name"`
    Stock int     `db:"stock"`
    Price float64 `db:"price"`
}
