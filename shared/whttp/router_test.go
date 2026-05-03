package whttp

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRouter_NotFound(t *testing.T) {
	r := NewRouter()

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		slog.Warn("404 Intercepted!", "method", req.Method, "path", req.URL.Path)
		ErrorJSON(w, http.StatusNotFound, "404 Not Found.")
	})

	t.Run("Should use cusom NotFound", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/pang", nil)

		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}

		actualBody := strings.TrimSpace(rec.Body.String())
		expectedBody := `{"error":"404 Not Found."}`

		if actualBody != expectedBody {
			t.Errorf("Expected body: %s, got %s", expectedBody, actualBody)
		}
	})

}

func TestRouter_Methods(t *testing.T) {

	r := NewRouter()

	testHandler := func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(req.Method + " OK"))
	}

	r.Get("/test", testHandler)
	r.Post("/test", testHandler)
	r.Put("/test", testHandler)
	r.Patch("/test", testHandler)
	r.Delete("/test", testHandler)
	r.Connect("/test", testHandler)
	r.Head("/test", testHandler)
	r.Options("/test", testHandler)
	r.Trace("/test", testHandler)

	tests := []struct {
		name         string
		method       string
		expectedBody string
	}{
		{"Should return GET", http.MethodGet, "GET OK"},
		{"Should return POST", http.MethodPost, "POST OK"},
		{"Should return PATCH", http.MethodPatch, "PATCH OK"},
		{"Should return PUT", http.MethodPut, "PUT OK"},
		{"Should return DELETE", http.MethodDelete, "DELETE OK"},
		{"Should return CONNECT", http.MethodConnect, "CONNECT OK"},
		{"Should return HEAD", http.MethodHead, "HEAD OK"},
		{"Should return OPTIONS", http.MethodOptions, "OPTIONS OK"},
		{"Should return TRACE", http.MethodTrace, "TRACE OK"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(tt.method, "/test", nil)

			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
			}

			actualBody := strings.TrimSpace(rec.Body.String())

			if actualBody != tt.expectedBody {
				t.Errorf("Expected body: %s, got %s", tt.expectedBody, actualBody)
			}
		})
	}
}

func TestRouter_Use(t *testing.T) {

	r := NewRouter()

	spyMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Spy-Middleware", "been here")
			next.ServeHTTP(w, req)
		})
	}

	r.Use(spyMiddleware)

	r.Get("/test", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
	}

	spyHeader := rec.Header().Get("X-Spy-Middleware")

	if spyHeader != "been here" {
		t.Errorf("Expected header 'been here', got %s", spyHeader)
	}
}

func TestRouter_Route(t *testing.T) {

	r := NewRouter()

	r.Route("/api", func(api *Router) {
		api.Get("/status", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	t.Run("Should return 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/status", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}

func TestRouter_Mount(t *testing.T) {
	r := NewRouter()

	sr := NewRouter()
	sr.Get("/status", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Mount("/api", sr)

	t.Run("Should return 200", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}
	})

	t.Run("Should return 404", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/status", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Errorf("Expected status code %d, got %d", http.StatusNotFound, rec.Code)
		}
	})
}

func TestRouter_Group(t *testing.T) {
	r := NewRouter()

	spyMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("X-Spy-Middleware", "been here")
			next.ServeHTTP(w, req)
		})
	}

	r.Get("/not-spied", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Group(func(api *Router) {
		api.Use(spyMiddleware)
		api.Get("/spied", func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
	})

	t.Run("Should not been spied on", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/not-spied", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		spyHeader := rec.Header().Get("X-Spy-Middleware")

		if spyHeader != "" {
			t.Errorf("Expected header '', got %s", spyHeader)
		}
	})

	t.Run("Should been spied on", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/spied", nil)
		rec := httptest.NewRecorder()

		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		spyHeader := rec.Header().Get("X-Spy-Middleware")

		if spyHeader != "been here" {
			t.Errorf("Expected header 'been here', got %s", spyHeader)
		}
	})
}
