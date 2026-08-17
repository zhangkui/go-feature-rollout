package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/zhangkui/go-feature-rollout/internal/model"
	"github.com/zhangkui/go-feature-rollout/internal/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(flagService *service.Service) http.Handler {
	handler := &Handler{service: flagService}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.health)
	mux.HandleFunc("/flags", handler.flags)
	mux.HandleFunc("/flags/", handler.flagAction)
	return mux
}

func (h *Handler) health(writer http.ResponseWriter, _ *http.Request) {
	writeJSON(writer, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) flags(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var input struct {
		Key         string `json:"key"`
		Description string `json:"description"`
	}
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	flag, err := h.service.CreateFlag(input.Key, input.Description)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, flag)
}

func (h *Handler) flagAction(writer http.ResponseWriter, request *http.Request) {
	parts := splitPath(strings.TrimPrefix(request.URL.Path, "/flags/"))
	if len(parts) == 0 {
		writeError(writer, http.StatusNotFound, "not found")
		return
	}
	key := parts[0]
	if len(parts) == 1 {
		switch request.Method {
		case http.MethodGet:
			flag, err := h.service.GetFlag(key)
			if err != nil {
				writeServiceError(writer, err)
				return
			}
			writeJSON(writer, http.StatusOK, flag)
		case http.MethodDelete:
			if err := h.service.ArchiveFlag(key); err != nil {
				writeServiceError(writer, err)
				return
			}
			writer.WriteHeader(http.StatusNoContent)
		default:
			writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		}
		return
	}

	switch parts[1] {
	case "rules":
		h.rules(writer, request, key)
	case "publish":
		h.publish(writer, request, key)
	case "rollback":
		h.rollback(writer, request, key, parts)
	case "evaluate":
		h.evaluate(writer, request, key)
	case "batch-evaluate":
		h.batchEvaluate(writer, request, key)
	default:
		writeError(writer, http.StatusNotFound, "not found")
	}
}

func (h *Handler) rules(writer http.ResponseWriter, request *http.Request, key string) {
	if request.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var rule model.Rule
	if err := decodeJSON(request, &rule); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.UpsertDraftRule(key, rule); err != nil {
		writeServiceError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (h *Handler) publish(writer http.ResponseWriter, request *http.Request, key string) {
	if request.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	version, err := h.service.Publish(key, time.Now().UTC())
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusCreated, version)
}

func (h *Handler) rollback(writer http.ResponseWriter, request *http.Request, key string, parts []string) {
	if request.Method != http.MethodPost || len(parts) != 3 {
		writeError(writer, http.StatusBadRequest, "version is required")
		return
	}
	number, err := strconv.Atoi(parts[2])
	if err != nil {
		writeError(writer, http.StatusBadRequest, "invalid version")
		return
	}
	if err := h.service.Rollback(key, number); err != nil {
		writeServiceError(writer, err)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
}

func (h *Handler) evaluate(writer http.ResponseWriter, request *http.Request, key string) {
	if request.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	version, err := optionalVersion(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	var input model.EvaluationRequest
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	decision, err := h.service.Evaluate(key, version, input)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, decision)
}

func (h *Handler) batchEvaluate(writer http.ResponseWriter, request *http.Request, key string) {
	if request.Method != http.MethodPost {
		writeError(writer, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	version, err := optionalVersion(request)
	if err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	var input []model.EvaluationRequest
	if err := decodeJSON(request, &input); err != nil {
		writeError(writer, http.StatusBadRequest, err.Error())
		return
	}
	decisions, err := h.service.BatchEvaluate(request.Context(), key, version, input)
	if err != nil {
		writeServiceError(writer, err)
		return
	}
	writeJSON(writer, http.StatusOK, decisions)
}

func optionalVersion(request *http.Request) (int, error) {
	value := request.URL.Query().Get("version")
	if value == "" {
		return 0, nil
	}
	version, err := strconv.Atoi(value)
	if err != nil || version < 1 {
		return 0, errors.New("invalid version")
	}
	return version, nil
}

func decodeJSON(request *http.Request, target any) error {
	decoder := json.NewDecoder(io.LimitReader(request.Body, 1<<20))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func splitPath(path string) []string {
	segments := strings.Split(strings.Trim(path, "/"), "/")
	result := make([]string, 0, len(segments))
	for _, segment := range segments {
		if segment != "" {
			result = append(result, segment)
		}
	}
	return result
}

func writeServiceError(writer http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrFlagNotFound), errors.Is(err, service.ErrVersionNotFound):
		writeError(writer, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrFlagExists):
		writeError(writer, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrFlagArchived):
		writeError(writer, http.StatusGone, err.Error())
	default:
		writeError(writer, http.StatusBadRequest, err.Error())
	}
}

func writeError(writer http.ResponseWriter, status int, message string) {
	writeJSON(writer, status, map[string]string{"error": message})
}

func writeJSON(writer http.ResponseWriter, status int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	_ = json.NewEncoder(writer).Encode(value)
}
