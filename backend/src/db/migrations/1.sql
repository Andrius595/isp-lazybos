-- +migrate Up
CREATE TABLE IF NOT EXISTS user (
	uuid TEXT PRIMARY KEY NOT NULL,
	email TEXT NOT NULL,
	password_hash TEXT NOT NULL,
	first_name TEXT NOT NULL,
	last_name TEXT NOT NULL,
	email_verified INTEGER NOT NULL,

	CONSTRAINT uniq_email UNIQUE(email)
);

CREATE TABLE IF NOT EXISTS bet_user (
	user_uuid TEXT PRIMARY KEY NOT NULL,
	identity_verified INTEGER NOT NULL,
	balance NUMERIC NOT NULL,

	CONSTRAINT fk_bet_user_user_uuid_user FOREIGN KEY(user_uuid) REFERENCES user(uuid)
);

CREATE TABLE IF NOT EXISTS identity_verification (
	uuid TEXT PRIMARY KEY NOT NULL,
	user_uuid TEXT NOT NULL,
	status TEXT NOT NULL,
	id_photo_base64 TEXT NOT NULL,
	portrait_photo_base64 TEXT NOT NULL,
	responded_at TIMESTAMP NOT NULL,
	created_at TIMESTAMP NOT NULL,

	CONSTRAINT fk_user_uuid_bet_user_user_uuid FOREIGN KEY(user_uuid) REFERENCES bet_user(user_uuid)
);

CREATE TABLE IF NOT EXISTS email_verification_token (
	token TEXT PRIMARY KEY NOT NULL,
	user_uuid TEXT NOT NULL,
	activated BOOLEAN NOT NULL,
	
	CONSTRAINT fk_user_uuid_user_user_uuid FOREIGN KEY(user_uuid) REFERENCES user(user_uuid)
);

CREATE TABLE IF NOT EXISTS deposit (
	uuid TEXT PRIMARY KEY NOT NULL,
	amount NUMERIC NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	user_uuid TEXT NOT NULL,

	CONSTRAINT fk_user_uuid_user_user_uuid FOREIGN KEY(user_uuid) REFERENCES bet_user(user_uuid)
);

CREATE TABLE IF NOT EXISTS withdrawal (
	uuid TEXT PRIMARY KEY NOT NULL,
	amount NUMERIC NOT NULL,
	timestamp TIMESTAMP NOT NULL,
	user_uuid TEXT NOT NULL,

	CONSTRAINT fk_user_uuid_user_user_uuid FOREIGN KEY(user_uuid) REFERENCES bet_user(user_uuid)
);

CREATE TABLE IF NOT EXISTS bet_event (
	uuid TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	begins_at TIMESTAMP NOT NULL,
	finished BOOLEAN NOT NULL,
	sport_name TEXT NOT NULL,

	home_team_uuid TEXT NOT NULL,
	away_team_uuid TEXT NOT NULL,

	CONSTRAINT fk_sport_sport_name FOREIGN KEY(sport_name) REFERENCES sport(name),
	CONSTRAINT fk_home_team_uuid_team_uuid FOREIGN KEY(home_team_uuid) REFERENCES team(uuid),
	CONSTRAINT fk_away_team_uuid_team_uuid FOREIGN KEY(away_team_uuid) REFERENCES team(uuid)
);

CREATE TABLE IF NOT EXISTS event_selection (
	uuid TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	odds_home NUMERIC NOT NULL,
	odds_away NUMERIC NOT NULL,
	winner TEXT NOT NULL,
	event_uuid TEXT NOT NULL,

	CONSTRAINT fk_event_uuid_event_uuid FOREIGN KEY(event_uuid) REFERENCES bet_event(uuid)
);

CREATE TABLE IF NOT EXISTS team (
	uuid TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS team_player (
	uuid TEXT PRIMARY KEY NOT NULL,
	name TEXT NOT NULL,
	team_uuid TEXT NOT NULL,

	CONSTRAINT fk_team_uuid_team_uuid FOREIGN KEY(team_uuid) REFERENCES team(uuid)
);

CREATE TABLE IF NOT EXISTS sport (
	name TEXT PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS bet (
	uuid TEXT PRIMARY KEY NOT NULL,
	user_uuid TEXT NOT NULL,
	stake NUMERIC NOT NULL,
	selection_uuid TEXT NOT NULL,
	selection_winner TEXT NOT NULL,

	CONSTRAINT fk_bet_selection_uuid_bet_selection_uuid FOREIGN KEY(selection_uuid) REFERENCES event_selection(uuid),
	CONSTRAINT fk_user_id_user_user_uuid FOREIGN KEY(user_uuid) REFERENCES bet_user(user_uuid)
);

INSERT INTO sport VALUES ("football"), ("basketball");

-- +migrate Down
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS bet_user;
DROP TABLE IF EXISTS identity_verification;
DROP TABLE IF EXISTS deposit;
DROP TABLE IF EXISTS withdrawal;
