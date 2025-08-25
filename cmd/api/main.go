// Filename: cmd/api/main.go

package main

import (
	"flag"
	"log/slog"
	"os"
)

// configuration holds all the runtime configuration settings for the app.
// Private (non-exportable) to this package (lowercase "configuration")
type configuration struct {
	port    int
	env     string // Application environment
	version string
}

// Hold dependencies shared across handlers,
// such as config and logger.
type application struct {
	config configuration
	logger *slog.Logger
}

// loadConfig reads configuration from command line flags
func loadConfig() configuration {
	var cfg configuration

	// Register CLI flags and bind them to cfg fields.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.StringVar(&cfg.version, "version", "1.0.0", "Application version")
	flag.Parse()

	return cfg
}

// setupLogger configures the application logger based on environment
func setupLogger() *slog.Logger {
	var logger *slog.Logger

	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	return logger
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
	// Initialize application struc with dependencies
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Run the application
	err := app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
