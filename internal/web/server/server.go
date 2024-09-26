package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/andersonribeir0/mouchat/internal/web"
	"github.com/andersonribeir0/mouchat/internal/web/handlers"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/nedpals/supabase-go"
)

type Server struct {
	db   *supabase.Client
	port string
}

func NewServer(listenPort string, db *supabase.Client) *http.Server {
	NewServer := &Server{
		port: listenPort,

		db: db,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Handle("/*", http.StripPrefix("/", http.FileServer(http.FS(web.Files))))

	r.Get("/", MakeHandler(handlers.HandleHomeIndex))

	return r
}

func MakeHandler(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("internal server error", "err", err, "path", r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func render(r *http.Request, w http.ResponseWriter, component templ.Component) error {
	return component.Render(r.Context(), w)
}
