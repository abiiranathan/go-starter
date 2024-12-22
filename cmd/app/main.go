package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/abiiranathan/go-starter/config"
	"github.com/abiiranathan/go-starter/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func loadConfiguration() (*config.Config, error) {
	var configPath string
	var env string

	flag.StringVar(&configPath, "config", "", "path to config file")
	flag.StringVar(&env, "env", "", "path to env file")
	flag.Parse()

	if configPath != "" && env != "" {
		return nil, fmt.Errorf("only one of -config or -env flags can be provided")
	}

	if configPath != "" {
		c, err := config.LoadFromYAML(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %v", err)
		}
		return c, nil
	}

	if env != "" {
		c, err := config.LoadFromEnv(env)
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %v", err)
		}
		return c, nil
	}

	// If .env is available, load from it
	stat, err := os.Stat(".env")
	if err == nil && !stat.IsDir() {
		c, err := config.LoadFromEnv(".env")
		if err != nil {
			return nil, fmt.Errorf("failed to load config: %v", err)
		}
		return c, nil
	}

	return nil, fmt.Errorf("no configuration provided")
}

func connectToDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(cfg.DatabaseURL)
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

func main() {
	cfg := Must(loadConfiguration())
	pool := Must(connectToDatabase(cfg))
	defer pool.Close()

	log.Println("connected to database")

	// Initialize the chi router
	mux := chi.NewRouter()

	// Initialize the handler and inject the sqlc.Querier
	handler := handler.NewHandler(mux, pool)

	// Setup the routes
	handler.SetupRoutes()

	// Start the server
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Println("server is running on port", cfg.Port)

	log.Fatalln(http.ListenAndServe(addr, mux))
}
