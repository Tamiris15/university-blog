package sqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/Tamiris15/university-blog/pkg/domain"
)

type UserRepository struct {
	databaseConnection *sql.DB
}

func NewUserRepository(databaseConnection *sql.DB) *UserRepository {
	return &UserRepository{databaseConnection: databaseConnection}
}

func (repo *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, email, password, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	result, err := repo.databaseConnection.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (repo *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at 
			  FROM users WHERE id = ?`

	user := &domain.User{}
	err := repo.databaseConnection.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at 
			  FROM users WHERE email = ?`

	user := &domain.User{}
	err := repo.databaseConnection.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users 
			  SET username = ?, email = ?, password = ?, updated_at = ? 
			  WHERE id = ?`

	user.UpdatedAt = time.Now()
	_, err := repo.databaseConnection.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.UpdatedAt,
		user.ID,
	)
	return err
}

func (repo *UserRepository) Delete(ctx context.Context, id int64) error {
	// Начинаем транзакцию
	transaction, err := repo.databaseConnection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	// Удаляем пользователя (посты удалятся автоматически благодаря ON DELETE CASCADE)
	_, err = transaction.ExecContext(ctx, "DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}

	// Подтверждаем транзакцию
	return transaction.Commit()
}

func (repo *UserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `SELECT id, username, email, password, created_at, updated_at FROM users`

	rows, err := repo.databaseConnection.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
