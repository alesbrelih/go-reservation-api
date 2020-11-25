CREATE TABLE "item" (
	"id" bigserial PRIMARY KEY,
	"title" varchar(255) NOT NULL,
	"show_from" timestamp,
	"show_to" timestamp
);