CREATE TABLE IF NOT EXISTS "tenant_has_reservation_user" (
	reservation_user_id bigint REFERENCES reservation_user(id) ON UPDATE CASCADE ON DELETE CASCADE,
	tenant_id bigint REFERENCES tenant(id) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT ruht_pk PRIMARY KEY (reservation_user_id, tenant_id)
);