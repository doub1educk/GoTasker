package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	db "github.com/doub1educk/gotasker/internal/database"
)

type TaskHandler struct {
	database *db.Database
	logger   *slog.Logger
}

func NewTaskHandler(database *db.Database, logger *slog.Logger) *TaskHandler {
	return &TaskHandler{
		database: database,
		logger:   logger,
	}
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("HTTP request", "method", r.Method, "path", r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "method not supported", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := h.database.GetAllTasks()
	if err != nil {
		h.logger.Error("error with receiving all tasks", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		h.logger.Error("error encoding", err)
		return
	}
	h.logger.Info("tasks_count", len(tasks))
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("http ", "method", r.Method, "path", r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.logger.Error("error parse form", "error", err)
		http.Error(w, "data is not correct", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	if title == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	id, err := h.database.CreateTask(title, description)
	if err != nil {
		h.logger.Error("error make task", "error", err)
		http.Error(w, "error save tak", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":      id,
		"message": "task is create",
		"title":   title,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("err with encode JSON", "error", err)
	}

	h.logger.Info("task create", "id", id, "title", title)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id not allowed", http.StatusBadRequest)
		return
	}
	if err := h.database.DeleteTask(id); err != nil {
		h.logger.Error("failed to delete task", "error", err)
		http.Error(w, "error", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "task was deleted", id)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		h.logger.Error("failed to parse form", "error", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task id", http.StatusBadRequest)
		return
	}
	status := r.FormValue("status")

	if status != "pending" && status != "done" {
		http.Error(w, "Use valid status", http.StatusBadRequest)
		return
	}
	if err := h.database.UpdateTask(id, status); err != nil {
		h.logger.Error("failed to update task:", id, "error", err)
		http.Error(w, "server err", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Task %d updated", id)
}
