package application

import (
	"context"
	"errors"

	"github.com/Tamiris15/university-blog/pkg/domain"
)

type CommentService struct {
	commentRepository domain.CommentRepository
}

func NewCommentService(commentRepository domain.CommentRepository) *CommentService {
	return &CommentService{
		commentRepository: commentRepository,
	}
}

func (s *CommentService) CreateComment(ctx context.Context, comment *domain.Comment) error {
	if comment.Content == "" || comment.PostID == 0 || comment.AuthorID == 0 {
		return errors.New("Недопустимые или отсутствующие данные комментария")
	}
	return s.commentRepository.Create(ctx, comment)
}

func (s *CommentService) GetCommentByID(ctx context.Context, id int64) (*domain.Comment, error) {
	return s.commentRepository.GetByID(ctx, id)
}

func (s *CommentService) GetCommentsByPostID(ctx context.Context, postID int64) ([]*domain.Comment, error) {
	return s.commentRepository.GetByPostID(ctx, postID)
}

func (s *CommentService) UpdateComment(ctx context.Context, comment *domain.Comment) error {
	if comment.ID == 0 {
		return errors.New("Недопустимый идентификатор комментария")
	}
	return s.commentRepository.Update(ctx, comment)
}

func (s *CommentService) DeleteComment(ctx context.Context, id int64) error {
	return s.commentRepository.Delete(ctx, id)
}
