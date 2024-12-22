package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/abiiranathan/go-starter/cmd/app/sqlc"
	"github.com/abiiranathan/go-starter/internal"
	"github.com/go-chi/chi/v5"
)

type userService struct {
	handler *handler
}

// UserRoutes registers all user routes.
func UserRoutes(h *handler) {
	userApi := &userService{handler: h}

	h.router.Route("/users", func(r chi.Router) {
		r.Get("/", userApi.listUsers)
		r.Post("/", userApi.createUser)
		r.Get("/{id}", userApi.getUser)
		r.Put("/{id}", userApi.updateUser)
		r.Delete("/{id}", userApi.deleteUser)
	})
}

func (svc *userService) listUsers(w http.ResponseWriter, r *http.Request) {
	const redisKey = "users"

	// Check if the data is cached in Redis
	data, err := svc.handler.redis.Get(r.Context(), redisKey).Result()
	if err == nil {
		var users []sqlc.User
		err = json.Unmarshal([]byte(data), &users)
		if err != nil {
			svc.handler.SendError(w, r, err)
			return
		}

		w.Header().Set("X-Cache", "HIT")
		svc.handler.Json(w, users, http.StatusOK)
		return
	}

	// If the data is not cached, fetch it from the database
	users, err := svc.handler.querier.ListUsers(r.Context())
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}

	// Cache the data in Redis
	b, err := json.Marshal(users)
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}

	// Cache the data in Redis
	err = svc.handler.redis.Set(r.Context(), redisKey, b, 0).Err()
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}

	w.Header().Set("X-Cache", "MISS")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (svc *userService) createUser(w http.ResponseWriter, r *http.Request) {
	var payload sqlc.CreateUserParams
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}

	// Hash the password before storing it in the database
	hashedPassword, err := internal.HashPassword(payload.PasswordHash)
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}
	payload.PasswordHash = hashedPassword

	// Create the user
	user, err := svc.handler.querier.CreateUser(r.Context(), payload)
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}
	svc.handler.Json(w, user, http.StatusCreated)
}

func (svc *userService) getUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}

	user, err := svc.handler.querier.GetUserByID(r.Context(), int64(id))
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}
	svc.handler.Json(w, user, http.StatusOK)
}

func (svc *userService) updateUser(w http.ResponseWriter, r *http.Request) {
	var payload sqlc.UpdateUserParams
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}

	userId := r.PathValue("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}

	payload.ID = int64(id)
	user, err := svc.handler.querier.UpdateUser(r.Context(), payload)
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}
	svc.handler.Json(w, user, http.StatusOK)
}

func (svc *userService) deleteUser(w http.ResponseWriter, r *http.Request) {
	userId := r.PathValue("id")
	id, err := strconv.Atoi(userId)
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}

	err = svc.handler.querier.DeleteUser(r.Context(), int64(id))
	if err != nil {
		svc.handler.SendError(w, r, err)
		return
	}
	svc.handler.Json(w, nil, http.StatusNoContent)
}
