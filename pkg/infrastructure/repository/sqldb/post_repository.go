package sqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/Tamiris15/university-blog/pkg/domain"
)

type PostRepository struct {
	databaseConnection *sql.DB
}

func NewPostRepository(databaseConnection *sql.DB) *PostRepository {
	return &PostRepository{databaseConnection: databaseConnection}
}

func (repo *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	query := `INSERT INTO posts (title, content, author_id, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	post.CreatedAt = now
	post.UpdatedAt = now

	result, err := repo.databaseConnection.ExecContext(ctx, query,
		post.Title,
		post.Content,
		post.AuthorID,
		post.CreatedAt,
		post.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = id
	return nil
}

func (repo *PostRepository) GetByID(ctx context.Context, id int64) (*domain.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at 
			  FROM posts WHERE id = ?`

	post := &domain.Post{}
	err := repo.databaseConnection.QueryRowContext(ctx, query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (repo *PostRepository) GetAll(ctx context.Context) ([]*domain.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at FROM posts`

	rows, err := repo.databaseConnection.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (repo *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	query := `UPDATE posts 
			  SET title = ?, content = ?, author_id = ?, updated_at = ? 
			  WHERE id = ?`

	post.UpdatedAt = time.Now()
	_, err := repo.databaseConnection.ExecContext(ctx, query,
		post.Title,
		post.Content,
		post.AuthorID,
		post.UpdatedAt,
		post.ID,
	)
	return err
}

func (repo *PostRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM posts WHERE id = ?`
	_, err := repo.databaseConnection.ExecContext(ctx, query, id)
	return err
}

func (repo *PostRepository) GetByAuthorID(ctx context.Context, authorID int64) ([]*domain.Post, error) {
	query := `SELECT id, title, content, author_id, created_at, updated_at 
			  FROM posts WHERE author_id = ?`

	rows, err := repo.databaseConnection.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}
