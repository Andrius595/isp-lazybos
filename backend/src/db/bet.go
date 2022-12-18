package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AutoBet struct {
	UUID            uuid.UUID       `db:"ab.uuid"`
	HighRisk        bool            `db:"ab.high_risk"`
	UserUUID        uuid.UUID       `db:"ab.user_uuid"`
	BalanceFraction decimal.Decimal `db:"ab.balance_fraction"`
}

type fetchEventCriteria func(b sq.SelectBuilder, prefix string) sq.SelectBuilder

func EventNotFinished() fetchEventCriteria {
	return func(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
		return b.Where(sq.NotEq{columnPredicate(prefix, "finished"): true})
	}
}

type fetchBetCriteria func(b sq.SelectBuilder, prefix string) sq.SelectBuilder

func UserBets(id uuid.UUID) fetchBetCriteria {
	return func(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
		return b.Where(sq.Eq{columnPredicate(prefix, "user_uuid"): id})
	}
}

func SelectionBets(id uuid.UUID) fetchBetCriteria {
	return func(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
		return b.Where(sq.Eq{columnPredicate(prefix, "selection_uuid"): id})
	}
}

type Bet struct {
	UUID            uuid.UUID       `db:"bt.uuid"`
	UserUUID        uuid.UUID       `db:"bt.user_uuid"`
	SelectionUUID   uuid.UUID       `db:"bt.selection_uuid"`
	Timestamp       time.Time       `db:"bt.timestamp"`
	SelectionWinner string          `db:"bt.selection_winner"`
	Stake           decimal.Decimal `db:"bt.stake"`
	Odds            decimal.Decimal `db:"bt.odds"`
	State           string          `db:"bt.state"`
}

type Event struct {
	UUID         uuid.UUID `db:"betev.uuid"`
	Name         string    `db:"betev.name"`
	Sport        string    `db:"betev.sport_name"`
	BeginsAt     time.Time `db:"betev.begins_at"`
	Finished     bool      `db:"betev.finished"`
	HomeTeamUUID uuid.UUID `db:"betev.home_team_uuid"`
	AwayTeamUUID uuid.UUID `db:"betev.away_team_uuid"`
}

type EventSelection struct {
	UUID      uuid.UUID       `db:"es.uuid"`
	EventUUID uuid.UUID       `db:"es.event_uuid"`
	Name      string          `db:"es.name"`
	OddsHome  decimal.Decimal `db:"es.odds_home"`
	OddsAway  decimal.Decimal `db:"es.odds_away"`
	AutoOdds  bool            `db:"es.auto_odds"`
	Winner    string          `db:"es.winner"`
}

type Team struct {
	UUID uuid.UUID `db:"tm.uuid"`
	Name string    `db:"tm.name"`
}

type TeamPlayer struct {
	UUID     uuid.UUID `db:"tmp.uuid"`
	TeamUUID uuid.UUID `db:"tmp.team_uuid"`
	Name     string    `db:"tmp.name"`
}

func (d *DB) InsertEvent(ctx context.Context, e sq.ExecerContext, ev Event) error {
	b := sq.Insert("bet_event").SetMap(map[string]interface{}{
		"uuid":           ev.UUID,
		"name":           ev.Name,
		"sport_name":     ev.Sport,
		"begins_at":      ev.BeginsAt,
		"finished":       ev.Finished,
		"home_team_uuid": ev.HomeTeamUUID,
		"away_team_uuid": ev.AwayTeamUUID,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) InsertEventSelection(ctx context.Context, e sq.ExecerContext, se EventSelection) error {
	b := sq.Insert("event_selection").SetMap(map[string]interface{}{
		"uuid":       se.UUID,
		"name":       se.Name,
		"odds_home":  se.OddsHome,
		"odds_away":  se.OddsAway,
		"auto_odds":  se.AutoOdds,
		"winner":     se.Winner,
		"event_uuid": se.EventUUID,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) InsertTeam(ctx context.Context, e sq.ExecerContext, tm Team) error {
	b := sq.Insert("team").SetMap(map[string]interface{}{
		"uuid": tm.UUID,
		"name": tm.Name,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) InsertTeamPlayer(ctx context.Context, e sq.ExecerContext, tp TeamPlayer) error {
	b := sq.Insert("team_player").SetMap(map[string]interface{}{
		"uuid":      tp.UUID,
		"team_uuid": tp.TeamUUID,
		"name":      tp.Name,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchEvents(ctx context.Context, q sq.QueryerContext, c fetchEventCriteria) ([]Event, error) {
	b := sq.Select()

	b = c(eventQuery(b, "betev").From("bet_event AS betev"), "betev")
	qr, args := b.MustSql()

	var ee []Event

	if err := d.d.SelectContext(ctx, &ee, qr, args...); err != nil {
		return nil, err
	}

	return ee, nil
}

func (d *DB) FetchEvent(ctx context.Context, q sq.QueryerContext, id uuid.UUID) (Event, bool, error) {
	b := sq.Select()

	b = eventQuery(b, "betev").From("bet_event AS betev").Where(sq.Eq{"betev.uuid": id})
	qr, _ := b.MustSql()

	var ee Event

	err := d.d.GetContext(ctx, &ee, qr)
	switch err {
	case nil:
		return ee, true, nil
	case sql.ErrNoRows:
		return Event{}, false, nil
	default:
		return Event{}, false, err
	}
}

func (d *DB) FetchSelectionsByEvent(ctx context.Context, q sq.QueryerContext, id uuid.UUID) ([]EventSelection, error) {
	b := sq.Select()

	b = selectionQuery(b, "es").From("event_selection AS es").Where(sq.Eq{"es.event_uuid": id})
	qr, args := b.MustSql()

	var ss []EventSelection

	if err := d.d.SelectContext(ctx, &ss, qr, args...); err != nil {
		return nil, err
	}

	return ss, nil
}

func (d *DB) FetchSelections(ctx context.Context) ([]EventSelection, error) {
	b := sq.Select()

	b = selectionQuery(b, "es").From("event_selection AS es")
	qr, _ := b.MustSql()

	var es []EventSelection

	if err := d.d.SelectContext(ctx, &es, qr); err != nil {
		return nil, err
	}

	return es, nil
}

func (d *DB) FetchSelectionByUUID(ctx context.Context, q sq.QueryerContext, id uuid.UUID) (EventSelection, bool, error) {
	b := sq.Select()

	b = selectionQuery(b, "es").From("event_selection AS es").Where(sq.Eq{"es.uuid": id})
	qr, args := b.MustSql()

	var es EventSelection

	err := d.d.GetContext(ctx, &es, qr, args...)
	switch err {
	case nil:
		return es, true, nil
	case sql.ErrNoRows:
		return EventSelection{}, false, nil
	default:
		return EventSelection{}, false, err
	}
}

func (d *DB) FetchTeamByUUID(ctx context.Context, q sq.QueryerContext, id uuid.UUID) (Team, bool, error) {
	b := sq.Select()

	b = teamQuery(b, "tm").From("team AS tm").Where(sq.Eq{"tm.uuid": id})
	qr, args := b.MustSql()

	var tm Team

	err := d.d.GetContext(ctx, &tm, qr, args...)
	switch err {
	case nil:
		return tm, true, nil
	case sql.ErrNoRows:
		return Team{}, false, nil
	default:
		return Team{}, false, err
	}
}

func (d *DB) FetchPlayersByTeam(ctx context.Context, q sq.QueryerContext, id uuid.UUID) ([]TeamPlayer, error) {
	b := sq.Select()

	b = teamQuery(b, "tmp").From("team_player AS tmp").Where(sq.Eq{"tmp.team_uuid": id})
	qr, args := b.MustSql()

	var pp []TeamPlayer

	if err := d.d.SelectContext(ctx, &pp, qr, args...); err != nil {
		return nil, err
	}

	return pp, nil
}

func (d *DB) InsertBet(ctx context.Context, e sq.ExecerContext, bt Bet) error {
	b := sq.Insert("bet").SetMap(map[string]interface{}{
		"uuid":             bt.UUID,
		"user_uuid":        bt.UserUUID,
		"selection_uuid":   bt.SelectionUUID,
		"selection_winner": bt.SelectionWinner,
		"stake":            bt.Stake,
		"odds":             bt.Odds,
		"state":            bt.State,
		"timestamp":        bt.Timestamp,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchBets(ctx context.Context, q sq.QueryerContext, c fetchBetCriteria) ([]Bet, error) {
	b := sq.Select()

	b = c(betQuery(b, "bt").From("bet AS bt"), "bt")
	qr, args := b.MustSql()

	var bb []Bet

	if err := d.d.SelectContext(ctx, &bb, qr, args...); err != nil {
		return nil, err
	}

	return bb, nil
}

func (d *DB) UpdateBet(ctx context.Context, e sq.ExecerContext, bt Bet) error {
	b := sq.Update("bet").SetMap(map[string]interface{}{
		"state": bt.State,
	}).Where(sq.Eq{"uuid": bt.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) UpdateEvent(ctx context.Context, e sq.ExecerContext, ev Event) error {
	b := sq.Update("event").SetMap(map[string]interface{}{
		"finished": ev.Finished,
	}).Where(sq.Eq{"uuid": ev.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) UpdateSelection(ctx context.Context, e sq.ExecerContext, sel EventSelection) error {
	b := sq.Update("event_selection").SetMap(map[string]interface{}{
		"winner":    sel.Winner,
		"odds_away": sel.OddsAway,
		"odds_home": sel.OddsHome,
	}).Where(sq.Eq{"uuid": sel.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchSelectionBest(ctx context.Context, highRisk bool) (EventSelection, bool, error) {
	b := sq.Select()

	asc := "DESC"
	if !highRisk {
		asc = "ASC"
	}

	orderBy := fmt.Sprintf("MAX(es.odds_away, es.odds_home) %s", asc)
	b = selectionQuery(b, "es").From("event_selection AS es").OrderBy(orderBy).Limit(1)
	qr, args := b.MustSql()

	var sel EventSelection

	err := d.d.GetContext(ctx, &sel, qr, &args)
	switch err {
	case nil:
		return sel, true, nil
	case sql.ErrNoRows:
		return EventSelection{}, false, nil
	default:
		return EventSelection{}, false, err
	}
}

func (d *DB) FetchAutoBets(ctx context.Context) ([]AutoBet, error) {
	b := sq.Select()

	b = autoBetQuery(b, "ab").From("auto_bet AS ab")
	qr, args := b.MustSql()

	var bb []AutoBet

	if err := d.d.SelectContext(ctx, &bb, qr, args...); err != nil {
		return nil, err
	}

	return bb, nil
}

func (d *DB) FetchUserAutoBets(ctx context.Context, id uuid.UUID) ([]AutoBet, error) {
	b := sq.Select()

	b = autoBetQuery(b, "ab").From("auto_bet AS ab").Where(sq.Eq{"ab.user_uuid": id})
	qr, args := b.MustSql()

	var bb []AutoBet

	if err := d.d.SelectContext(ctx, &bb, qr, args...); err != nil {
		return nil, err
	}

	return bb, nil
}

func (d *DB) InsertAutoBet(ctx context.Context, e sq.ExecerContext, ab AutoBet) error {
	b := sq.Insert("auto_bet").SetMap(map[string]interface{}{
		"uuid":             ab.UUID,
		"user_uuid":        ab.UserUUID,
		"high_risk":        ab.HighRisk,
		"balance_fraction": ab.BalanceFraction,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) DeleteAutoBet(ctx context.Context, e sq.ExecerContext, id uuid.UUID) error {
	b := sq.Delete("auto_bet").Where(sq.Eq{"uuid": id})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func autoBetQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
		column(prefix, "high_risk"),
		column(prefix, "user_uuid"),
		column(prefix, "balance_fraction"),
	)
}

func eventQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
		column(prefix, "sport_name"),
		column(prefix, "name"),
		column(prefix, "begins_at"),
		column(prefix, "finished"),
		column(prefix, "home_team_uuid"),
		column(prefix, "away_team_uuid"),
	)
}

func selectionQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
		column(prefix, "name"),
		column(prefix, "odds_home"),
		column(prefix, "auto_odds"),
		column(prefix, "odds_away"),
		column(prefix, "winner"),
		column(prefix, "event_uuid"),
	)
}

func teamQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
		column(prefix, "name"),
	)
}

func teamPlayerQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
		column(prefix, "name"),
		column(prefix, "team_uuid"),
	)
}

func betQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
		column(prefix, "user_uuid"),
		column(prefix, "selection_uuid"),
		column(prefix, "selection_winner"),
		column(prefix, "stake"),
		column(prefix, "odds"),
		column(prefix, "state"),
		column(prefix, "timestamp"),
	)
}
