CREATE TABLE IF NOT EXISTS users (
	id serial NOT NULL,
	name varchar(32) NULL,
	email varchar(64) NOT NULL,
	password varchar(256) NOT NULL,
	is_active bool DEFAULT true NOT NULL,
  otp varchar(32) NULL,
	last_login timestamp NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	CONSTRAINT users_pkey PRIMARY KEY (email)
);

CREATE TABLE IF NOT EXISTS project (
		id UUID PRIMARY KEY DEFAULT default_random_uuid(),
		name varchar(64) NOT NULL,
		visibility varchar(32) DEFAULT 'public' NOT NULL,
		user_email varchar(64) REFERENCES users(email) NOT NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS monitor (
		id UUID PRIMARY KEY DEFAULT default_random_uuid(),
		name varchar(64) NOT NULL,
		period integer NOT NULL DEFAULT 600,
		grace_period integer NOT NULL DEFAULT 300,
		user_email varchar(64) REFERENCES users(email) NOT NULL,
		project_id UUID NOT NULL,
		ping_url varchar(512) NOT NULL,
		status varchar(64) NULL DEFAULT 'up', -- up, down, grace_period
		type varchar(64) NULL,
		last_ping timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL
);

CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
		NEW.updated_at = now();
		RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER updated_at BEFORE UPDATE ON monitor FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER updated_at BEFORE UPDATE ON project FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
