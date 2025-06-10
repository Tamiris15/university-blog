package sqldb

import (
	"context"
	"database/sql"
	"time"

	"github.com/Tamiris15/university-blog/pkg/domain"
)

type CommentRepository struct {
	databaseConnection *sql.DB
}

func NewCommentRepository(databaseConnection *sql.DB) *CommentRepository {
	return &CommentRepository{databaseConnection: databaseConnection}
}

func (repo *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	query := `INSERT INTO comments (content, post_id, author_id, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	comment.CreatedAt = now
	comment.UpdatedAt = now

	result, err := repo.databaseConnection.ExecContext(ctx, query,
		comment.Content,
		comment.PostID,
		comment.AuthorID,
		comment.CreatedAt,
		comment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	comment.ID = id
	return nil
}

func (repo *CommentRepository) GetByID(ctx context.Context, id int64) (*domain.Comment, error) {
	query := `SELECT id, content, post_id, author_id, created_at, updated_at 
			  FROM comments WHERE id = ?`

	comment := &domain.Comment{}
	err := repo.databaseConnection.QueryRowContext(ctx, query, id).Scan(
		&comment.ID,
		&comment.Content,
		&comment.PostID,
		&comment.AuthorID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (repo *CommentRepository) GetByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error) {
	query := `SELECT id, content, post_id, author_id, created_at, updated_at 
			  FROM comments WHERE post_id = ?`

	rows, err := repo.databaseConnection.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		comment := &domain.Comment{}
		err := rows.Scan(
			&comment.ID,
			&comment.Content,
			&comment.PostID,
			&comment.AuthorID,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (repo *CommentRepository) Update(ctx context.Context, comment *domain.Comment) error {
	query := `UPDATE comments 
			  SET content = ?, post_id = ?, author_id = ?, updated_at = ? 
			  WHERE id = ?`

	comment.UpdatedAt = time.Now()
	_, err := repo.databaseConnection.ExecContext(ctx, query,
		comment.Content,
		comment.PostID,
		comment.AuthorID,
		comment.UpdatedAt,
		comment.ID,
	)
	return err
}

func (repo *CommentRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM comments WHERE id = ?`
	_, err := repo.databaseConnection.ExecContext(ctx, query, id)
	return err
}
