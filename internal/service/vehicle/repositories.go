package vehicle

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type VehicleRepository interface {
	Create(ctx context.Context, v Vehicle) error
	UpdatePartial(ctx context.Context, id int, updates map[string]interface{}) error
	GetAll(ctx context.Context) ([]VehicleResponse, error)
	GetByID(ctx context.Context, id int) (VehicleResponse, error)
	GetLowStock(ctx context.Context) ([]VehicleLowStock, error)
}

type vehicleRepository struct {
	db *pgxpool.Pool
}

func NewVehicleRepository(db *pgxpool.Pool) VehicleRepository {
	return &vehicleRepository{db}
}

func (r *vehicleRepository) Create(ctx context.Context, v Vehicle) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO vehicles (type_id, brand_id, name, fuel_type, transmission, price, stock)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
	`,
		v.TypeID, v.BrandID, v.Name, v.FuelType, v.Transmission, v.Price, v.Stock)

	return err
}

func (r *vehicleRepository) UpdatePartial(ctx context.Context, id int, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	q := "UPDATE vehicles SET "
	args := []interface{}{}
	i := 1

	for col, val := range updates {
		if i > 1 {
			q += ", "
		}
		q += col + " = $" + fmt.Sprint(i)
		args = append(args, val)
		i++
	}

	q += " WHERE vehicle_id = $" + fmt.Sprint(i)
	args = append(args, id)

	_, err := r.db.Exec(ctx, q, args...)
	return err
}

func (r *vehicleRepository) GetAll(ctx context.Context) ([]VehicleResponse, error) {
	rows, err := r.db.Query(ctx, `
		SELECT v.vehicle_id, v.name, b.brand_name, t.type_name,
			v.fuel_type, v.transmission, v.price, v.stock
		FROM vehicles v
		JOIN brands b ON v.brand_id = b.brand_id
		JOIN vehicle_types t ON v.type_id = t.type_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []VehicleResponse

	for rows.Next() {
		var v VehicleResponse
		rows.Scan(
			&v.ID, &v.Name, &v.Brand, &v.Type,
			&v.FuelType, &v.Transmission, &v.Price, &v.Stock,
		)
		list = append(list, v)
	}

	return list, nil
}

func (r *vehicleRepository) GetByID(ctx context.Context, id int) (VehicleResponse, error) {
	var v VehicleResponse
	err := r.db.QueryRow(ctx, `
		SELECT v.vehicle_id, v.name, b.brand_name, t.type_name,
			v.fuel_type, v.transmission, v.price, v.stock
		FROM vehicles v
		JOIN brands b ON v.brand_id = b.brand_id
		JOIN vehicle_types t ON v.type_id = t.type_id
		WHERE v.vehicle_id = $1
	`, id).Scan(
		&v.ID, &v.Name, &v.Brand, &v.Type,
		&v.FuelType, &v.Transmission, &v.Price, &v.Stock,
	)

	return v, err
}

func (r *vehicleRepository) GetLowStock(ctx context.Context) ([]VehicleLowStock, error) {
	rows, err := r.db.Query(ctx, `
		SELECT v.vehicle_id, v.name, b.brand_name, t.type_name, v.stock, v.price
		FROM vehicles v
		JOIN brands b ON v.brand_id = b.brand_id
		JOIN vehicle_types t ON v.type_id = t.type_id
		WHERE v.stock <= 5
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []VehicleLowStock
	for rows.Next() {
		var lv VehicleLowStock
		rows.Scan(&lv.ID, &lv.Name, &lv.Brand, &lv.Type, &lv.Stock, &lv.Price)
		list = append(list, lv)
	}

	return list, nil
}
