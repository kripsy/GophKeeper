BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS users
(
	id bigint NOT NULL,
	username VARCHAR(255) NOT NULL,
	password VARCHAR(60) NOT NULL,
	PRIMARY KEY (id),
	UNIQUE(username)
);

COMMIT;
