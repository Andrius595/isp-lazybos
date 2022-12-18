package db

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type fetchEventCriteria func(b sq.SelectBuilder, prefix string) sq.SelectBuilder

func EventNotFinished() fetchEventCriteria {
	return func(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
		return b.Where(sq.NotEq{columnPredicate(prefix, "finished"): true})
	}
}

type Bet struct {
	UUID            uuid.UUID       `db:"bt.uuid"`
	UserUUID        uuid.UUID       `db:"bt.user_uuid"`
	SelectionUUID   uuid.UUID       `db:"bt.selection_uuid"`
	SelectionWinner string          `db:"bt.selection_winner"`
	Stake           decimal.Decimal `db:"bt.stake"`
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
		"state":            bt.State,
	})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
}

func (d *DB) FetchBetsBySelection(ctx context.Context, q sq.QueryerContext, id uuid.UUID) ([]Bet, error) {
	b := sq.Select()

	b = betQuery(b, "bt").From("bet AS bt").Where(sq.Eq{"bt.selection_uuid": id})
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
		"winner": sel.Winner,
	}).Where(sq.Eq{"uuid": sel.UUID})

	_, err := sq.ExecContextWith(ctx, e, b)
	return err
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
		column(prefix, "state"),
	)
}
