package main

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/autobet"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/db"
	"github.com/ramasauskas/ispbet/purse"
	"github.com/ramasauskas/ispbet/report"
	"github.com/ramasauskas/ispbet/server"
	"github.com/ramasauskas/ispbet/user"
)

type serverDBAdapter struct {
	db *db.DB
}

func (a *serverDBAdapter) FetchBetUserByEmail(ctx context.Context, email string) (user.BetUser, bool, error) {
	bu, ok, err := a.db.FetchBetUser(ctx, a.db.NoTX(), db.FetchUserByEmail(email))
	if err != nil {
		return user.BetUser{}, false, err
	}

	if !ok {
		return user.BetUser{}, false, nil
	}

	return decodeBetUser(bu), true, nil
}

func (a *serverDBAdapter) FetchBetUserByUUID(ctx context.Context, userUUID uuid.UUID) (user.BetUser, bool, error) {
	bu, ok, err := a.db.FetchBetUser(ctx, a.db.NoTX(), db.FetchUserByUUID(userUUID))

	if err != nil {
		return user.BetUser{}, false, err
	}

	if !ok {
		return user.BetUser{}, false, nil
	}

	return decodeBetUser(bu), true, nil
}

func (a *serverDBAdapter) InsertBetUser(ctx context.Context, u user.BetUser) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.InsertUser(ctx, tx, encodeUser(u.User)); err != nil {
		return err
	}

	if err = a.db.InsertBetUser(ctx, tx, encodeBetUser(u)); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (a *serverDBAdapter) InsertBetUserIdentityVerification(ctx context.Context, ver user.IdentityVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.InsertIdentityVerification(ctx, tx, encodeIdentityVerification(ver)); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchIdentityVerification(ctx context.Context, id uuid.UUID) (user.IdentityVerification, bool, error) {
	verif, ok, err := a.db.FetchIdentityVerification(ctx, a.db.NoTX(), id)
	if err != nil {
		return user.IdentityVerification{}, false, err
	}

	if !ok {
		return user.IdentityVerification{}, false, nil
	}

	return decodeIdentityVerification(verif), true, nil
}

func (a *serverDBAdapter) FetchIdentityVerifications(ctx context.Context) ([]user.IdentityVerification, error) {
	verifs, err := a.db.FetchIdentityVerifications(ctx, a.db.NoTX())
	if err != nil {
		return nil, err
	}

	vv := make([]user.IdentityVerification, len(verifs))

	for i := range verifs {
		vv[i] = decodeIdentityVerification(verifs[i])
	}

	return vv, nil
}

func (a *serverDBAdapter) InsertIdentityVerificationUpdate(ctx context.Context, u user.BetUser, ver user.IdentityVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.UpdateBetUser(ctx, tx, encodeBetUser(u)); err != nil {
		return err
	}

	if err = a.db.UpdateIdentityVerification(ctx, tx, encodeIdentityVerification(ver)); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) InsertEmailVerification(ctx context.Context, ve user.EmailVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.InsertEmailVerification(ctx, tx, encodeEmailVerification(ve)); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchEmailVerification(ctx context.Context, token string) (user.EmailVerification, bool, error) {
	vee, ok, err := a.db.FetchEmailVerification(ctx, a.db.NoTX(), token)
	if err != nil {
		return user.EmailVerification{}, false, err
	}

	if !ok {
		return user.EmailVerification{}, false, nil
	}

	return decodeEmailVerifcation(vee), true, nil
}

func (a *serverDBAdapter) InsertUserVerification(ctx context.Context, u user.User, ve user.EmailVerification) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := a.db.UpdateUser(ctx, tx, encodeUser(u)); err != nil {
		return err
	}

	if err := a.db.UpdateEmailVerification(ctx, tx, encodeEmailVerification(ve)); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchUserByUUID(ctx context.Context, uuid uuid.UUID) (user.User, bool, error) {
	u, ok, err := a.db.FetchUser(ctx, a.db.NoTX(), db.FetchUserByUUID(uuid))
	if err != nil {
		return user.User{}, false, err
	}

	if !ok {
		return user.User{}, false, nil
	}

	return decodeUser(u), true, nil
}

func (a *serverDBAdapter) InsertDeposit(ctx context.Context, u user.BetUser, d purse.Deposit) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.UpdateBetUser(ctx, tx, encodeBetUser(u)); err != nil {
		return err
	}

	if err = a.db.InsertDeposit(ctx, tx, encodeDeposit(d)); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) InsertWithdrawal(ctx context.Context, u user.BetUser, wd purse.Withdrawal) error {
	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err = a.db.UpdateBetUser(ctx, tx, encodeBetUser(u)); err != nil {
		return err
	}

	if err = a.db.InsertWithdrawal(ctx, tx, encodeWithdrawal(wd)); err != nil {
		return err
	}

	return tx.Commit()
}

func (a *serverDBAdapter) FetchBetUsers(ctx context.Context) ([]user.BetUser, error) {
	uu, err := a.db.FetchBetUsers(ctx, a.db.NoTX())
	if err != nil {
		return nil, err
	}

	uut := make([]user.BetUser, len(uu))

	for i := range uu {
		uut[i] = decodeBetUser(uu[i])
	}

	return uut, nil
}

func (a *serverDBAdapter) InsertEvent(ctx context.Context, ev bet.Event) error {
	homeTeam := encodeTeam(ev.HomeTeam)
	awayTeam := encodeTeam(ev.AwayTeam)

	dbEv := encodeEvent(ev)

	var (
		awayPlayers []db.TeamPlayer
		homePlayers []db.TeamPlayer
		selections  []db.EventSelection
	)

	for _, p := range ev.HomeTeam.Players {
		homePlayers = append(homePlayers, encodePlayer(p, ev.HomeTeam.UUID))
	}

	for _, p := range ev.AwayTeam.Players {
		awayPlayers = append(awayPlayers, encodePlayer(p, ev.AwayTeam.UUID))
	}

	for _, s := range ev.Selections {
		selections = append(selections, encodeSelection(s, ev.UUID))
	}

	tx, err := a.db.NewTX(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if err := a.db.InsertTeam(ctx, tx, homeTeam); err != nil {
		return err
	}

	if err := a.db.InsertTeam(ctx, tx, awayTeam); err != nil {
		return err
	}

	for _, p := range awayPlayers {
		if err := a.db.InsertTeamPlayer(ctx, tx, p); err != nil {
			return err
		}
	}

	for _, p := range homePlayers {
		if err := a.db.InsertTeamPlayer(ctx, tx, p); err != nil {
			return err
		}
	}

	if err := a.db.InsertEvent(ctx, tx, dbEv); err != nil {
		return err
	}

	for _, s := range selections {
		if err := a.db.InsertEventSelection(ctx, tx, s); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (a *serverDBAdapter) UpdateEvent(ctx context.Context, ev bet.Event) error {
	return a.db.UpdateEvent(ctx, a.db.NoTX(), encodeEvent(ev))
}

func (a *serverDBAdapter) UpdateSelection(ctx context.Context, sel bet.EventSelection) error {
	return a.db.UpdateSelection(ctx, a.db.NoTX(), encodeSelection(sel, sel.EventUUID))
}

func (a *serverDBAdapter) FetchEvents(ctx context.Context) ([]bet.Event, error) {
	evs, err := a.db.FetchEvents(ctx, a.db.NoTX(), db.EventNotFinished())
	if err != nil {
		return nil, err
	}

	var decodedEvs []bet.Event

	for i := range evs {
		ev, err := fillEvent(ctx, a.db, a.db.NoTX(), evs[i])
		if err != nil {
			return nil, err
		}

		decodedEvs = append(decodedEvs, ev)
	}

	return decodedEvs, nil
}

func (a *serverDBAdapter) FetchSelection(ctx context.Context, id uuid.UUID) (bet.EventSelection, bool, error) {
	sel, ok, err := a.db.FetchSelectionByUUID(ctx, a.db.NoTX(), id)
	if err != nil {
		return bet.EventSelection{}, false, err
	}

	if !ok {
		return bet.EventSelection{}, false, nil
	}

	return decodeSelection(sel), true, nil
}

func (a *serverDBAdapter) FetchAdminUserByEmail(ctx context.Context, email string) (user.AdminUser, bool, error) {
	u, ok, err := a.db.FetchAdminUser(ctx, a.db.NoTX(), db.FetchUserByEmail(email))
	if err != nil {
		return user.AdminUser{}, false, err
	}

	if !ok {
		return user.AdminUser{}, false, nil
	}

	return decodeAdminUser(u), true, nil
}

func (a *serverDBAdapter) FetchAdminUserByUUID(ctx context.Context, id uuid.UUID) (user.AdminUser, bool, error) {
	u, ok, err := a.db.FetchAdminUser(ctx, a.db.NoTX(), db.FetchUserByUUID(id))
	if err != nil {
		return user.AdminUser{}, false, err
	}

	if !ok {
		return user.AdminUser{}, false, nil
	}

	return decodeAdminUser(u), true, nil
}

func (a *serverDBAdapter) InsertAdminLog(ctx context.Context, lg user.AdminLog) error {
	return a.db.InsertAdminLog(ctx, a.db.NoTX(), encodeAdminLog(lg))
}

func (a *serverDBAdapter) FetchAdminsLogs(ctx context.Context) ([]user.AdminLog, error) {
	logs, err := a.db.FetchAdminsLogs(ctx, a.db.NoTX())
	if err != nil {
		return nil, err
	}

	var ll []user.AdminLog

	for _, l := range logs {
		ll = append(ll, decodeAdminLog(l))
	}

	return ll, nil
}

func (a *serverDBAdapter) FetchUserBets(ctx context.Context, id uuid.UUID) ([]bet.Bet, error) {
	bets, err := a.db.FetchBets(ctx, a.db.NoTX(), db.UserBets(id))
	if err != nil {
		return nil, err
	}

	var bb []bet.Bet

	for _, b := range bets {
		bb = append(bb, decodeBet(b))
	}

	return bb, nil
}

func (a *serverDBAdapter) FetchEvent(ctx context.Context, id uuid.UUID) (bet.Event, bool, error) {
	ev, ok, err := a.db.FetchEvent(ctx, a.db.NoTX(), id)
	if err != nil {
		return bet.Event{}, false, err
	}

	if !ok {
		return bet.Event{}, false, nil
	}

	filled, err := fillEvent(ctx, a.db, a.db.NoTX(), ev)
	if err != nil {
		return bet.Event{}, false, err
	}

	return filled, true, nil
}

func (a *serverDBAdapter) FetchEventBySelection(ctx context.Context, id uuid.UUID) (bet.Event, bool, error) {
	ev, ok, err := a.db.FetchEventBySelection(ctx, a.db.NoTX(), id)
	if err != nil {
		return bet.Event{}, false, err
	}

	if !ok {
		return bet.Event{}, false, nil
	}

	filled, err := fillEvent(ctx, a.db, a.db.NoTX(), ev)
	if err != nil {
		return bet.Event{}, false, err
	}

	return filled, true, nil
}

func (a *serverDBAdapter) FetchProfit(ctx context.Context, po server.ProfitOpts) (server.ProfitReport, error) {
	pr, err := a.db.ProfitReport(ctx, db.ProfitOpts{
		From: po.From,
		To:   po.To,
	})
	if err != nil {
		return server.ProfitReport{}, err
	}

	return server.ProfitReport{
		Profit: pr.Profit,
		Loss:   pr.Loss,
		Final:  pr.Final,
	}, nil
}

func (a *serverDBAdapter) FetchAdminLogs(ctx context.Context, id uuid.UUID) ([]user.AdminLog, error) {
	ll, err := a.db.FetchAdminLog(ctx, id)
	if err != nil {
		return nil, err
	}

	var decoded []user.AdminLog

	for _, l := range ll {
		decoded = append(decoded, decodeAdminLog(l))
	}

	return decoded, nil
}

func (a *serverDBAdapter) FetchAdminUsers(ctx context.Context) ([]user.AdminUser, error) {
	uu, err := a.db.FetchAdmins(ctx)
	if err != nil {
		return nil, err
	}

	var decoded []user.AdminUser

	for _, u := range uu {
		decoded = append(decoded, decodeAdminUser(u))
	}

	return decoded, nil
}

func (a *serverDBAdapter) FetchBetReport(ctx context.Context, from, to time.Time) ([]bet.Bet, error) {
	bb, err := a.db.FetchBetReport(ctx, from, to)
	if err != nil {
		return nil, err
	}

	var bets []bet.Bet

	for _, b := range bb {
		bets = append(bets, decodeBet(b))
	}

	return bets, nil
}

func (a *serverDBAdapter) InsertAutoReport(ctx context.Context, r report.AutoReport) error {
	return a.db.InsertAutoReport(ctx, a.db.NoTX(), db.AutoReport{
		UUID:   r.UUID,
		Type:   string(r.Type),
		SendTo: r.SendTo,
	})
}

func (a *serverDBAdapter) InsertAutoBet(ctx context.Context, au autobet.AutoBet) error {
	return a.db.InsertAutoBet(ctx, a.db.NoTX(), db.AutoBet{
		UUID:            au.UUID,
		HighRisk:        au.HighRisk,
		UserUUID:        au.UserUUID,
		BalanceFraction: au.BalanceFraction,
	})
}

func (a *serverDBAdapter) DeleteAutoBet(ctx context.Context, id uuid.UUID) error {
	return a.db.DeleteAutoBet(ctx, a.db.NoTX(), id)
}

func (a *serverDBAdapter) FetchUserAutoBets(ctx context.Context, id uuid.UUID) ([]autobet.AutoBet, error) {
	bb, err := a.db.FetchUserAutoBets(ctx, id)
	if err != nil {
		return nil, err
	}

	var bbe []autobet.AutoBet

	for _, b := range bb {
		bbe = append(bbe, autobet.AutoBet{
			UUID:            b.UUID,
			HighRisk:        b.HighRisk,
			UserUUID:        b.UserUUID,
			BalanceFraction: b.BalanceFraction,
		})
	}

	return bbe, nil
}

func fillEvent(ctx context.Context, db *db.DB, tx db.TX, ev db.Event) (bet.Event, error) {
	home, ok, err := db.FetchTeamByUUID(ctx, tx, ev.HomeTeamUUID)
	if err != nil {
		return bet.Event{}, err
	}

	if !ok {
		return bet.Event{}, errors.New("not found")
	}

	away, ok, err := db.FetchTeamByUUID(ctx, tx, ev.AwayTeamUUID)
	if err != nil {
		return bet.Event{}, err
	}

	if !ok {
		return bet.Event{}, errors.New("not found")
	}

	homePlayers, err := db.FetchPlayersByTeam(ctx, tx, home.UUID)
	if err != nil {
		return bet.Event{}, err
	}

	awayPlayers, err := db.FetchPlayersByTeam(ctx, tx, away.UUID)
	if err != nil {
		return bet.Event{}, err
	}

	sels, err := db.FetchSelectionsByEvent(ctx, tx, ev.UUID)
	if err != nil {
		return bet.Event{}, err
	}

	var (
		decodedSels []bet.EventSelection
		decodedAway []bet.Player
		decodedHome []bet.Player
	)

	for _, sel := range sels {
		decodedSels = append(decodedSels, decodeSelection(sel))
	}

	for _, p := range awayPlayers {
		decodedAway = append(decodedAway, decodePlayer(p))
	}

	for _, p := range homePlayers {
		decodedHome = append(decodedHome, decodePlayer(p))
	}

	return bet.Event{
		UUID:       ev.UUID,
		Name:       ev.Name,
		Sport:      bet.Sport(ev.Sport),
		Selections: decodedSels,
		BeginsAt:   ev.BeginsAt,
		HomeTeam:   decodeTeam(home, decodedHome),
		AwayTeam:   decodeTeam(away, decodedAway),
	}, nil
}

func encodeAdminLog(lg user.AdminLog) db.AdminLog {
	return db.AdminLog{
		UUID:      lg.UUID,
		AdminUUID: lg.AdminUUID,
		Action:    lg.Action,
		Timestamp: lg.Timestamp,
	}
}

func decodeAdminLog(lg db.AdminLog) user.AdminLog {
	return user.AdminLog{
		UUID:      lg.UUID,
		AdminUUID: lg.AdminUUID,
		Action:    lg.Action,
		Timestamp: lg.Timestamp,
	}
}

func decodeAdminUser(u db.AdminUser) user.AdminUser {
	return user.AdminUser{
		User: decodeUser(u.User),
		Role: user.Role(u.Role),
	}
}

func decodeUser(u db.User) user.User {
	return user.User{
		UUID:             u.UUID,
		Email:            u.Email,
		PasswordHash:     u.PasswordHash,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		EmailVerified:    u.EmailVerified,
		IdentityVerified: u.IdentityVerified,
	}
}

func decodeBetUser(u db.BetUser) user.BetUser {
	return user.BetUser{
		User:             decodeUser(u.User),
		IdentityVerified: u.IdentityVerified,
		Balance:          u.Balance,
	}
}

func encodeUser(u user.User) db.User {
	return db.User{
		UUID:             u.UUID,
		Email:            u.Email,
		PasswordHash:     u.PasswordHash,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		EmailVerified:    u.EmailVerified,
		IdentityVerified: u.IdentityVerified,
	}
}

func encodeBetUser(u user.BetUser) db.BetUser {
	return db.BetUser{
		User:             encodeUser(u.User),
		IdentityVerified: u.IdentityVerified,
		Balance:          u.Balance,
	}
}

func encodeEmailVerification(ev user.EmailVerification) db.EmailVerification {
	return db.EmailVerification{
		UserUUID:  ev.UserUUID,
		Token:     ev.Token,
		Activated: ev.Activated,
	}
}

func decodeEmailVerifcation(ev db.EmailVerification) user.EmailVerification {
	return user.EmailVerification{
		UserUUID:  ev.UserUUID,
		Token:     ev.Token,
		Activated: ev.Activated,
	}
}

func encodeIdentityVerification(idv user.IdentityVerification) db.IdentityVerification {
	return db.IdentityVerification{
		UUID:                idv.UUID,
		UserUUID:            idv.UserUUID,
		Status:              string(idv.Status),
		IDPhotoBase64:       idv.IDPhotoBase64,
		PortraitPhotoBase64: idv.PortraitPhotoBase64,
		RespondedAt:         idv.RespondedAt,
		CreatedAt:           idv.CreatedAt,
	}
}

func decodeIdentityVerification(idv db.IdentityVerification) user.IdentityVerification {
	return user.IdentityVerification{
		UUID:                idv.UUID,
		UserUUID:            idv.UserUUID,
		Status:              user.IdentityVerificationStatus(idv.Status),
		IDPhotoBase64:       idv.IDPhotoBase64,
		PortraitPhotoBase64: idv.PortraitPhotoBase64,
		RespondedAt:         idv.RespondedAt,
		CreatedAt:           idv.CreatedAt,
	}
}

func encodeDeposit(d purse.Deposit) db.Deposit {
	return db.Deposit{
		UUID:      d.UUID,
		Amount:    d.Amount,
		Timestamp: d.Timestamp,
		UserUUID:  d.UserUUID,
	}
}

func decodeDeposit(d db.Deposit) purse.Deposit {
	return purse.Deposit{
		UUID:      d.UUID,
		Amount:    d.Amount,
		Timestamp: d.Timestamp,
		UserUUID:  d.UserUUID,
	}
}

func encodeWithdrawal(wd purse.Withdrawal) db.Withdrawal {
	return db.Withdrawal{
		UUID:      wd.UUID,
		Amount:    wd.Amount,
		Timestamp: wd.Timestamp,
		UserUUID:  wd.UserUUID,
	}
}

func decodeWithdrawal(wd db.Withdrawal) purse.Withdrawal {
	return purse.Withdrawal{
		UUID:      wd.UUID,
		Amount:    wd.Amount,
		Timestamp: wd.Timestamp,
		UserUUID:  wd.UserUUID,
	}
}

func encodeEvent(ev bet.Event) db.Event {
	return db.Event{
		UUID:         ev.UUID,
		Name:         ev.Name,
		Sport:        string(ev.Sport),
		BeginsAt:     ev.BeginsAt,
		Finished:     ev.Finished(),
		HomeTeamUUID: ev.HomeTeam.UUID,
		AwayTeamUUID: ev.AwayTeam.UUID,
	}
}

func decodeEvent(ev db.Event, home bet.Team, away bet.Team) bet.Event {
	return bet.Event{
		UUID:     ev.UUID,
		Name:     ev.Name,
		Sport:    bet.Sport(ev.Sport),
		BeginsAt: ev.BeginsAt,
		HomeTeam: home,
		AwayTeam: away,
	}
}

func encodeSelection(sel bet.EventSelection, eventUUID uuid.UUID) db.EventSelection {
	return db.EventSelection{
		UUID:      sel.UUID,
		EventUUID: eventUUID,
		Name:      sel.Name,
		OddsHome:  sel.OddsHome,
		AutoOdds:  sel.AutoOdds,
		OddsAway:  sel.OddsAway,
		Winner:    string(sel.Winner),
	}
}

func decodeSelection(sel db.EventSelection) bet.EventSelection {
	return bet.EventSelection{
		UUID:      sel.UUID,
		Name:      sel.Name,
		OddsHome:  sel.OddsHome,
		OddsAway:  sel.OddsAway,
		EventUUID: sel.EventUUID,
		AutoOdds:  sel.AutoOdds,
		Winner:    bet.Winner(sel.Winner),
	}
}

func encodeTeam(t bet.Team) db.Team {
	return db.Team{
		UUID: t.UUID,
		Name: t.Name,
	}
}

func decodeTeam(t db.Team, pp []bet.Player) bet.Team {
	return bet.Team{
		UUID:    t.UUID,
		Name:    t.Name,
		Players: pp,
	}
}

func encodePlayer(tp bet.Player, teamUUID uuid.UUID) db.TeamPlayer {
	return db.TeamPlayer{
		UUID:     tp.UUID,
		TeamUUID: teamUUID,
		Name:     tp.Name,
	}
}

func decodePlayer(tp db.TeamPlayer) bet.Player {
	return bet.Player{
		UUID: tp.UUID,
		Name: tp.Name,
	}
}
