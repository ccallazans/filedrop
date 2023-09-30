package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/labstack/gommon/log"
)

type UserStore interface {
	DB() *sql.DB
	FindAll(ctx context.Context) ([]*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	Exists(ctx context.Context, email string) (bool, error)
	Save(ctx context.Context, user *domain.User) error
	DeleteByID(ctx context.Context, id string) error
}

// IMPLEMENTATION

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) UserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (r *PostgresUserStore) DB() *sql.DB {
	return r.db
}

func (r *PostgresUserStore) FindByID(ctx context.Context, id string) (*domain.User, error) {
	query := `
		SELECT id, first_name, email, password, role_id, created_at, updated_at
		FROM users
		WHERE id = $1
		LIMIT 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to prepare query")
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	var user domain.User
	if err := row.Scan(&user.ID, &user.FirstName, &user.Email, &user.Password, &user.RoleID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Error(err)
		return nil, errors.New("failed to scan row")
	}

	return &user, nil
}

func (r *PostgresUserStore) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, first_name, email, password, role_id, created_at, updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to prepare query")
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)

	var user domain.User
	if err := row.Scan(&user.ID, &user.FirstName, &user.Email, &user.Password, &user.RoleID, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Error(err)
		return nil, errors.New("failed to scan row")
	}

	return &user, nil
}

func (r *PostgresUserStore) Exists(ctx context.Context, email string) (bool, error) {
	query := `
		SELECT 1
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return false, errors.New("failed to prepare query")
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, email)

	var exists bool
	if err := row.Scan(&exists); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		log.Error(err)
		return false, errors.New("failed to scan row")
	}

	return exists, nil
}

func (r *PostgresUserStore) FindAll(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, first_name, email, password, role_id, created_at, updated_at
		FROM users
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to prepare query")
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to query users")
	}
	defer rows.Close()

	users := []*domain.User{}
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.Email, &user.Password, &user.RoleID, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Error(err)
			return nil, errors.New("failed to scan row")
		}

		users = append(users, &user)
	}

	return users, nil
}

func (r *PostgresUserStore) Save(ctx context.Context, user *domain.User) error {
	query := `
			INSERT INTO users (first_name, email, password, role_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return errors.New("failed to prepare query")
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, user.FirstName, user.Email, user.Password, user.RoleID, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		log.Error(err)
		return errors.New("failed to execute query")
	}

	return nil
}

func (r *PostgresUserStore) DeleteByID(ctx context.Context, id string) error {
	query := `
		DELETE FROM users
		WHERE id = $7
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return errors.New("failed to prepare query")
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, id)
	if err != nil {
		log.Error(err)
		return errors.New("failed to execute query")
	}

	return nil
}
