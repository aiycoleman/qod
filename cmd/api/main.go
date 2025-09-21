// Filename: cmd/api/main.go

package main

import (
	"context"
	"database/sql"
	"flag"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/aiycoleman/qod/internal/data"
	_ "github.com/lib/pq"
)

// configuration holds all the runtime configuration settings for the app.
// Private (non-exportable) to this package (lowercase "configuration")
type configuration struct {
	port    int
	env     string // Application environment
	version string // Version number of the API
	db      struct {
		dsn string
	}
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
}

// Hold dependencies shared across handlers,
// such as config and logger.
type application struct {
	config     configuration
	logger     *slog.Logger
	quoteModel data.QuoteModel
}

// loadConfig reads configuration from command line flags
func loadConfig() configuration {
	var cfg configuration

	// Register CLI flags and bind them to cfg fields.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.StringVar(&cfg.version, "version", "1.0.0", "Application version")

	// Read in the dsn
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://quotes:whyme@localhost/quotes", "PostgreSQL DSN")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate Limiter maximum requests per second")

	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 5, "Rate Limiter maximum burst")

	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// Allow us to access space-seperted origins.
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space seperated)",
		func(val string) error {
			cfg.cors.trustedOrigins = strings.Fields(val)
			return nil
		})

	flag.Parse()

	return cfg
}

// setupLogger configures the application logger based on environment
func setupLogger() *slog.Logger {
	var logger *slog.Logger

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	return logger
}

func openDB(settings configuration) (*sql.DB, error) {
	// open a connection pool
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}

	// set a context to ensure DB operations don't take too long
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test if the connection pool was created
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	// return the connection pool (sql.DB)
	return db, nil

}

// printUB is a small test function, not used in production (testing).
func printUB() string {
	return "Hello, UB!"
}

func main() {
	// greeting := printUB()
	// fmt.Println(greeting)

	// Initialize configuration
	cfg := loadConfig()
	// Initialize logger
	logger := setupLogger()

	// Call to openDB() sets up our connection pool
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// release the database resources before exiting
	defer db.Close()

	logger.Info("database connection pool established")

	// Initialize application struc with dependencies
	app := &application{
		config:     cfg,
		logger:     logger,
		quoteModel: data.QuoteModel{DB: db},
	}

	// Run the application
	if err := app.serve(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
