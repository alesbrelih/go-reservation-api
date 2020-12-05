CREATE TABLE IF NOT EXISTS "inquiry" (
	id bigserial primary key,
	inquirer varchar(255) NOT NULL,
	email varchar(100),
	phone varchar(25),
	date_reservation timestamp NOT NULL,
	date_created timestamp NOT NULL
);