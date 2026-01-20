package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]UserResponse, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	GetByID(ctx context.Context, id int) (User, error)
	Create(ctx context.Context, u User) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(ctx context.Context, u User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (name, email, password, role)
		 VALUES ($1, $2, $3, $4)`,
		u.Name, u.Email, u.Password, u.Role,
	)
	return err
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`SELECT user_id, name, email, password, role, created_at, updated_at
		 FROM users WHERE email=$1`,
		email,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

func (r *userRepository) GetByID(ctx context.Context, id int) (User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`SELECT user_id, name, email, password, role, created_at, updated_at
		 FROM users WHERE user_id=$1`,
		id,
	).Scan(
		&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt,
	)
	return u, err
}

func (r *userRepository) GetAll(ctx context.Context) ([]UserResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT user_id, name, email, role FROM users`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []UserResponse
	for rows.Next() {
		var u UserResponse
		rows.Scan(&u.ID, &u.Name, &u.Email, &u.Role)
		users = append(users, u)
	}

	return users, nil
}
