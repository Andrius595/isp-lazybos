package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/user"
)

func (s *Server) adminRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/identity-verifications", s.identityVerifications)
	r.Post("/finalize-identity-verification", s.finalizeIdentityVerification)

	return r
}

func (s *Server) identityVerifications(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("identityVerifications")

	verifications, err := s.db.FetchIdentityVerifications(ctx)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch verifications")
		respondErr(w, internalErr())

		return
	}

	respondJSON(w, http.StatusOK, verifications)
}

func (s *Server) finalizeIdentityVerification(w http.ResponseWriter, r *http.Request) {
	var input struct {
		VerificationUUID uuid.UUID `json:"verification_uuid"`
		Accept           bool      `json:"accept"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ctx := r.Context()
	log := s.logger("finalizeIdentityVerification")

	ver, ok, err := s.db.FetchIdentityVerification(ctx, input.VerificationUUID)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch identity verification")
		respondErr(w, internalErr())

		return
	}

	if !ok {
		respondErr(w, notFoundErr())
		return
	}

	bu, ok, err := s.db.FetchBetUserByUUID(ctx, ver.UserUUID)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch user")
		respondErr(w, internalErr())

		return
	}

	if !ok {
		respondErr(w, notFoundErr())
		return
	}

	if input.Accept {
		if err = user.VerifyBetUserIdentity(&bu, &ver); err != nil {
			respondErr(w, badRequestErr(err))
			return
		}
	} else {
		if err = ver.Reject(); err != nil {
			respondErr(w, badRequestErr(err))
			return
		}
	}

	if err = s.db.InsertIdentityVerificationUpdate(ctx, bu, ver); err != nil {
		log.Error().Err(err).Msg("cannot insert verification update")
		respondErr(w, internalErr())

		return
	}

	respondOK(w)
}
