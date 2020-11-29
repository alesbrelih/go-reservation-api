# Requisites

Using golang-migrate for db migrations. Link: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

# Enviroment variables

## Docker compose

Database (compose/.postgres.env):
* POSTGRES_USER (user to be created on db and database with same name aswell)
* POSTGRES_PASSWORD (user password)

More info/configuration: [Postgres Docker Hub](https://hub.docker.com/_/postgres) 


## Application
* APPLICATION_PORT
* DB_HOST
* DB_PORT
* DB_USER
* DB_PASS


# Migrations (using CLI)

### Create migration
migrate create -ext sql -dir db/migrations -seq my-new-migration

### Up migration
migrate -path db/migrations -database "postgresql://myuser:mypwd@dbhost:dbpass/db?sslmode=disable" -verbose up

Ssl disable is since we are using local env

### Down migration
migrate -path db/migrations -database "postgresql://myuser:mypwd@dbhost:dbpass/db?sslmode=disable" -verbose up

# Makefile

Contains commands:
* make composeup
* make composedown
* make migratecreate (creates new migration with specified name)
* make migrateup
* make migratedown

It needs enviroment variables set to work (migrations part). 
There is another script (./setenv.sh) which sets variables for current shell using source ./setenv.sh, if you dont want to set it to local machine
You need to set variables inside file though.


