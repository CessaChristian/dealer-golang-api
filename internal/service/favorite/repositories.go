package favorite

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	AddFavorite(ctx context.Context, fav Favorite) error
	GetFavoritesByUser(ctx context.Context, userID int) ([]FavoriteReport, error)
	GetAllFavoritesAdmin(ctx context.Context) ([]FavoriteReport, error)
}

type favoriteRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &favoriteRepository{db}
}

func (r *favoriteRepository) AddFavorite(ctx context.Context, fav Favorite) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO favorites (user_id, vehicle_id, created_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, vehicle_id) DO NOTHING
	`, fav.UserID, fav.VehicleID)
	return err
}

func (r *favoriteRepository) GetFavoritesByUser(ctx context.Context, userID int) ([]FavoriteReport, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			u.name AS name,
			v.name AS vehicle_name,
			b.brand_name,
			t.type_name,
			v.price
		FROM favorites f
		JOIN users u ON f.user_id = u.user_id
		JOIN vehicles v ON f.vehicle_id = v.vehicle_id
		JOIN brands b ON v.brand_id = b.brand_id
		JOIN vehicle_types t ON v.type_id = t.type_id
		WHERE f.user_id = $1
		ORDER BY v.name ASC
	`, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []FavoriteReport

	for rows.Next() {
		var fr FavoriteReport
		rows.Scan(&fr.Name, &fr.VehicleName, &fr.Brand, &fr.Type, &fr.Price)
		list = append(list, fr)
	}

	return list, nil
}

func (r *favoriteRepository) GetAllFavoritesAdmin(ctx context.Context) ([]FavoriteReport, error) {
	rows, err := r.db.Query(ctx, `
		SELECT 
			u.name AS name,
			v.name AS vehicle_name,
			b.brand_name,
			t.type_name,
			v.price
		FROM favorites f
		JOIN users u ON f.user_id = u.user_id
		JOIN vehicles v ON f.vehicle_id = v.vehicle_id
		JOIN brands b ON v.brand_id = b.brand_id
		JOIN vehicle_types t ON v.type_id = t.type_id
		ORDER BY u.name ASC, v.name ASC
	`)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []FavoriteReport

	for rows.Next() {
		var fr FavoriteReport
		rows.Scan(&fr.Name, &fr.VehicleName, &fr.Brand, &fr.Type, &fr.Price)
		list = append(list, fr)
	}

	return list, nil
}
