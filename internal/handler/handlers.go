package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/abiiranathan/go-starter/cmd/app/sqlc"
	"github.com/abiiranathan/go-starter/templates"
	"github.com/abiiranathan/rex"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Handler interface {
	// SetupRoutes sets up all routes for the application
	SetupRoutes()

	// Conn returns the database connection pool
	Pool() *pgxpool.Pool
}

type handler struct {
	router  *rex.Router   // chi router
	querier sqlc.Querier  // sqlc querier
	pool    *pgxpool.Pool // database connection pool
	redis   *redis.Client // redis client
}

func NewHandler(r *rex.Router, pool *pgxpool.Pool) Handler {
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

func (h *handler) SetupRoutes() {
	AttachMiddleware(h.router)
	StaticRoutes(h.router)

	IndexRoute(h)
	UserRoutes(h)

	// Add more routes here
}

func RenderPage(c *rex.Context, component templ.Component) error {
	c.Response.Header().Set("Content-Type", "text/html")
	c.WriteHeader(http.StatusOK)
	return templates.Layout(component).Render(c.Request.Context(), c.Response)
}
