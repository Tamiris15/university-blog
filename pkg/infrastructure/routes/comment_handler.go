package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Tamiris15/university-blog/pkg/application"
	"github.com/Tamiris15/university-blog/pkg/domain"
)

type CommentHandler struct {
	commentService *application.CommentService
}

func NewCommentHandler(commentService *application.CommentService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
	}
}

func (ch *CommentHandler) Create(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["post_id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор публикации", http.StatusBadRequest)
		return
	}

	var comment domain.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Некорректный формат данных для создания комментария", http.StatusBadRequest)
		return
	}

	comment.PostID = postID

	if err := ch.commentService.CreateComment(r.Context(), &comment); err != nil {
		http.Error(w, "Произошла ошибка при создании комментария", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

func (ch *CommentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор комментария", http.StatusBadRequest)
		return
	}

	comment, err := ch.commentService.GetCommentByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Произошла ошибка при получении комментария", http.StatusInternalServerError)
		return
	}
	if comment == nil {
		http.Error(w, "Комментарий не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(comment)
}

func (ch *CommentHandler) GetByPostID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.ParseInt(vars["post_id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор публикации", http.StatusBadRequest)
		return
	}

	comments, err := ch.commentService.GetCommentsByPostID(r.Context(), postID)
	if err != nil {
		http.Error(w, "Произошла ошибка при получении комментариев для публикации", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comments)
}

func (ch *CommentHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор комментария", http.StatusBadRequest)
		return
	}

	var comment domain.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Некорректный формат данных для обновления комментария", http.StatusBadRequest)
		return
	}
	comment.ID = id

	if err := ch.commentService.UpdateComment(r.Context(), &comment); err != nil {
		http.Error(w, "Произошла ошибка при обновлении комментария", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(comment)
}

func (ch *CommentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор комментария", http.StatusBadRequest)
		return
	}

	if err := ch.commentService.DeleteComment(r.Context(), id); err != nil {
		http.Error(w, "Произошла ошибка при удалении комментария", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
