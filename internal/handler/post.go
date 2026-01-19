package handler

import (
	"encoding/json"
	"lemara_blog/internal/domain"
	"lemara_blog/internal/service"
	"net/http"
)

type PostHandler struct {
	service service.PostService
}

func NewPostHandler(service service.PostService) *PostHandler {
	return &PostHandler{service: service}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userId := GetUserIDFromContext(r.Context())
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Получаем данные из запроса
	var updateReq domain.PostCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Присваиваем пользователя из контекста
	updateReq.Author = userId

	post, err := h.service.CreatePost(r.Context(), &updateReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Записываем и возвращаем ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(post)
}
