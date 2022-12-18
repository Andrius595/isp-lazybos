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
	"github.com/ramasauskas/ispbet/purse"
	"github.com/ramasauskas/ispbet/user"
	"github.com/shopspring/decimal"
)

type newBetEvent struct {
	Name       string `json:"name"`
	Sport      string `json:"sport"`
	Selections []struct {
		Name     string          `json:"name"`
		OddsHome decimal.Decimal `json:"odds_home"`
		OddsAway decimal.Decimal `json:"odds_away"`
	} `json:"selections"`
	AwayTeam struct {
		Name    string   `json:"name"`
		Players []string `json:"players"`
	} `json:"away_team"`
	HomeTeam struct {
		Name    string   `json:"name"`
		Players []string `json:"players"`
	} `json:"home_team"`
	BeginsAt time.Time `json:"begins_at"`
}

func (be newBetEvent) validate() error {
	if len(be.Selections) == 0 {
		return errors.New("no selections provided")
	}

	if len(be.AwayTeam.Players) == 0 {
		return errors.New("away team has no players")
	}

	if len(be.HomeTeam.Players) == 0 {
		return errors.New("home team has no players")
	}

	return nil
}

type betEvent struct {
	UUID       uuid.UUID           `json:"uuid"`
	Name       string              `json:"name"`
	Selections []betEventSelection `json:"selections"`
	Sport      string              `json:"sport"`
	BeginsAt   time.Time           `json:"begins_at"`
	Finished   bool                `json:"finished"`
	HomeTeam   betEventTeam        `json:"home_team"`
	AwayTeam   betEventTeam        `json:"away_team"`
}

func betEventView(e bet.Event) betEvent {
	var selections []betEventSelection

	for _, s := range e.Selections {
		selections = append(selections, betEventSelection{
			UUID:     s.UUID,
			Name:     s.Name,
			OddsHome: s.OddsHome,
			OddsAway: s.OddsAway,
			Winner:   s.Winner,
		})
	}

	var (
		homePlayers []betEventPlayer
		awayPlayers []betEventPlayer
	)

	for _, p := range e.AwayTeam.Players {
		awayPlayers = append(awayPlayers, betEventPlayer{
			UUID: p.UUID,
			Name: p.Name,
		})
	}

	for _, p := range e.HomeTeam.Players {
		homePlayers = append(homePlayers, betEventPlayer{
			UUID: p.UUID,
			Name: p.Name,
		})
	}

	home := betEventTeam{
		UUID:    e.HomeTeam.UUID,
		Name:    e.HomeTeam.Name,
		Players: homePlayers,
	}

	away := betEventTeam{
		UUID:    e.HomeTeam.UUID,
		Name:    e.HomeTeam.Name,
		Players: awayPlayers,
	}

	return betEvent{
		UUID:       e.UUID,
		Name:       e.Name,
		Selections: selections,
		Sport:      string(e.Sport),
		BeginsAt:   e.BeginsAt,
		Finished:   e.Finished(),
		HomeTeam:   home,
		AwayTeam:   away,
	}
}

type betEventTeam struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`

	Players []betEventPlayer `json:"players"`
}

type betEventPlayer struct {
	UUID uuid.UUID `json:"uuid"`
	Name string    `json:"name"`
}

type betEventSelection struct {
	UUID     uuid.UUID       `json:"uuid"`
	Name     string          `json:"name"`
	OddsHome decimal.Decimal `json:"odds_home"`
	OddsAway decimal.Decimal `json:"odds_away"`
	Winner   bet.Winner      `json:"winner"`
}

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
	r.Get("/identity-verifications", s.identityVerifications)
	r.Post("/finalize-identity-verification", s.finalizeIdentityVerification)
	r.Post("/deposit", s.createDeposit)
	r.Post("/withdraw", s.createWithdrawal)
	r.Post("/event", s.createEvent)
	r.Post("/resolve", s.resolveEventSelection)

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
		bu, ok, err := s.db.FetchBetUserByUUID(ctx, v.UserUUID)
		if err != nil {
			log.Error().Err(err).Msg("cannot fetch user")
			respondErr(w, internalErr())

			return
		}

		if !ok {
			log.Warn().Stringer("uuid", bu.UUID).Msg("cannot find user")
			continue
		}

		verViews[i] = identityVerification{
			UUID:                v.UUID,
			User:                betUserView(bu),
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

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent newBetEvent

	if err := json.NewDecoder(r.Body).Decode(&newEvent); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	if err := newEvent.validate(); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ev := bet.Event{
		UUID:     uuid.New(),
		Name:     newEvent.Name,
		Sport:    bet.Sport(newEvent.Sport),
		BeginsAt: newEvent.BeginsAt,
	}

	for _, s := range newEvent.Selections {
		ev.Selections = append(ev.Selections, bet.EventSelection{
			UUID:     uuid.New(),
			Name:     s.Name,
			OddsHome: s.OddsHome,
			OddsAway: s.OddsAway,
			Winner:   bet.WinnerTBD,
		})
	}

	var (
		homePlayers []bet.Player
		awayPlayers []bet.Player
	)

	for _, p := range newEvent.HomeTeam.Players {
		homePlayers = append(homePlayers, bet.Player{
			UUID: uuid.New(),
			Name: p,
		})
	}

	for _, p := range newEvent.AwayTeam.Players {
		awayPlayers = append(awayPlayers, bet.Player{
			UUID: uuid.New(),
			Name: p,
		})
	}

	ev.AwayTeam = bet.Team{
		UUID:    uuid.New(),
		Name:    newEvent.AwayTeam.Name,
		Players: awayPlayers,
	}

	ev.HomeTeam = bet.Team{
		UUID:    uuid.New(),
		Name:    newEvent.HomeTeam.Name,
		Players: homePlayers,
	}

	ctx := r.Context()
	log := s.logger("createEvent")

	if err := s.db.InsertEvent(ctx, ev); err != nil {
		log.Error().Err(err).Msg("cannot insert event")
		respondErr(w, internalErr())

		return
	}

	respondJSON(w, http.StatusOK, betEventView(ev))
}

func (s *Server) resolveEventSelection(w http.ResponseWriter, r *http.Request) {
	var input struct {
		SelectionUUID uuid.UUID  `json:"selection_uuid"`
		Winner        bet.Winner `json:"winner"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ctx := r.Context()
	log := s.logger("resolveEventSelection")

	sel, ok, err := s.db.FetchSelection(ctx, input.SelectionUUID)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch event selection")
		respondErr(w, internalErr())

		return
	}

	if !ok {
		respondErr(w, notFoundErr())
		return
	}

	if sel.Winner.Finalized() {
		respondOK(w)
		return
	}

	sel.Winner = input.Winner

	if err := s.resolver.Resolve(ctx, sel); err != nil {
		log.Error().Err(err).Msg("cannot resolve")
		respondErr(w, internalErr())

		return
	}

	respondOK(w)
}

type Resolver interface {
	Resolve(context.Context, bet.EventSelection) error
}
