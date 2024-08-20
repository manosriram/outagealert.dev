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
		id varchar(22) PRIMARY KEY,
		name varchar(64) NOT NULL,
		visibility varchar(32) DEFAULT 'public' NOT NULL,
		user_email varchar(64) REFERENCES users(email) NOT NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS monitor (
		id varchar(22) PRIMARY KEY,
		name varchar(64) NOT NULL,
		period integer NOT NULL DEFAULT 600,
		grace_period integer NOT NULL DEFAULT 300,
		user_email varchar(64) REFERENCES users(email) NOT NULL,
		project_id varchar(22) REFERENCES project(id) NOT NULL,
		ping_url varchar(512) NOT NULL,
		status varchar(64) NOT NULL DEFAULT 'up', -- up, down, grace_period, paused
		type varchar(64) NOT NULL,
		last_ping timestamp,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS ping (
		id varchar(22) PRIMARY KEY,
		monitor_id varchar(22) REFERENCES monitor(id) NOT NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS event (
		id varchar(22) primary key,
		monitor_id varchar(22) references monitor(id) NOT NULL,
		from_status varchar(64) NOT NULL,
		to_status varchar(64) NOT NULL,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp default current_timestamp
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
CREATE TRIGGER updated_at BEFORE UPDATE ON ping FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER updated_at BEFORE UPDATE ON event FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
