package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/user"
	"github.com/shopspring/decimal"
)

type newUserBet struct {
	SelectionUUID uuid.UUID       `json:"selection_uuid"`
	Stake         decimal.Decimal `json:"stake"`
	Winner        string          `json:"winner"`
}

type userBet struct {
	UUID          uuid.UUID       `json:"uuid"`
	Stake         decimal.Decimal `json:"stake"`
	SelectionUUID uuid.UUID       `json:"selection_uuid"`
}

func userBetView(b bet.Bet) userBet {
	return userBet{
		UUID:          b.UUID,
		Stake:         b.Stake,
		SelectionUUID: b.SelectionUUID,
	}
}

type newBetUser struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (bu newBetUser) Materialize() (user.BetUser, error) {
	if bu.Email == "" {
		return user.BetUser{}, errors.New("email not provided")
	}

	if bu.FirstName == "" {
		return user.BetUser{}, errors.New("first name not provided")
	}

	if bu.LastName == "" {
		return user.BetUser{}, errors.New("last name not provided")
	}

	if bu.Password == "" {
		return user.BetUser{}, errors.New("password not provided")
	}

	u := user.BetUser{
		User: user.User{
			UUID:      uuid.New(),
			Email:     bu.Email,
			FirstName: bu.FirstName,
			LastName:  bu.LastName,
		},
	}

	if err := u.SetPassword(bu.Password); err != nil {
		return user.BetUser{}, err
	}

	return u, nil
}

type betUser struct {
	UUID               uuid.UUID `json:"uuid"`
	Email              string    `json:"email"`
	FirstName          string    `json:"first_name"`
	LastName           string    `json:"last_name"`
	EmailVerified      bool      `json:"email_verified"`
	IdentitityVerified bool      `json:"identitity_verified"`
}

func betUserView(u user.BetUser) betUser {
	return betUser{
		UUID:               u.UUID,
		Email:              u.Email,
		FirstName:          u.FirstName,
		LastName:           u.LastName,
		EmailVerified:      u.EmailVerified,
		IdentitityVerified: u.IdentityVerified,
	}
}

type newIdentityVerification struct {
	IDPhoto       string `json:"id_photo"`
	PortraitPhoto string `json:"portrait_photo"`
}

func (id newIdentityVerification) Validate() error {
	if len(id.IDPhoto) == 0 {
		return errors.New("no ID photo provided")
	}

	if len(id.PortraitPhoto) == 0 {
		return errors.New("no portrait photo provided")
	}

	return nil
}

type identityVerification struct {
	UUID                uuid.UUID                       `json:"uuid"`
	User                betUser                         `json:"user"`
	Status              user.IdentityVerificationStatus `json:"status"`
	IDPhotoBase64       string                          `json:"id_photo_base_64"`
	PortraitPhotoBase64 string                          `json:"portrait_photo_base_64"`
	RespondedAt         time.Time                       `json:"responded_at"`
	CreatedAt           time.Time                       `json:"created_at"`
}

func identityVerificationView(id user.IdentityVerification, bu betUser) identityVerification {
	return identityVerification{
		UUID:                id.UUID,
		User:                bu,
		Status:              id.Status,
		IDPhotoBase64:       id.IDPhotoBase64,
		PortraitPhotoBase64: id.PortraitPhotoBase64,
		RespondedAt:         id.RespondedAt,
		CreatedAt:           id.CreatedAt,
	}
}

func (s *Server) betUserRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/register", s.registerBetUser)
	r.Post("/login", s.loginBetUser)
	r.Post("/logout", s.logout)

	r.Group(func(r chi.Router) {
		r.Use(s.sessions.Auth)

		r.Get("/me", s.withBetUser(s.betUserMe))
		r.Post("/identity-verification", s.withBetUser(s.createVerificationRequest))
		r.Post("/bet", s.withBetUser(s.bet))
	})

	return r
}

func (s *Server) betUserMe(w http.ResponseWriter, r *http.Request, bu user.BetUser) {
	respondJSON(w, http.StatusOK, betUserView(bu))
}

func (s *Server) registerBetUser(w http.ResponseWriter, r *http.Request) {
	var newUser newBetUser

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	u, err := newUser.Materialize()
	if err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ctx := r.Context()
	log := s.logger("registerBetUser")

	_, ok, err := s.db.FetchBetUserByEmail(ctx, u.Email)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch bet user")
		respondErr(w, internalErr())

		return
	}

	if ok {
		respondErr(w, badRequestErr(errors.New("user already exists with provided email")))
		return
	}

	if err = s.db.InsertBetUser(ctx, u); err != nil {
		log.Error().Err(err).Msg("cannot insert bet user")
		respondErr(w, internalErr())

		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		log = s.logger("registerBetUser:defer")

		if err := s.sendEmailVerification(ctx, u.User); err != nil {
			log.Err(err).Msg("cannot send email verification")
		}
	}()

	respondJSON(w, http.StatusCreated, betUserView(u))
}

func (s *Server) loginBetUser(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ctx := r.Context()
	log := s.logger("loginBetUser")

	u, ok, err := s.db.FetchBetUserByEmail(ctx, input.Email)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch user")
		respondErr(w, internalErr())

		return
	}

	if !ok {
		respondErr(w, notFoundErr())
		return
	}

	if !u.Login(input.Password) {
		respondErr(w, badRequestErr(errors.New("invalid password")))
		return
	}

	if err = s.sessions.Init(w, r, u.UUID.String()); err != nil {
		log.Error().Err(err).Msg("cannot initialize session")
		respondErr(w, internalErr())

		return
	}

	respondJSON(w, http.StatusOK, betUserView(u))
}

func (s *Server) createVerificationRequest(w http.ResponseWriter, r *http.Request, bu user.BetUser) {
	var req newIdentityVerification

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	if err := req.Validate(); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ver, err := bu.CreateVerificationRequest(req.IDPhoto, req.PortraitPhoto)
	if err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ctx := r.Context()
	log := s.logger("createVerificationRequest")

	if err = s.db.InsertBetUserIdentityVerification(ctx, ver); err != nil {
		log.Error().Err(err).Msg("cannot insert identity verification")
		respondErr(w, internalErr())

		return
	}

	respondJSON(w, http.StatusCreated, identityVerificationView(ver, betUserView(bu)))
}

func (s *Server) bet(w http.ResponseWriter, r *http.Request, u user.BetUser) {
	var nb newUserBet

	if err := json.NewDecoder(r.Body).Decode(&nb); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	b := bet.Bet{
		UUID:            uuid.New(),
		UserUUID:        u.UUID,
		SelectionUUID:   nb.SelectionUUID,
		SelectionWinner: bet.Winner(nb.Winner),
		Stake:           nb.Stake,
		State:           bet.BetStateTBD,
	}

	ctx := r.Context()
	log := s.logger("bet")

	resp, err := s.better.Bet(ctx, &b, &u)
	if err != nil {
		log.Error().Err(err).Msg("cannot place bet")
		respondErr(w, internalErr())

		return
	}

	if !resp.Ok {
		respondErr(w, badRequestErr(errors.New(resp.ErrorMessage)))
		return
	}

	respondJSON(w, http.StatusCreated, userBetView(b))
}

type BetResponse struct {
	Ok           bool
	ErrorMessage string
}

type Better interface {
	Bet(context.Context, *bet.Bet, *user.BetUser) (BetResponse, error)
}
