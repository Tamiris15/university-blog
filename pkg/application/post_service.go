package application

import (
	"context"
	"errors"

	"github.com/Tamiris15/university-blog/pkg/domain"
)

type PostService struct {
	postRepository domain.PostRepository
}

func NewPostService(postRepository domain.PostRepository) *PostService {
	return &PostService{
		postRepository: postRepository,
	}
}

func (s *PostService) CreatePost(ctx context.Context, post *domain.Post) error {
	if post.Title == "" || post.Content == "" || post.AuthorID == 0 {
		return errors.New("Недопустимые или отсутствующие данные публикации")
	}
	return s.postRepository.Create(ctx, post)
}

func (s *PostService) GetPostByID(ctx context.Context, id int64) (*domain.Post, error) {
	return s.postRepository.GetByID(ctx, id)
}

func (s *PostService) GetAllPosts(ctx context.Context) ([]*domain.Post, error) {
	return s.postRepository.GetAll(ctx)
}

func (s *PostService) UpdatePost(ctx context.Context, post *domain.Post) error {
	if post.ID == 0 {
		return errors.New("Недопустимый идентификатор публикации")
	}
	return s.postRepository.Update(ctx, post)
}

func (s *PostService) DeletePost(ctx context.Context, id int64) error {
	return s.postRepository.Delete(ctx, id)
}

func (s *PostService) GetPostsByAuthorID(ctx context.Context, authorID int64) ([]*domain.Post, error) {
	return s.postRepository.GetByAuthorID(ctx, authorID)
}
