package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/abiiranathan/go-starter/config"
	"github.com/abiiranathan/go-starter/handlers"
	"github.com/abiiranathan/go-starter/views"
	"github.com/abiiranathan/rex"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var tmpl = Must(views.Templates(template.FuncMap{}))
var tmplOptions = []rex.RouterOption{
	rex.WithTemplates(tmpl),
	rex.BaseLayout("templates/layout.html"),
	rex.ErrorTemplate("templates/error.html"),
	rex.ContentBlock("Content"),
}

func main() {
	cfg := Must(config.Load())
	pool := Must(ConnectToDatabase(cfg.DatabaseURL))
	defer pool.Close()

	log.Println("Connected to database")

	// Router options
	routerOptions := []rex.RouterOption{}
	routerOptions = append(routerOptions, tmplOptions...)

	// Initialize the router
	router := rex.NewRouter(routerOptions...)

	// Init redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:             cfg.RedisClientAddr,
		PoolSize:         10,
		MaxRetries:       3,
		DisableIndentity: true,
	})

	defer redisClient.Close()

	// Initialize the handler.
	handler := handlers.NewHandler(router, pool, redisClient)

	// Setup the routes
	handler.SetupRoutes()

	// Start the server
	fmt.Printf("Server is running on http://0.0.0.0:%s\n", cfg.Port)

	log.Fatalln(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}

// ConnectToDatabase Connects to the postgres database with
// the provided DATABASE_URL in the configuration.
// The connection timeout is 10 seconds.
func ConnectToDatabase(connString string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Must[T any](t T, err error) T {
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return t
}
