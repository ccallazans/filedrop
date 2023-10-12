package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ccallazans/filedrop/internal/domain"
	"github.com/labstack/gommon/log"
)

type FileStore interface {
	DB() *sql.DB
	FindByID(ctx context.Context, id string) (*domain.File, error)
	FindByHash(ctx context.Context, hash string) (*domain.File, error)
	Exists(ctx context.Context, hash string) (bool, error)
	Save(ctx context.Context, file *domain.File) error
	DeleteByID(ctx context.Context, id string) error
}

// IMPLEMENTATION

type PostgresFileStore struct {
	db *sql.DB
}

func NewPostgresFileStore(db *sql.DB) FileStore {
	return &PostgresFileStore{
		db: db,
	}
}

func (r *PostgresFileStore) DB() *sql.DB {
	return r.db
}

func (r *PostgresFileStore) FindByID(ctx context.Context, id string) (*domain.File, error) {
	query := `
		SELECT id, filename, password, location, hash, is_active, user_id, created_at, updated_at
		FROM files
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

	var file domain.File
	if err := row.Scan(&file.ID, &file.Filename, &file.Password, &file.Location, &file.Hash, &file.IsActive, &file.UserID, &file.CreatedAt, &file.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Error(err)
		return nil, errors.New("failed to scan row")
	}

	return &file, nil
}

func (r *PostgresFileStore) FindByHash(ctx context.Context, hash string) (*domain.File, error) {
	query := `
		SELECT id, filename, password, location, hash, is_active, user_id, created_at, updated_at
		FROM files
		WHERE hash = $1
		LIMIT 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return nil, errors.New("failed to prepare query")
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, hash)

	var file domain.File
	if err := row.Scan(&file.ID, &file.Filename, &file.Password, &file.Location, &file.Hash, &file.IsActive, &file.UserID, &file.CreatedAt, &file.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		log.Error(err)
		return nil, errors.New("failed to scan row")
	}

	return &file, nil
}

func (r *PostgresFileStore) Exists(ctx context.Context, hash string) (bool, error) {
	query := `
		SELECT 1
		FROM files
		WHERE hash = $1
		LIMIT 1
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return false, errors.New("failed to prepare query")
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, hash)

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

func (r *PostgresFileStore) Save(ctx context.Context, file *domain.File) error {
	query := `
			INSERT INTO files (id, filename, password, location, hash, is_active, user_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		log.Error(err)
		return errors.New("failed to prepare query")
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, file.ID, file.Filename, file.Password, file.Location, file.Hash, file.IsActive, file.UserID, time.Now(), time.Now())
	if err != nil {
		log.Error(err)
		return errors.New("failed to execute query")
	}

	return nil
}

func (r *PostgresFileStore) DeleteByID(ctx context.Context, id string) error {
	query := `
		DELETE FROM files
		WHERE id = $1
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
