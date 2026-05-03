package fleet

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jefjesuswt/walroos/shared/whttp"
)

type NodeManager interface {
	RegisterNode(ctx context.Context, name, address string) (Node, error)
	ListNodes(ctx context.Context) ([]Node, error)
	RemoveNode(ctx context.Context, id string) error
	PingNode(ctx context.Context, id, token string) error
}

type Handler struct {
	nm NodeManager
}

func NewHandler(nm NodeManager) *Handler {
	return &Handler{
		nm: nm,
	}
}

func (h *Handler) RegisterRoutes(router *whttp.Router) {

	router.NotFound(func(w http.ResponseWriter, req *http.Request) {
		whttp.ErrorJSON(w, http.StatusNotFound, "Endpoint not found in Walroos Hub.")
	})

	router.Group(func(r *whttp.Router) {
		r.Route("/api/nodes", func(n *whttp.Router) {
			n.Post("/", h.handleNewNode)
			n.Post("/{id}/ping", h.handlePing)
			n.Get("/", h.handleList)
			n.Delete("/{id}", h.handleRemove)
		})
	})
}

func (h *Handler) handleNewNode(w http.ResponseWriter, req *http.Request) {
	type requestPayload struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	}

	var rp requestPayload

	err := json.NewDecoder(req.Body).Decode(&rp)
	if err != nil {
		whttp.ErrorJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	node, err := h.nm.RegisterNode(req.Context(), rp.Name, rp.Address)
	if err != nil {

		if errors.Is(err, ErrNodeNameAlreadyExists) {
			whttp.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, ErrNodeAddressAlreadyExists) {
			whttp.ErrorJSON(w, http.StatusBadRequest, err.Error())
			return
		}

		whttp.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	whttp.JSON(w, http.StatusCreated, node)
}

func (h *Handler) handleList(w http.ResponseWriter, req *http.Request) {
	nodes, err := h.nm.ListNodes(req.Context())
	if err != nil {
		whttp.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	whttp.JSON(w, http.StatusOK, nodes)
}

func (h *Handler) handleRemove(w http.ResponseWriter, req *http.Request) {
	nodeID := req.PathValue("id")

	if err := h.nm.RemoveNode(req.Context(), nodeID); err != nil {
		if errors.Is(err, ErrNodeNotFound) {
			whttp.ErrorJSON(w, http.StatusNotFound, err.Error())
			return
		}

		whttp.ErrorJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	whttp.JSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "node removed",
	})

}

func (h *Handler) handlePing(w http.ResponseWriter, req *http.Request) {
	nodeID := req.PathValue("id")

	token, err := getBearerToken(req)
	if err != nil {
		whttp.ErrorJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	if err := h.nm.PingNode(req.Context(), nodeID, token); err != nil {
		if errors.Is(err, ErrNodeNotFound) {
			whttp.ErrorJSON(w, http.StatusNotFound, err.Error())
			return
		}

		whttp.ErrorJSON(w, http.StatusUnauthorized, err.Error())
		return
	}

	whttp.JSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "ping respected",
	})
}

func getBearerToken(req *http.Request) (string, error) {
	authHeader := req.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	if token == "" || token == authHeader {
		return "", errors.New("Missing or invalid authorization token.")
	}

	return token, nil
}
