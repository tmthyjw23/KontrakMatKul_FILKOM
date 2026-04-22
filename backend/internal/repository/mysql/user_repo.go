package mysql

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"sistemkontrakmatkul/backend/internal/domain/models"
	"sistemkontrakmatkul/backend/internal/domain/repositories"
)

type UserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

var _ repositories.UserRepository = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) GetByStudentNumber(ctx context.Context, studentNumber string) (*models.User, error) {
	const query = `SELECT id, student_number, full_name, email, password, role, max_sks, created_at, updated_at FROM users WHERE student_number = ?`
	
	var user models.User
	err := r.db.QueryRowContext(ctx, query, studentNumber).Scan(
		&user.ID, &user.StudentNumber, &user.FullName, &user.Email, &user.Password, &user.Role, &user.MaxSKS, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil user and nil error to distinguish from system error
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint64) (*models.User, error) {
	const query = `SELECT id, student_number, full_name, email, password, role, max_sks, created_at, updated_at FROM users WHERE id = ?`
	
	var user models.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.StudentNumber, &user.FullName, &user.Email, &user.Password, &user.Role, &user.MaxSKS, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	const query = `INSERT INTO users (student_number, full_name, email, password, role, max_sks) VALUES (?, ?, ?, ?, ?, ?)`
	
	_, err := r.db.ExecContext(ctx, query, user.StudentNumber, user.FullName, user.Email, user.Password, user.Role, user.MaxSKS)
	return err
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	const query = `UPDATE users SET full_name = ?, email = ?, password = ?, role = ?, max_sks = ? WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, user.FullName, user.Email, user.Password, user.Role, user.MaxSKS, user.ID)
	return err
}

func (r *UserRepository) Delete(ctx context.Context, id uint64) error {
	const query = `DELETE FROM users WHERE id = ?`
	
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
