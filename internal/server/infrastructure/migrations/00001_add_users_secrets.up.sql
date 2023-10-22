BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users
(
	id bigint NOT NULL,
	username VARCHAR(255) NOT NULL,
	password VARCHAR(60) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE(username)
);

CREATE TABLE IF NOT EXISTS secrets
(
	id SERIAL PRIMARY KEY,
	external_id UUID NOT NULL UNIQUE,
	user_id INT REFERENCES users(id) ON DELETE CASCADE,
	lastUpdate TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	hash CHAR(64) NOT NULL
);

CREATE INDEX idx_secrets_user_id ON secrets(user_id);

COMMIT;
