package server

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/user"
	"github.com/swithek/sessionup"
)

type betUserHandler func(http.ResponseWriter, *http.Request, user.BetUser)
type userHandler func(http.ResponseWriter, *http.Request, user.User)
type adminHandler (func(http.ResponseWriter, *http.Request, user.AdminUser))

func (s *Server) withBetUser(h betUserHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := s.logger("withBetUser")

		session, ok := sessionup.FromContext(ctx)
		if !ok {
			respondErr(w, unauthorizedErr())
			return
		}

		userUUID, err := uuid.Parse(session.UserKey)
		if err != nil {
			log.Error().Err(err).Msg("cannot parse user key")
			respondErr(w, internalErr())

			return
		}

		u, ok, err := s.db.FetchBetUserByUUID(ctx, userUUID)
		if err != nil {
			log.Error().Err(err).Msg("cannot fetch bet user")
			respondErr(w, internalErr())

			return
		}

		if !ok {
			respondErr(w, notFoundErr())
			return
		}

		h(w, r, u)
	})
}

func (s *Server) withUser(h userHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := s.logger("withUser")

		session, ok := sessionup.FromContext(ctx)
		if !ok {
			log.Info().Msg("not found")
			respondErr(w, unauthorizedErr())
			return
		}

		userUUID, err := uuid.Parse(session.UserKey)
		if err != nil {
			log.Error().Err(err).Msg("cannot parse user key")
			respondErr(w, internalErr())

			return
		}

		u, ok, err := s.db.FetchUserByUUID(ctx, userUUID)
		if err != nil {
			log.Error().Err(err).Msg("cannot fetch user")
			respondErr(w, internalErr())

			return
		}

		if !ok {
			respondErr(w, notFoundErr())
			return
		}

		h(w, r, u)
	})
}

func (s *Server) authorizeAdmin(role user.Role, action string, hdl adminHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := s.logger("authorizeAdmin")

		session, ok := sessionup.FromContext(ctx)
		if !ok {
			log.Info().Msg("not found")
			respondErr(w, unauthorizedErr())
			return
		}

		userUUID, err := uuid.Parse(session.UserKey)
		if err != nil {
			log.Error().Err(err).Msg("cannot parse user key")
			respondErr(w, internalErr())

			return
		}

		adm, ok, err := s.db.FetchAdminUser(ctx, userUUID)
		if err != nil {
			log.Error().Err(err).Msg("cannot fetch admin")
			respondErr(w, internalErr())

			return
		}

		if !ok {
			respondErr(w, notFoundErr())
			return
		}

		if !adm.Permit(role) {
			respondErr(w, notFoundErr())
			return
		}

		lg := adm.Log(action)

		if err := s.db.InsertActionLog(ctx, lg); err != nil {
			log.Error().Err(err).Msg("cannot insert action log")
			respondErr(w, internalErr())

			return
		}

		hdl(w, r, adm)
	})
}
