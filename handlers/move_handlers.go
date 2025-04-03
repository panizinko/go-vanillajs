package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"adrianzinko.com/gomovies/data"
	"adrianzinko.com/gomovies/logger"
	"adrianzinko.com/gomovies/models"
)

type MovieHandler struct {
	repository data.MovieRepository
	logger     *logger.Logger
}

func NewMovieHandler(repository data.MovieRepository, log *logger.Logger) *MovieHandler {
	return &MovieHandler{
		repository: repository,
		logger:     log,
	}
}

func (h *MovieHandler) writeJSONResponse(w http.ResponseWriter, data any) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode response", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *MovieHandler) handleRepositoryError(w http.ResponseWriter, err error, context string) bool {
	if err != nil {
		if err == data.ErrMovieNotFound {
			http.Error(w, context, http.StatusNotFound)
			return true
		}
		h.logger.Error(context, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return true
	}
	return false
}

func (h *MovieHandler) parseID(w http.ResponseWriter, idStr string) (int, bool) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.logger.Error("Invalid ID format", err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func (h *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.repository.GetTopMovies()

	if h.handleRepositoryError(w, err, "Failed to get movies") {
		return
	}
	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served top movies")
	}
}

func (h *MovieHandler) GetRandomMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.repository.GetRandomMovies()
	if h.handleRepositoryError(w, err, "Failed to get movies") {
		return
	}
	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served random movies")
	}
}

func (h *MovieHandler) SearchMovies(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	order := r.URL.Query().Get("order")
	genreStr := r.URL.Query().Get("genre")

	var genre *int
	if genreStr != "" {
		genreInt, ok := h.parseID(w, genreStr)
		if !ok {
			return
		}
		genre = &genreInt
	}

	var movies []models.Movie
	var err error
	if query != "" {
		movies, err = h.repository.SearchMoviesByName(query, order, genre)
	}
	if h.handleRepositoryError(w, err, "Failed to get movies") {
		return
	}
	if h.writeJSONResponse(w, movies) == nil {
		h.logger.Info("Successfully served movies")
	}
}

func (h *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/api/v1/movies/"):]
	id, ok := h.parseID(w, idStr)

	if !ok {
		return
	}

	movie, err := h.repository.GetMovieByID(id)
	if h.handleRepositoryError(w, err, "Failed to get movie by ID") {
		return
	}
	if h.writeJSONResponse(w, movie) == nil {
		h.logger.Info("Successfully served movie with ID: " + idStr)
	}
}

func (h *MovieHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := h.repository.GetAllGenres()
	if h.handleRepositoryError(w, err, "Failed to get genres") {
		return
	}
	if h.writeJSONResponse(w, genres) == nil {
		h.logger.Info("Successfully served genres")
	}
}
