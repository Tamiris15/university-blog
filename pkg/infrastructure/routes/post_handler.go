package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Tamiris15/university-blog/pkg/application"
	"github.com/Tamiris15/university-blog/pkg/domain"
)

type PostHandler struct {
	postService *application.PostService
}

func NewPostHandler(postService *application.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (ph *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	var post domain.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Некорректный формат данных для создания публикации", http.StatusBadRequest)
		return
	}

	if err := ph.postService.CreatePost(r.Context(), &post); err != nil {
		http.Error(w, "Произошла ошибка при создании публикации", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (ph *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор публикации", http.StatusBadRequest)
		return
	}

	post, err := ph.postService.GetPostByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Произошла ошибка при получении публикации", http.StatusInternalServerError)
		return
	}
	if post == nil {
		http.Error(w, "Публикация не найдена", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (ph *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := ph.postService.GetAllPosts(r.Context())
	if err != nil {
		http.Error(w, "Произошла ошибка при получении всех публикаций", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

func (ph *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор публикации", http.StatusBadRequest)
		return
	}

	var post domain.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Некорректный формат данных для обновления публикации", http.StatusBadRequest)
		return
	}
	post.ID = id

	if err := ph.postService.UpdatePost(r.Context(), &post); err != nil {
		http.Error(w, "Произошла ошибка при обновлении публикации", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(post)
}

func (ph *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор публикации", http.StatusBadRequest)
		return
	}

	if err := ph.postService.DeletePost(r.Context(), id); err != nil {
		http.Error(w, "Произошла ошибка при удалении публикации", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (ph *PostHandler) GetByAuthorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authorID, err := strconv.ParseInt(vars["author_id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор автора", http.StatusBadRequest)
		return
	}

	posts, err := ph.postService.GetPostsByAuthorID(r.Context(), authorID)
	if err != nil {
		http.Error(w, "Произошла ошибка при получении публикаций по идентификатору автора", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(posts)
}
