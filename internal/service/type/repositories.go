package vtype

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, name string) error
	GetAll(ctx context.Context) ([]Type, error)
	GetByName(ctx context.Context, name string) (Type, error)
	Ensure(ctx context.Context, name string) (int, error)
}

type typeRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &typeRepository{db}
}

func (r *typeRepository) Create(ctx context.Context, name string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO vehicle_types (type_name) VALUES ($1)`,
		name,
	)
	return err
}

func (r *typeRepository) GetAll(ctx context.Context) ([]Type, error) {
	rows, err := r.db.Query(ctx,
		`SELECT type_id, type_name 
		 FROM vehicle_types 
		 ORDER BY type_name ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Type
	for rows.Next() {
		var t Type
		rows.Scan(&t.ID, &t.Name)
		list = append(list, t)
	}

	return list, nil
}

func (r *typeRepository) GetByName(ctx context.Context, name string) (Type, error) {
	var t Type
	err := r.db.QueryRow(ctx,
		`SELECT type_id, type_name
		 FROM vehicle_types 
		 WHERE LOWER(type_name) = LOWER($1)`,
		name,
	).Scan(&t.ID, &t.Name)
	return t, err
}

func (r *typeRepository) Ensure(ctx context.Context, name string) (int, error) {
	name = strings.TrimSpace(strings.ToLower(name))

	t, err := r.GetByName(ctx, name)
	if err == nil {
		return t.ID, nil
	}

	_, err = r.db.Exec(ctx,
		`INSERT INTO vehicle_types (type_name) VALUES ($1)`,
		name,
	)
	if err != nil {
		return 0, err
	}

	t, err = r.GetByName(ctx, name)
	if err != nil {
		return 0, err
	}

	return t.ID, nil
}

