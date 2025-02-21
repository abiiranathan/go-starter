package handlers

import (
	"github.com/abiiranathan/go-starter/sqlc"
	"github.com/abiiranathan/rex"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler interface {
	// SetupRoutes sets up all routes for the application
	SetupRoutes()

	// Pool returns the database connection pool
	Pool() *pgxpool.Pool
}

// Struct implementation for the routes.Handler interface.
type handler struct {
	router  *rex.Router   // chi router
	querier sqlc.Querier  // sqlc querier
	pool    *pgxpool.Pool // database connection pool
	redis   *redis.Client // redis client
}

func NewHandler(r *rex.Router, pool *pgxpool.Pool, redisClient *redis.Client) Handler {
	querier := sqlc.New(pool)

	return &handler{
		router:  r,
		querier: querier,
		pool:    pool,
		redis:   redisClient,
	}
}

// Returns the database connection pool.
func (h *handler) Pool() *pgxpool.Pool {
	return h.pool
}

// Register all application routes here.
func (h *handler) SetupRoutes() {
	AttachMiddleware(h.router)
	StaticRoutes(h.router)

	IndexRoute(h)
	UserRoutes(h)
}
