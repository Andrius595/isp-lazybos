package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/purse"
	"github.com/ramasauskas/ispbet/user"
	"github.com/shopspring/decimal"
)

type newDeposit struct {
	Amount   decimal.Decimal `json:"amount"`
	UserUUID uuid.UUID       `json:"user_uuid"`
}

type newWithdrawal struct {
	Amount   decimal.Decimal `json:"amount"`
	UserUUID uuid.UUID       `json:"user_uuid"`
}

type deposit struct {
	UUID      uuid.UUID       `json:"uuid"`
	Amount    decimal.Decimal `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
	UserUUID  uuid.UUID       `json:"user_uuid"`
}

type withdrawal struct {
	UUID      uuid.UUID       `json:"uuid"`
	Amount    decimal.Decimal `json:"amount"`
	Timestamp time.Time       `json:"timestamp"`
	UserUUID  uuid.UUID       `json:"user_uuid"`
}

func withdrawalView(w purse.Withdrawal) withdrawal {
	return withdrawal{
		UUID:      w.UUID,
		Amount:    w.Amount,
		Timestamp: w.Timestamp,
		UserUUID:  w.UserUUID,
	}
}

func depositView(d purse.Deposit) deposit {
	return deposit{
		UUID:      d.UUID,
		Amount:    d.Amount,
		Timestamp: d.Timestamp,
		UserUUID:  d.UserUUID,
	}
}

func (d *newDeposit) validate() error {
	if d.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount cannot be less than or equal to 0")
	}

	return nil
}

func (w *newWithdrawal) validate() error {
	if w.Amount.LessThanOrEqual(decimal.Zero) {
		return errors.New("amount cannot be less than or equal to 0")
	}

	return nil
}

func (s *Server) adminRouter() http.Handler {
	r := chi.NewRouter()

	r.Get("/bet-users", s.betUsers)
	r.Post("/identity-verifications", s.identityVerifications)
	r.Post("/finalize-identity-verification", s.finalizeIdentityVerification)
	r.Post("/deposit", s.createDeposit)
	r.Post("/withdraw", s.createWithdrawal)

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

	verViews := make([]identityVerification, len(verifications))

	for i, v := range verifications {
		verViews[i] = identityVerification{
			UUID:                v.UUID,
			UserUUID:            v.UserUUID,
			Status:              v.Status,
			IDPhotoBase64:       v.IDPhotoBase64,
			PortraitPhotoBase64: v.PortraitPhotoBase64,
			RespondedAt:         v.RespondedAt,
			CreatedAt:           v.CreatedAt,
		}
	}

	respondJSON(w, http.StatusOK, verViews)
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

func (s *Server) betUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("betUsers")

	uu, err := s.db.FetchBetUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch bet users")
		respondErr(w, internalErr())

		return
	}

	uuView := make([]betUser, len(uu))
	for i, u := range uu {
		uuView[i] = betUserView(u)
	}

	respondJSON(w, http.StatusOK, uuView)
}

func (s *Server) createDeposit(w http.ResponseWriter, r *http.Request) {
	var nd newDeposit

	if err := json.NewDecoder(r.Body).Decode(&nd); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ctx := r.Context()
	log := s.logger("createDeposit")

	if err := nd.validate(); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	d := purse.Deposit{
		UUID:      uuid.New(),
		Amount:    nd.Amount,
		Timestamp: time.Now(),
		UserUUID:  nd.UserUUID,
	}

	u, ok, err := s.db.FetchBetUserByUUID(ctx, d.UserUUID)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch bet user")
		respondErr(w, internalErr())

		return
	}

	if !ok {
		respondErr(w, notFoundErr())
		return
	}

	if err = u.Credit(d.Amount); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	if err = s.db.InsertDeposit(ctx, u, d); err != nil {
		log.Error().Err(err).Msg("cannot insert deposit")
		respondErr(w, internalErr())

		return
	}

	respondJSON(w, http.StatusCreated, depositView(d))
}

func (s *Server) createWithdrawal(w http.ResponseWriter, r *http.Request) {
	var nw newWithdrawal

	if err := json.NewDecoder(r.Body).Decode(&nw); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ctx := r.Context()
	log := s.logger("createWithdrawal")

	if err := nw.validate(); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	wd := purse.Withdrawal{
		UUID:      uuid.New(),
		Amount:    nw.Amount,
		Timestamp: time.Now(),
		UserUUID:  nw.UserUUID,
	}

	u, ok, err := s.db.FetchBetUserByUUID(ctx, wd.UserUUID)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch bet user")
		respondErr(w, internalErr())

		return
	}

	if !ok {
		respondErr(w, notFoundErr())
		return
	}

	if err = u.Debit(wd.Amount); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	if err = s.db.InsertWithdrawal(ctx, u, wd); err != nil {
		log.Error().Err(err).Msg("cannot insert withdrawal")
		respondErr(w, internalErr())

		return
	}

	respondJSON(w, http.StatusCreated, withdrawalView(wd))
}
