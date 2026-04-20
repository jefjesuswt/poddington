package fleet

import (
	"log/slog"
	"net/http"
	"strings"

	podHTTP "github.com/jefjesuswt/poddington/shared/http"
)

type Handler struct {
	svc *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		svc: s,
	}
}

func (h *Handler) RegisterRoutes(router *podHTTP.Router) {

	router.NotFound(func(w http.ResponseWriter, req *http.Request) {
		slog.Warn("🚨 404 Interceptado", "method", req.Method, "path", req.URL.Path)
		podHTTP.ErrorJSON(w, http.StatusNotFound, "Endpoint not found in Poddington Hub.")
	})

	router.Group(func(r *podHTTP.Router) {
		r.Route("/api/nodes", func(n *podHTTP.Router) {
			n.Post("/{id}/ping", h.handlePing)
		})
	})
}

func (h *Handler) handlePing(w http.ResponseWriter, req *http.Request) {
	nodeID := req.PathValue("id")
	authHeader := req.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if token == "" || token == authHeader {
		podHTTP.ErrorJSON(w, http.StatusUnauthorized, "Missing or invalid authorization token.")
	}

	if err := h.svc.PingNode(req.Context(), nodeID, token); err != nil {
		if err.Error() == "node not found" {
			podHTTP.ErrorJSON(w, http.StatusNotFound, err.Error())
			return
		}

		podHTTP.ErrorJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	podHTTP.JSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "ping respected",
	})
}
