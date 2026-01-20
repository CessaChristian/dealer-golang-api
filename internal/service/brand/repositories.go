package brand

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, name string) error
	GetAll(ctx context.Context) ([]Brand, error)
	GetByName(ctx context.Context, name string) (Brand, error)
	Ensure(ctx context.Context, name string) (int, error)
}

type brandRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &brandRepository{db}
}

func (r *brandRepository) Create(ctx context.Context, name string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO brands (brand_name) VALUES ($1)`,
		name,
	)
	return err
}

func (r *brandRepository) GetAll(ctx context.Context) ([]Brand, error) {
	rows, err := r.db.Query(ctx,
		`SELECT brand_id, brand_name FROM brands ORDER BY brand_name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Brand
	for rows.Next() {
		var b Brand
		rows.Scan(&b.ID, &b.Name)
		list = append(list, b)
	}

	return list, nil
}

func (r *brandRepository) GetByName(ctx context.Context, name string) (Brand, error) {
	var b Brand
	err := r.db.QueryRow(ctx,
		`SELECT brand_id, brand_name FROM brands WHERE LOWER(brand_name) = LOWER($1)`,
		name,
	).Scan(&b.ID, &b.Name)

	return b, err
}

func (r *brandRepository) Ensure(ctx context.Context, name string) (int, error) {
	// normalisasi dulu
	name = strings.TrimSpace(strings.ToLower(name))

	// cek apakah sudah ada
	b, err := r.GetByName(ctx, name)
	if err == nil {
		return b.ID, nil
	}

	// buat baru
	_, err = r.db.Exec(ctx,
		`INSERT INTO brands (brand_name) VALUES ($1)`,
		name,
	)
	if err != nil {
		return 0, err
	}

	// ambil lagi ID barunya
	b, err = r.GetByName(ctx, name)
	if err != nil {
		return 0, err
	}

	return b.ID, nil
}
