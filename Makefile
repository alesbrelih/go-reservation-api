DBUSER=${DB_USER}
DBPASS=${DB_PASS}
DBHOST=${DB_HOST}
DBPORT=${DB_PORT}
DBNAME=${DB_NAME}

ifndef DBUSER
$(error DB_USER enviroment variable is not set)
endif

ifndef DBPASS
$(error DB_PASS enviroment variable is not set)
endif

ifndef DBHOST
$(error DB_HOST enviroment variable is not set)
endif

ifndef DBPORT
$(error DB_PORT enviroment variable is not set)
endif

ifndef DBNAME
$(error DB_NAME enviroment variable is not set)
endif

composeup:
	docker-compose -f compose/docker-compose.yml up

composedown:
	docker-compose -f compose/docker-compose.yml down

migratecreate:
	@echo "Enter migration name";
	@read MIGRATION; migrate create -ext sql -dir db/migrations -seq $$MIGRATION

migrateup:
	migrate -path db/migrations -database "postgresql://${DBUSER}:${DBPASS}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://${DBUSER}:${DBPASS}@${DBHOST}:${DBPORT}/${DBNAME}?sslmode=disable" -verbose down

.PHONY: composeup composedown migratecreate migrateup migratedown