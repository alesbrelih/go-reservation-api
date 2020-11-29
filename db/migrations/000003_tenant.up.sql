CREATE TABLE IF NOT EXISTS "tenant" (
	id bigserial primary key,
	title varchar(255) NOT NULL,
	email varchar(255) NOT NULL
);