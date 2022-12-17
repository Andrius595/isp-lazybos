package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"
	"github.com/swithek/sessionup"
)

type serverErr struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

type Server struct {
	srv      *http.Server
	db       DB
	log      zerolog.Logger
	sessions *sessionup.Manager
	email    EmailSender

	wg sync.WaitGroup
}

func NewServer(port uint, sessionStore sessionup.Store, email EmailSender, db DB, log zerolog.Logger) *Server {
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	sessions := sessionup.NewManager(sessionStore)

	return &Server{
		srv:      srv,
		log:      log,
		db:       db,
		sessions: sessions,
		email:    email,
	}
}

func (s *Server) Run() error {
	s.log.Info().Msg("starting HTTP server")
	defer s.log.Info().Msg("HTTP server stopped")

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
	}))

	r.Mount("/bet-user", s.betUserRouter())
	r.Mount("/user", s.userRouter())
	r.Mount("/admin", s.adminRouter())
	r.Mount("/betting", s.betRouter())

	s.srv.Handler = r

	err := s.srv.ListenAndServe()
	switch err {
	case nil, http.ErrServerClosed:
		return nil
	default:
		return err
	}
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	s.log.Info().Msg("stopping HTTP server")
	if err := s.srv.Shutdown(ctx); err != nil {
		return err
	}

	s.wg.Wait()

	return nil
}

func (s *Server) logger(fn string) zerolog.Logger {
	return s.log.With().Str("function", fn).Logger()
}

func respondJSON(w http.ResponseWriter, code int, data any) {
	w.Header().Add("Content-type", "application/json")

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

func respondOK(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func respondErr(w http.ResponseWriter, se serverErr) {
	w.Header().Add("Content-type", "application/json")

	w.WriteHeader(se.Code)
	json.NewEncoder(w).Encode(se)
}

func internalErr() serverErr {
	return serverErr{
		Code:    http.StatusInternalServerError,
		Message: "internal server error",
	}
}

func notFoundErr() serverErr {
	return serverErr{
		Code:    http.StatusNotFound,
		Message: "not found",
	}
}

func badRequestErr(err error) serverErr {
	return serverErr{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}

func unauthorizedErr() serverErr {
	return serverErr{
		Code:    http.StatusUnauthorized,
		Message: "unauthorized",
	}
}
