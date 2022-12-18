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
	"github.com/ramasauskas/ispbet/report"
	"github.com/ramasauskas/ispbet/user"
	"github.com/shopspring/decimal"
)

type adminLog struct {
	UUID      uuid.UUID `json:"uuid"`
	AdminUUID uuid.UUID `json:"admin_uuid"`
	Action    string    `json:"action"`
	Timestamp time.Time `json:"timestamp"`
}

func adminLogView(l user.AdminLog) adminLog {
	return adminLog{
		UUID:      l.UUID,
		AdminUUID: l.AdminUUID,
		Action:    l.Action,
		Timestamp: l.Timestamp,
	}
}

type adminUser struct {
	UUID          uuid.UUID `json:"uuid"`
	Email         string    `json:"email"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	EmailVerified bool      `json:"email_verified"`
	Role          string    `json:"role"`
}

func adminUserView(u user.AdminUser) adminUser {
	return adminUser{
		UUID:          u.UUID,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		EmailVerified: u.EmailVerified,
		Role:          string(u.Role),
	}
}

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

	if be.BeginsAt.Before(time.Now()) {
		return errors.New("begins at cannot be before now")
	}

	if len(be.Name) == 0 {
		return errors.New("name not provided")
	}

	if len(be.Sport) == 0 {
		return errors.New("sport not provided")
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

func betEventSelectionView(s bet.EventSelection) betEventSelection {
	return betEventSelection{
		UUID:     s.UUID,
		Name:     s.Name,
		OddsHome: s.OddsHome,
		OddsAway: s.OddsAway,
		Winner:   s.Winner,
	}
}

func betEventView(e bet.Event) betEvent {
	var selections []betEventSelection

	for _, s := range e.Selections {
		selections = append(selections, betEventSelectionView(s))
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

	r.Post("/login", s.adminLogin)

	r.Group(func(r chi.Router) {
		r.Use(s.sessions.Auth)
		r.Get("/bet-users", s.betUsers)
		r.Get("/admin-logs", s.adminsLogs)
		r.Get("/identity-verifications", s.identityVerifications)

		r.Post("/auto-report", s.createAutoReport)

		r.Post("/finalize-identity-verification", s.authorizeAdmin(user.RoleUsers, "finalize-identity", s.finalizeIdentityVerification))
		r.Post("/deposit", s.authorizeAdmin(user.RoleUsers, "deposit", s.createDeposit))
		r.Post("/withdraw", s.authorizeAdmin(user.RoleUsers, "withdraw", s.createWithdrawal))
		r.Post("/event", s.authorizeAdmin(user.RoleMatches, "create-event", s.createEvent))
		r.Post("/resolve", s.authorizeAdmin(user.RoleMatches, "resolve-event", s.resolveEventSelection))

		r.Route("/report", func(r chi.Router) {
			r.Post("/profit", s.profitReport)
			r.Post("/admins", s.admins)
			r.Post("/admin-logs/{uuid}", s.adminLogs)
			r.Post("/user-bets/{uuid}", s.userBets)
			r.Post("/bets", s.betReport)
		})
	})

	return r
}

func (s *Server) adminLogin(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ctx := r.Context()
	log := s.logger("loginAdmin")

	u, ok, err := s.db.FetchAdminUserByEmail(ctx, input.Email)
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

	respondJSON(w, http.StatusOK, adminUserView(u))
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

func (s *Server) finalizeIdentityVerification(w http.ResponseWriter, r *http.Request, _ user.AdminUser) {
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

func (s *Server) createDeposit(w http.ResponseWriter, r *http.Request, _ user.AdminUser) {
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

func (s *Server) createWithdrawal(w http.ResponseWriter, r *http.Request, _ user.AdminUser) {
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

func (s *Server) createEvent(w http.ResponseWriter, r *http.Request, _ user.AdminUser) {
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

func (s *Server) resolveEventSelection(w http.ResponseWriter, r *http.Request, _ user.AdminUser) {
	var input struct {
		SelectionUUID uuid.UUID  `json:"selection_uuid"`
		Winner        bet.Winner `json:"winner"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	if err := input.Winner.Validate(); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	if input.Winner == bet.WinnerTBD {
		respondErr(w, badRequestErr(errors.New("winner cannot be tbd")))
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

func (s *Server) createAutoReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("createAutoReport")

	var input struct {
		SendTo string `json:"send_to"`
		Type   string `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	if input.SendTo == "" {
		respondErr(w, badRequestErr(errors.New("no send email provided")))
		return
	}

	switch input.Type {
	case "profit", "deposit":
	default:
		respondErr(w, badRequestErr(errors.New("type must be profit or debit")))
	}

	rep := report.AutoReport{
		UUID:   uuid.New(),
		Type:   report.ReportType(input.Type),
		SendTo: input.SendTo,
	}

	if err := s.db.InsertAutoReport(ctx, rep); err != nil {
		log.Error().Err(err).Msg("cannot insert auto report")
		respondErr(w, internalErr())

		return
	}

	respondOK(w)
}

func (s *Server) adminsLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("adminsLog")

	ll, err := s.db.FetchAdminsLogs(ctx)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch admin logs")
		respondErr(w, internalErr())
		return
	}

	views := make([]adminLog, 0)

	for _, l := range ll {
		views = append(views, adminLogView(l))
	}

	respondJSON(w, http.StatusOK, views)
}

func (s *Server) profitReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("adminLog")

	var input struct {
		From time.Time
		To   time.Time
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	profit, err := s.db.FetchProfit(ctx, ProfitOpts{
		From: input.From,
		To:   input.To,
	})
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch profit report")
		respondErr(w, internalErr())

		return
	}

	respondJSON(w, http.StatusOK, profit)
}

func (s *Server) admins(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("admins")

	uu, err := s.db.FetchAdminUsers(ctx)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch admin users")
		respondErr(w, internalErr())
		return
	}

	views := make([]adminUser, 0)

	for _, u := range uu {
		views = append(views, adminUserView(u))
	}

	respondJSON(w, http.StatusOK, views)
}

func (s *Server) adminLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("adminLogs")

	id, err := uuid.Parse(chi.URLParamFromCtx(ctx, "uuid"))
	if err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	ll, err := s.db.FetchAdminLogs(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch admin logs")
		respondErr(w, internalErr())
		return
	}

	views := make([]adminLog, 0)

	for _, l := range ll {
		views = append(views, adminLogView(l))
	}

	respondJSON(w, http.StatusOK, views)
}

func (s *Server) userBets(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("userBets")

	id, err := uuid.Parse(chi.URLParamFromCtx(ctx, "uuid"))
	if err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	bets, err := s.db.FetchUserBets(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch user bets")
		respondErr(w, internalErr())

		return
	}

	views := make([]userBet, 0)

	for _, b := range bets {
		sel, ok, err := s.db.FetchSelection(ctx, b.SelectionUUID)
		if err != nil {
			log.Error().Err(err).Msg("cannot fetch event")
			continue
		}

		if !ok {
			continue
		}

		ev, ok, err := s.db.FetchEvent(ctx, sel.EventUUID)
		if err != nil {
			log.Error().Err(err).Msg("cannot fetch event")
			respondErr(w, internalErr())

			return
		}

		if !ok {
			continue
		}

		betView := userBetView(b, betEventView(ev), betEventSelectionView(sel))
		views = append(views, betView)
	}

	respondJSON(w, http.StatusOK, views)
}

func (s *Server) betReport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := s.logger("betReport")

	var input struct {
		From time.Time
		To   time.Time
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondErr(w, badRequestErr(err))
		return
	}

	bb, err := s.db.FetchBetReport(ctx, input.From, input.To)
	if err != nil {
		log.Error().Err(err).Msg("cannot fetch bet report")
		respondErr(w, internalErr())

		return
	}

	views := make([]userBet, 0)

	for _, b := range bb {
		sel, ok, err := s.db.FetchSelection(ctx, b.SelectionUUID)
		if err != nil {
			log.Error().Err(err).Msg("cannot fetch event")
			continue
		}

		if !ok {
			continue
		}

		ev, ok, err := s.db.FetchEvent(ctx, sel.EventUUID)
		if err != nil {
			log.Error().Err(err).Msg("cannot fetch event")
			respondErr(w, internalErr())

			return
		}

		if !ok {
			continue
		}

		betView := userBetView(b, betEventView(ev), betEventSelectionView(sel))
		views = append(views, betView)
	}

	respondJSON(w, http.StatusOK, views)
}

type Resolver interface {
	Resolve(context.Context, bet.EventSelection) error
}
