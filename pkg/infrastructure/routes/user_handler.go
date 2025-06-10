package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Tamiris15/university-blog/pkg/application"
	"github.com/Tamiris15/university-blog/pkg/domain"
)

type UserHandler struct {
	userService *application.UserService
}

func NewUserHandler(userService *application.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Printf("Не удалось декодировать данные запроса для пользователя: %v", err)
		http.Error(w, "Некорректный формат данных пользователя", http.StatusBadRequest)
		return
	}

	log.Printf("Получены входящие данные пользователя: %+v", user)

	if err := uh.userService.CreateUser(r.Context(), &user); err != nil {
		log.Printf("Произошла ошибка при регистрации пользователя: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при создании пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор пользователя", http.StatusBadRequest)
		return
	}

	user, err := uh.userService.GetUserByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Произошла ошибка при получении данных пользователя", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Пользователь с указанным ID не найден", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор пользователя", http.StatusBadRequest)
		return
	}

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Некорректный формат данных для обновления пользователя", http.StatusBadRequest)
		return
	}
	user.ID = id

	if err := uh.userService.UpdateUser(r.Context(), &user); err != nil {
		http.Error(w, "Произошла ошибка при обновлении пользователя", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор пользователя", http.StatusBadRequest)
		return
	}

	if err := uh.userService.DeleteUser(r.Context(), id); err != nil {
		http.Error(w, "Произошла ошибка при удалении пользователя", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := uh.userService.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "Произошла ошибка при получении списка пользователей", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}
