package db

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

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

func (d *DB) FetchEvents(ctx context.Context, q sq.QueryerContext) ([]Event, error) {
	b := sq.Select()

	b = eventQuery(b, "betev").From("bet_event AS betev")
	qr, _ := b.MustSql()

	var ee []Event

	if err := d.d.SelectContext(ctx, &ee, qr); err != nil {
		return nil, err
	}

	return ee, nil
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
