package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/golkity/Calc_2.0/internal/orchenstrator"
	"github.com/golkity/Calc_2.0/pkg/logger"
)

type Handler struct {
	log *logger.Logger
}

func NewHandler(log *logger.Logger) *Handler {
	return &Handler{log: log}
}

func (h *Handler) RegRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Expression == "" {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	id, err := orchestrator.AddNewExpression(req.Expression)
	if err != nil {
		h.log.WithError(err).Error("Failed to add expression")
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{"id": id}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.WithError(err).Error("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *Handler) ListExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expressions := orchestrator.GetAllExpressions()

	out := make([]map[string]interface{}, 0, len(expressions))
	for _, e := range expressions {
		out = append(out, map[string]interface{}{
			"id":     e.ID,
			"status": e.Status,
			"result": e.Result,
		})
	}

	resp := map[string]interface{}{"expressions": out}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.WithError(err).Error("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetExpressionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/expressions/")
	if path == "" {
		http.Error(w, "Expression ID is required", http.StatusBadRequest)
		return
	}
	exprID := path

	expression, err := orchestrator.GetExpressionByID(exprID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Expression not found: %v", err), http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"expression": map[string]interface{}{
			"id":     expression.ID,
			"status": expression.Status,
			"result": expression.Result,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.WithError(err).Error("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task := orchestrator.GetNextTask()
	if task == nil {
		http.Error(w, "No tasks available", http.StatusNotFound)
		return
	}

	resp := map[string]interface{}{
		"task": map[string]interface{}{
			"id":             task.ID,
			"arg1":           task.Arg1,
			"arg2":           task.Arg2,
			"operation":      task.Operation,
			"operation_time": task.OperationTime,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		h.log.WithError(err).Error("Error encoding response")
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *Handler) PostTaskResultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ID     int     `json:"id"`
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusUnprocessableEntity)
		return
	}

	if err := orchestrator.CompleteTask(req.ID, req.Result); err != nil {
		http.Error(w, "Failed to save result", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": "Result saved successfully",
	})
}
