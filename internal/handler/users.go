package handler

import (
	"encoding/json"
	"net/http"

	"github.com/abiiranathan/go-starter/cmd/app/sqlc"
	"github.com/abiiranathan/go-starter/internal"
	"github.com/abiiranathan/rex"
)

type userService struct {
	handler *handler
}

// UserRoutes registers all user routes.
func UserRoutes(h *handler) {
	userApi := &userService{handler: h}

	users := h.router.Group("/users")

	users.GET("/", userApi.listUsers)
	users.POST("/", userApi.createUser)
	users.GET("/{id}", userApi.getUser)
	users.PUT("/{id}", userApi.updateUser)
	users.DELETE("/{id}", userApi.deleteUser)
}

func (svc *userService) listUsers(c *rex.Context) error {
	const redisKey = "users"

	// Check if the data is cached in Redis
	data, err := svc.handler.redis.Get(c.Request.Context(), redisKey).Result()
	if err == nil {
		var users []sqlc.User
		err = json.Unmarshal([]byte(data), &users)
		if err != nil {
			return err
		}
		c.SetHeader("X-Cache", "HIT")
		return c.JSON(users)
	}

	// If the data is not cached, fetch it from the database
	users, err := svc.handler.querier.ListUsers(c.Request.Context())
	if err != nil {
		return err
	}

	// Cache the data in Redis
	b, err := json.Marshal(users)
	if err != nil {
		return err
	}

	// Cache the data in Redis
	err = svc.handler.redis.Set(c.Request.Context(), redisKey, b, 0).Err()
	if err != nil {
		return err
	}

	c.SetHeader("X-Cache", "MISS")
	return c.JSON(users)
}

func (svc *userService) createUser(c *rex.Context) error {
	var payload sqlc.CreateUserParams
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	// Hash the password before storing it in the database
	hashedPassword, err := internal.HashPassword(payload.PasswordHash)
	if err != nil {
		return err
	}
	payload.PasswordHash = hashedPassword

	// Create the user
	user, err := svc.handler.querier.CreateUser(c.Request.Context(), payload)
	if err != nil {
		return err
	}

	c.WriteHeader(http.StatusCreated)
	return c.JSON(user)
}

func (svc *userService) getUser(c *rex.Context) error {
	userId := c.ParamInt("id")

	user, err := svc.handler.querier.GetUserByID(c.Request.Context(), int64(userId))
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (svc *userService) updateUser(c *rex.Context) error {
	var payload sqlc.UpdateUserParams
	err := c.BodyParser(&payload)
	if err != nil {
		return err
	}

	userId := c.ParamInt("id")

	payload.ID = int64(userId)
	user, err := svc.handler.querier.UpdateUser(c.Request.Context(), payload)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

func (svc *userService) deleteUser(c *rex.Context) error {
	userId := c.ParamInt("id")

	err := svc.handler.querier.DeleteUser(c.Request.Context(), int64(userId))
	if err != nil {
		return err
	}

	c.WriteHeader(http.StatusNoContent)
	return c.JSON(nil)
}
