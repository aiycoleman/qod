// Filename: cmd/api/main.go

package main

import (
	"flag"
	"log/slog"
	"os"
)

// The 'configuration' type is lowercase to
// signal that it is private (non-exportable) to the
// main package
type configuration struct {
	port int
	env  string
}

type application struct {
	config configuration
	logger *slog.Logger
}

// loadConfig reads configuration from command line flags
func loadConfig() configuration {
	var cfg configuration

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.Parse()

	return cfg
}

// setupLogger configures the application logger based on environment
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	return logger
}

func printUB() string {
	return "Hello, UB!"
}

func main() {
	// greeting := printUB()
	// fmt.Println(greeting)
	// Initialize configuration

	cfg := loadConfig()
	// Initialize logger
	logger := setupLogger(cfg.env)
	// Initialize application with dependencies
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Run the application
	app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
