CREATE TABLE USERS (
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
