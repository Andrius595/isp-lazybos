package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramasauskas/ispbet/bet"
	"github.com/ramasauskas/ispbet/db"
	"github.com/ramasauskas/ispbet/purse"
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

func (a *serverDBAdapter) FetchEvents(ctx context.Context) ([]bet.Event, error) {
	evs, err := a.db.FetchEvents(ctx, a.db.NoTX())
	if err != nil {
		return nil, err
	}

	var decodedEvs []bet.Event

	for i := range evs {
		home, ok, err := a.db.FetchTeamByUUID(ctx, a.db.NoTX(), evs[i].HomeTeamUUID)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		away, ok, err := a.db.FetchTeamByUUID(ctx, a.db.NoTX(), evs[i].AwayTeamUUID)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		homePlayers, err := a.db.FetchPlayersByTeam(ctx, a.db.NoTX(), home.UUID)
		if err != nil {
			return nil, err
		}

		awayPlayers, err := a.db.FetchPlayersByTeam(ctx, a.db.NoTX(), away.UUID)
		if err != nil {
			return nil, err
		}

		sels, err := a.db.FetchSelectionsByEvent(ctx, a.db.NoTX(), evs[i].UUID)
		if err != nil {
			return nil, err
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

		ev := bet.Event{
			UUID:       evs[i].UUID,
			Name:       evs[i].Name,
			Selections: decodedSels,
			BeginsAt:   evs[i].BeginsAt,
			Finished:   evs[i].Finished,
			HomeTeam:   decodeTeam(home, decodedHome),
			AwayTeam:   decodeTeam(away, decodedAway),
		}

		decodedEvs = append(decodedEvs, ev)
	}

	return decodedEvs, nil
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
		BeginsAt:     ev.BeginsAt,
		Finished:     ev.Finished,
		HomeTeamUUID: ev.HomeTeam.UUID,
		AwayTeamUUID: ev.AwayTeam.UUID,
	}
}

func decodeEvent(ev db.Event, home bet.Team, away bet.Team) bet.Event {
	return bet.Event{
		UUID:     ev.UUID,
		Name:     ev.Name,
		BeginsAt: ev.BeginsAt,
		Finished: ev.Finished,
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
		OddsAway:  sel.OddsAway,
		Winner:    string(sel.Winner),
	}
}

func decodeSelection(sel db.EventSelection) bet.EventSelection {
	return bet.EventSelection{
		UUID:     sel.UUID,
		Name:     sel.Name,
		OddsHome: sel.OddsHome,
		OddsAway: sel.OddsAway,
		Winner:   bet.Winner(sel.Winner),
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
