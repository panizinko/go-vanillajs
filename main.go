package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"adrianzinko.com/gomovies/data"
	"adrianzinko.com/gomovies/handlers"
	"adrianzinko.com/gomovies/logger"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func initializeLogger() *logger.Logger {
	logger, err := logger.NewLogger("log.txt")
	logger.Info("Logger initialized")

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Close()

	return logger
}

func main() {
	// init Logger
	logInstance := initializeLogger()

	// init DB
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbConnString := os.Getenv("DATABASE_URL")

	if dbConnString == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	defer db.Close()

	// init Data repository for movie
	movieRepository, err := data.NewMovieRepository(db, logInstance)

	if err != nil {
		log.Fatalf("Failed to initialize movie repository: %v", err)
	}

	// init MovieHandler
	movieHandler := handlers.NewMovieHandler(movieRepository, logInstance)

	// API endpoints - movies
	http.HandleFunc("/api/v1/movies/top", movieHandler.GetTopMovies)
	http.HandleFunc("/api/v1/movies/random", movieHandler.GetRandomMovies)
	http.HandleFunc("/api/v1/movies/search", movieHandler.SearchMovies)
	http.HandleFunc("/api/v1/movies/genres", movieHandler.GetGenres)
	http.HandleFunc("/api/v1/movies/", movieHandler.GetMovie)

	// Serve static files
	http.Handle("/", http.FileServer(http.Dir("public")))
	logInstance.Info("Server started at http://localhost:8080")

	const addr = ":8080"

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		logInstance.Error("Failed to start server", err)
	}
}
