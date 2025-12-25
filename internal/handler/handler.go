package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

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
	h.logger.Info("HTTP запрос", "method", r.Method, "path", r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.logger.Error("Ошибка парсинга формы", "error", err)
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	if title == "" {
		http.Error(w, "Название задачи обязательно", http.StatusBadRequest)
		return
	}

	id, err := h.database.CreateTask(title, description)
	if err != nil {
		h.logger.Error("Ошибка создания задачи", "error", err)
		http.Error(w, "Ошибка сохранения задачи", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"id":      id,
		"message": "Задача создана",
		"title":   title,
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Ошибка кодирования JSON", "error", err)
	}

	h.logger.Info("Задача создана", "id", id, "title", title)
}
