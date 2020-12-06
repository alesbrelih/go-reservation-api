CREATE TABLE IF NOT EXISTS "accepted" (
	id bigserial primary key,
	inquirer varchar(255) NOT NULL,
	inquirer_email varchar(100),
	inquirer_phone varchar(25),
	inquirer_comment text,
	item_id bigint REFERENCES item(id) ON UPDATE CASCADE ON DELETE CASCADE,
	item_title varchar(255),
	item_price bigint,
	notes text,
	date_reservation timestamp NOT NULL,
	date_inquiry_created timestamp NOT NULL,
	date_accepted timestamp NOT NULL
);