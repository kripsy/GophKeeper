BEGIN TRANSACTION ;

	
	CREATE TABLE IF NOT EXISTS users
	(
		id bigint NOT NULL,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(60) NOT NULL,
		CONSTRAINT users_pkey PRIMARY KEY (id)
	);
	
	ALTER TABLE users ADD CONSTRAINT username_unq UNIQUE(username);

	CREATE INDEX users_username_key ON users USING HASH (username);


	CREATE TABLE IF NOT EXISTS secrets
	(
		id SERIAL PRIMARY KEY,
		type VARCHAR(255) NOT NULL CHECK (type IN ('text', 'binary', 'card', 'login_password')),
		data BYTEA NOT NULL,
		meta TEXT,
		chunk_num INT DEFAULT 0,
		total_chunks INT DEFAULT 0,
		user_id INT REFERENCES users(id) ON DELETE CASCADE
	);

	CREATE INDEX idx_secrets_user_id ON secrets(user_id);
	
COMMIT ;