package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/kam1k88/gokeenapi/pkg/goarapi"
)

// Server exposes REST API backed by AnyRouterAPI facade.
type Server struct {
	api *goarapi.AnyRouterAPI
	mux *http.ServeMux
}

// New constructs HTTP server for provided API facade.
func New(api *goarapi.AnyRouterAPI) *Server {
	s := &Server{api: api, mux: http.NewServeMux()}
	s.routes()
	return s
}

// Handler returns configured http.Handler.
func (s *Server) Handler() http.Handler {
	return s.mux
}

// ListenAndServe starts HTTP server on given address.
func (s *Server) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/routers", s.handleRouters)
	s.mux.HandleFunc("/api/routers/", s.handleRouterActions)
}

func (s *Server) handleRouters(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		names := s.api.Routers()
		routers := make([]goarapi.DeviceInfo, 0, len(names))
		for _, name := range names {
			info, err := s.api.DeviceInfo(r.Context(), name)
			if err != nil {
				writeError(w, http.StatusBadGateway, err)
				return
			}
			routers = append(routers, info)
		}
		writeJSON(w, http.StatusOK, map[string]any{"routers": routers})
	default:
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method))
	}
}

func (s *Server) handleRouterActions(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.TrimPrefix(r.URL.Path, "/api/routers/")
	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 || segments[0] == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("router name is required"))
		return
	}
	name := segments[0]
	if len(segments) == 1 {
		writeError(w, http.StatusNotFound, fmt.Errorf("unknown router action"))
		return
	}

	switch segments[1] {
	case "routes":
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method))
			return
		}
		routes, err := s.api.ListRoutes(r.Context(), name)
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"routes": routes})
	case "route":
		switch r.Method {
		case http.MethodPost:
			var route goarapi.Route
			if err := json.NewDecoder(r.Body).Decode(&route); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			if err := s.api.AddRoute(r.Context(), name, route); err != nil {
				writeError(w, http.StatusBadGateway, err)
				return
			}
			writeJSON(w, http.StatusCreated, map[string]any{"status": "created"})
		case http.MethodDelete:
			if len(segments) < 3 {
				writeError(w, http.StatusBadRequest, fmt.Errorf("route identifier is required"))
				return
			}
			key, err := url.PathUnescape(segments[2])
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			if err := s.api.DeleteRoute(r.Context(), name, key); err != nil {
				writeError(w, http.StatusBadGateway, err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"status": "deleted"})
		default:
			writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method %s not allowed", r.Method))
		}
	default:
		writeError(w, http.StatusNotFound, fmt.Errorf("unknown router action"))
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}
