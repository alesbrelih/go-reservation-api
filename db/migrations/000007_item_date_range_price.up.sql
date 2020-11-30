CREATE TABLE IF NOT EXISTS "item_date_range_price" (
	id bigserial primary key,
	item_id bigint references item(id) ON UPDATE CASCADE ON DELETE CASCADE,
	date_from timestamp NOT NULL,
	date_to timestamp NOT NULL,
	price bigint NOT NULL
);