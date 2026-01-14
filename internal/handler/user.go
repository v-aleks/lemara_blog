package handler

import (
	"encoding/json"
	"net/http"

	"lemara_blog/internal/domain"
	"lemara_blog/internal/repository"
)

type UserHandler struct {
    userRepo repository.UserRepository
}

func NewUserHandler(userRepo repository.UserRepository) *UserHandler {
    return &UserHandler{userRepo: userRepo}
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
    userID := GetUserIDFromContext(r.Context())
    if userID == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    user, err := h.userRepo.FindByID(r.Context(), userID)
    if err != nil {
        http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
        return
    }

    if user == nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    response := domain.UserResponse{
        ID:        user.ID,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
    userID := GetUserIDFromContext(r.Context())
    if userID == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    var updateReq struct {
        Email    *string `json:"email"`
        Password *string `json:"password"`
        FirstName *string `json:"first_name"`
        LastName *string `json:"last_name"`
    }

    if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // Fetch existing user
    user, err := h.userRepo.FindByID(r.Context(), userID)
    if err != nil || user == nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Update fields if provided
    if updateReq.Email != nil && *updateReq.Email != "" {
        // Check if email is already taken
        existingUser, err := h.userRepo.FindByEmail(r.Context(), *updateReq.Email)
        if err == nil && existingUser != nil && existingUser.ID != userID {
            http.Error(w, "Email already in use", http.StatusBadRequest)
            return
        }
        user.Email = *updateReq.Email
    }

    if updateReq.Password != nil && *updateReq.Password != "" {
        if len(*updateReq.Password) < 8 {
            http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
            return
        }
        // In production, you would hash the password here
        user.PasswordHash = *updateReq.Password
    }

    if err := h.userRepo.Update(r.Context(), user); err != nil {
        http.Error(w, "Failed to update user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "Profile updated successfully"})
}

func (h *UserHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
    userID := GetUserIDFromContext(r.Context())
    if userID == "" {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    // Verify password (optional - for extra security)
    var deleteReq struct {
        Password string `json:"password"`
    }

    if err := json.NewDecoder(r.Body).Decode(&deleteReq); err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    // TODO: Verify password

    if err := h.userRepo.Delete(r.Context(), userID); err != nil {
        http.Error(w, "Failed to delete user", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
