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

-- +migrate Down
DROP TABLE IF EXISTS user;
DROP TABLE IF EXISTS bet_user;
DROP TABLE IF EXISTS identity_verification;
DROP TABLE IF EXISTS deposit;
DROP TABLE IF EXISTS withdrawal;
