package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zhangkui/go-feature-rollout/internal/service"
)

func TestCreateAndReadFlag(t *testing.T) {
	handler := NewHandler(service.New())
	body := bytes.NewBufferString(`{"key":"checkout-v2","description":"new checkout"}`)
	createRequest := httptest.NewRequest(http.MethodPost, "/flags", body)
	createResponse := httptest.NewRecorder()
	handler.ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", createResponse.Code, createResponse.Body.String())
	}

	readRequest := httptest.NewRequest(http.MethodGet, "/flags/checkout-v2", nil)
	readResponse := httptest.NewRecorder()
	handler.ServeHTTP(readResponse, readRequest)
	if readResponse.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", readResponse.Code, readResponse.Body.String())
	}
	var response map[string]any
	if err := json.Unmarshal(readResponse.Body.Bytes(), &response); err != nil {
		t.Fatal(err)
	}
	if response["key"] != "checkout-v2" {
		t.Fatalf("unexpected response: %v", response)
	}
}

func TestHealth(t *testing.T) {
	handler := NewHandler(service.New())
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", response.Code)
	}
}
