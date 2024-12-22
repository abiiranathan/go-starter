package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/abiiranathan/go-starter/cmd/app/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler interface {
	// SendError sends an error response to the client.
	// Depending on the request's Accept header, it sends either a JSON or HTML response
	// Default status code is 500, but it can be overridden by passing a custom status code
	SendError(w http.ResponseWriter, r *http.Request, err error, status ...int)

	// Json sends a JSON response to the client
	Json(w http.ResponseWriter, data interface{}, status int)

	//Send sends html response to the client
	Send(w http.ResponseWriter, data []byte, status int)

	// SetupRoutes sets up all routes for the application
	SetupRoutes()

	// Conn returns the database connection pool
	Pool() *pgxpool.Pool
}

type handler struct {
	router  *chi.Mux      // chi router
	querier sqlc.Querier  // sqlc querier
	pool    *pgxpool.Pool // database connection pool
	redis   *redis.Client // redis client
}

func NewHandler(r *chi.Mux, pool *pgxpool.Pool) Handler {
	querier := sqlc.New(pool)

	redisClient := redis.NewClient(&redis.Options{
		Addr:             "localhost:6379",
		PoolSize:         10,
		MaxRetries:       3,
		DisableIndentity: true,
	})

	return &handler{
		router:  r,
		querier: querier,
		pool:    pool,
		redis:   redisClient,
	}
}

func (h *handler) Pool() *pgxpool.Pool {
	return h.pool
}

func (h *handler) SendError(w http.ResponseWriter, r *http.Request, err error, status ...int) {
	code := http.StatusInternalServerError
	if len(status) > 0 {
		code = status[0]
	}

	accept := r.Header.Get("Accept")
	if strings.Contains(accept, "application/json") {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"error": "` + err.Error() + `"}`))
		return
	}

	http.Error(w, err.Error(), code)
}

func (h *handler) Json(w http.ResponseWriter, data interface{}, status int) {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(data)
}

func (h *handler) Send(w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(status)
	w.Write(data)
}

func (h *handler) SetupRoutes() {
	UserRoutes(h)

	// Add more routes here
}
