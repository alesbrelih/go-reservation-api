CREATE TABLE IF NOT EXISTS "user" (
	id bigserial primary key,
	first_name varchar(255),
	last_name varchar(255),
	username varchar(255) UNIQUE,
	email varchar(255) UNIQUE,
	pass text
);