package fleet

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jefjesuswt/walroos/shared/whttp"
)

type mockNodeManager struct {
	RegisterNodeFn func(ctx context.Context, name, address string) (Node, error)
	ListNodesFn    func(ctx context.Context) ([]Node, error)
	RemoveNodeFn   func(ctx context.Context, id string) error
	PingNodeFn     func(ctx context.Context, id, token string) error
}

func (m *mockNodeManager) RegisterNode(ctx context.Context, name, address string) (Node, error) {
	if m.RegisterNodeFn != nil {
		return m.RegisterNodeFn(ctx, name, address)
	}
	panic("not implemented")
}

func (m *mockNodeManager) ListNodes(ctx context.Context) ([]Node, error) {
	if m.ListNodesFn != nil {
		return m.ListNodesFn(ctx)
	}
	panic("not implemented")
}

func (m *mockNodeManager) RemoveNode(ctx context.Context, id string) error {
	if m.RemoveNodeFn != nil {
		return m.RemoveNodeFn(ctx, id)
	}
	panic("not implemented")
}

func (m *mockNodeManager) PingNode(ctx context.Context, id, token string) error {
	if m.PingNodeFn != nil {
		return m.PingNodeFn(ctx, id, token)
	}
	panic("not implemented")
}

func TestHandlers_RegisterRoutes(t *testing.T) {
	h := NewHandler(&mockNodeManager{})
	router := whttp.NewRouter()
	h.RegisterRoutes(router)

	t.Run("Should return 404.", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/wrong/route", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}

func TestHandlers_handlerNewNode(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(*mockNodeManager)
		expectedStatus int
	}{
		{
			name:        "Node successfully created.",
			requestBody: `{"name": "worker-1", "address": "10.0.0.1"}`,
			mockSetup: func(mnm *mockNodeManager) {
				mnm.RegisterNodeFn = func(ctx context.Context, name, address string) (Node, error) {
					return Node{
						ID:      "node-1",
						Name:    name,
						Address: address,
					}, nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:        "Node name already exists.",
			requestBody: `{"name": "workerk-1", "address": "10.0.0.1}`,
			mockSetup: func(mnm *mockNodeManager) {
				mnm.RegisterNodeFn = func(ctx context.Context, name, address string) (Node, error) {
					return Node{}, ErrNodeNameAlreadyExists
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:        "Node address already exists.",
			requestBody: `{"name": "workerk-1", "address": "10.0.0.1}`,
			mockSetup: func(mnm *mockNodeManager) {
				mnm.RegisterNodeFn = func(ctx context.Context, name, address string) (Node, error) {
					return Node{}, ErrNodeAddressAlreadyExists
				}
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid request body.",
			requestBody:    `{"name": "worker-1", "address":`,
			mockSetup:      func(mnm *mockNodeManager) {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockNodeManager{}

			if tt.mockSetup != nil {
				tt.mockSetup(&mock)
			}

			h := NewHandler(&mock)

			router := whttp.NewRouter()

			router.Post("/api/nodes", h.handleNewNode)

			req := httptest.NewRequest(http.MethodPost, "/api/nodes", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestHandlers_handleList(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mockNodeManager)
		expectedStatus int
	}{
		{

			name: "200. Showing list.",
			mockSetup: func(mnm *mockNodeManager) {
				mnm.ListNodesFn = func(ctx context.Context) ([]Node, error) {
					nodes := []Node{
						{
							ID:      "node-1",
							Name:    "worker-1",
							Address: "10.0.0.1",
						},
						{
							ID:      "node-2",
							Name:    "worker-2",
							Address: "10.0.0.2",
						},
					}
					return nodes, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "500. Internal Server Error",
			mockSetup: func(mnm *mockNodeManager) {
				mnm.ListNodesFn = func(ctx context.Context) ([]Node, error) {
					return nil, errors.New("error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := mockNodeManager{}
			if tt.mockSetup != nil {
				tt.mockSetup(&mock)
			}

			h := NewHandler(&mock)

			router := whttp.NewRouter()

			router.Get("/api/nodes", h.handleList)

			req := httptest.NewRequest(http.MethodGet, "/api/nodes", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestHandlers_handleRemove(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
		mockSetup      func(*mockNodeManager)
	}{
		{

			name:           "Failed to remove, ErrNodeNotFound",
			expectedStatus: http.StatusNotFound,
			mockSetup: func(mnm *mockNodeManager) {
				mnm.RemoveNodeFn = func(ctx context.Context, nodeID string) error {
					return ErrNodeNotFound
				}
			},
		},
		{
			name:           "Failed to remove, 500.",
			expectedStatus: http.StatusInternalServerError,
			mockSetup: func(mnm *mockNodeManager) {
				mnm.RemoveNodeFn = func(ctx context.Context, nodeID string) error {
					return errors.New("error")
				}
			},
		},
		{
			name:           "Successfully removed.",
			expectedStatus: http.StatusOK,
			mockSetup: func(mnm *mockNodeManager) {
				mnm.RemoveNodeFn = func(ctx context.Context, nodeID string) error {
					return nil
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := mockNodeManager{}
			if tt.mockSetup != nil {
				tt.mockSetup(&mock)
			}

			h := NewHandler(&mock)

			router := whttp.NewRouter()

			router.Delete("/api/nodes/{id}", h.handleRemove)

			req := httptest.NewRequest(http.MethodDelete, "/api/nodes/node-1", nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestHandlers_handlePing(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mockNodeManager)
		expectedStatus int
		headers        map[string]any
	}{
		{
			name: "Valid token, updating...",
			mockSetup: func(mnm *mockNodeManager) {
				mnm.PingNodeFn = func(ctx context.Context, id, token string) error {
					return nil
				}
			},
			expectedStatus: http.StatusOK,
			headers: map[string]any{
				"Authorization": "Bearer valid-token",
			},
		},
		{
			name: "Invalid token. Unauthorized",
			mockSetup: func(mnm *mockNodeManager) {
				mnm.PingNodeFn = func(ctx context.Context, id, token string) error {
					return errors.New("error")
				}
			},
			expectedStatus: http.StatusUnauthorized,
			headers: map[string]any{
				"Authorization": "Bearer invalid-token",
			},
		},
		{
			name:           "No token in headers.",
			expectedStatus: http.StatusUnauthorized,
			headers: map[string]any{
				"Authorization": nil,
			},
		},
		{
			name: "Node not found.",
			mockSetup: func(mnm *mockNodeManager) {
				mnm.PingNodeFn = func(ctx context.Context, id, token string) error {
					return ErrNodeNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
			headers: map[string]any{
				"Authorization": "Bearer valid-token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock := mockNodeManager{}
			if tt.mockSetup != nil {
				tt.mockSetup(&mock)
			}

			h := NewHandler(&mock)

			router := whttp.NewRouter()

			router.Post("/api/nodes/{id}/ping", h.handlePing)

			req := httptest.NewRequest(http.MethodPost, "/api/nodes/node-1/ping", nil)
			rec := httptest.NewRecorder()

			for k, v := range tt.headers {
				req.Header.Set(k, fmt.Sprintf("%v", v))
			}

			router.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}

}
