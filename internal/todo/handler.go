package todo

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Handler wires the repository to HTTP routes.
type Handler struct {
	repo *Repository
}

// NewHandler creates a new Handler instance.
func NewHandler(repo *Repository) *Handler {
	return &Handler{repo: repo}
}

// RegisterRoutes attaches todo and health routes to the provided router.
func (h *Handler) RegisterRoutes(r chi.Router) {
	r.Get("/healthz", h.health)

	r.Route("/todos", func(r chi.Router) {
		r.Get("/", h.list)
		r.Post("/", h.create)
		r.Get("/{id}", h.get)
		r.Put("/{id}", h.update)
		r.Delete("/{id}", h.delete)
	})
}

func (h *Handler) health(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repo.List(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not list todos")
		return
	}
	respondJSON(w, http.StatusOK, todos)
}

func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	todo, err := h.repo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			respondError(w, http.StatusNotFound, "todo not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "could not fetch todo")
		return
	}

	respondJSON(w, http.StatusOK, todo)
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}

	if err := decodeJSON(r, &payload); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	title := strings.TrimSpace(payload.Title)
	if title == "" {
		respondError(w, http.StatusBadRequest, "title is required")
		return
	}

	todo, err := h.repo.Create(r.Context(), title, payload.Completed)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not create todo")
		return
	}

	respondJSON(w, http.StatusCreated, todo)
}

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var payload struct {
		Title     *string `json:"title"`
		Completed *bool   `json:"completed"`
	}

	if err := decodeJSON(r, &payload); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if payload.Title == nil && payload.Completed == nil {
		respondError(w, http.StatusBadRequest, "no fields to update")
		return
	}

	if payload.Title != nil {
		trimmed := strings.TrimSpace(*payload.Title)
		if trimmed == "" {
			respondError(w, http.StatusBadRequest, "title cannot be empty")
			return
		}
		payload.Title = &trimmed
	}

	todo, err := h.repo.Update(r.Context(), id, payload.Title, payload.Completed)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			respondError(w, http.StatusNotFound, "todo not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "could not update todo")
		return
	}

	respondJSON(w, http.StatusOK, todo)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		if errors.Is(err, ErrNotFound) {
			respondError(w, http.StatusNotFound, "todo not found")
			return
		}
		respondError(w, http.StatusInternalServerError, "could not delete todo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.ParseInt(idStr, 10, 64)
}

func decodeJSON(r *http.Request, dest any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dest); err != nil {
		return err
	}
	return nil
}

func respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
