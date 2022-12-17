package db

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Event struct {
	UUID         uuid.UUID `db:"betev.uuid"`
	Name         string    `db:"betev.name"`
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

func eventQuery(b sq.SelectBuilder, prefix string) sq.SelectBuilder {
	return b.Columns(
		column(prefix, "uuid"),
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
