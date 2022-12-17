package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) betRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/event", s.events)

	return r
}

func (s *Server) events(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("events")

	evs, err := s.db.FetchEvents(ctx)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch events")
		respondErr(w, internalErr())

		return
	}

	var views []betEvent

	for _, e := range evs {
		views = append(views, betEventView(e))
	}

	respondJSON(w, http.StatusOK, views)
}
