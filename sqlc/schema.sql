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
	is_verified boolean NOT NULL,
	password varchar(256) NOT NULL,
	is_active bool DEFAULT true NOT NULL,
  otp varchar(32) NULL,
  magic_token varchar(22) NULL,
	last_login timestamp NULL,
	created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
	updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
  plan varchar(64) REFERENCES plan(name) DEFAULT 'free',
	CONSTRAINT users_pkey PRIMARY KEY (email),
		uuid varchar(22) NOT NULL
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
		period integer NOT NULL DEFAULT 25,
		period_text varchar(16) NOT NULL DEFAULT 'minutes',
		grace_period_text varchar(16) NOT NULL DEFAULT 'minutes',
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

CREATE TABLE slack_users (
		monitor_id varchar(22) REFERENCES monitor(id) NOT NULL,
		channel_url varchar(512),
		channel_id varchar(64),
		channel_name varchar(256),
		configuration_url varchar(512),
		created_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		updated_at timestamp DEFAULT CURRENT_TIMESTAMP NULL,
		primary key(monitor_id)
);


CREATE TABLE contact_us (
		name varchar(164),
		email varchar(64),
		message text,
		created_at timestamp DEFAULT CURRENT_TIMESTAMP,
		updated_at timestamp default CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_orders (
		order_id varchar(256) PRIMARY KEY,
		user_email varchar(64) REFERENCES users(email) NOT NULL,
		order_status varchar(16) NOT NULL,
		order_payment_session_id varchar(256),
		plan varchar(22),
		order_expiry_time timestamp,
		order_currency varchar(16),
		order_amount integer,
		order_metadata json,
		recharge_date timestamp,
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
CREATE TRIGGER updated_at BEFORE UPDATE ON ping FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER updated_at BEFORE UPDATE ON event FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER updated_at BEFORE UPDATE ON alert_integration FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER updated_at BEFORE UPDATE ON contact_us FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER updated_at BEFORE UPDATE ON user_orders FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
CREATE TRIGGER updated_at BEFORE UPDATE ON slack_users FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
