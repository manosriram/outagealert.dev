CREATE TABLE plan (
		id serial,
		name varchar(64) UNIQUE,
		price int,
		validity int,
		PRIMARY KEY (id)
);

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
  plan varchar(64) REFERENCES plan(name) DEFAULT 'free',
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
		status_before_pause varchar(64) NULL, -- status before pause inorder to restore when resuming
		is_active boolean DEFAULT true,
		type varchar(64) NOT NULL,
		total_pause_time integer DEFAULT 0,
		last_ping timestamp DEFAULT NULL,
		last_paused_at timestamp,
		last_resumed_at timestamp,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL
);

CREATE TABLE IF NOT EXISTS ping (
		id varchar(22) PRIMARY KEY,
		monitor_id varchar(22) REFERENCES monitor(id) NOT NULL,
		status integer,
		metadata jsonb,
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

CREATE TYPE ALERT_TYPE as ENUM ('email', 'slack', 'webhook');
CREATE TABLE IF NOT EXISTS alert_integration (
		id varchar(22),
		monitor_id varchar(22) REFERENCES monitor(id) NOT NULL,
		is_active boolean DEFAULT false NOT NULL,
		email_alert_sent boolean DEFAULT false NOT NULL,
		slack_alert_sent boolean DEFAULT false NOT NULL,
		webhook_alert_sent boolean DEFAULT false NOT NULL,
		alert_type ALERT_TYPE NOT NULL,
		alert_target varchar(512),
		created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		primary key(monitor_id, alert_type)
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
CREATE TRIGGER updated_at BEFORE UPDATE ON alert_integration FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
