# Anotaciones

## Migrations

go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

### Create Migration

migrate create -ext sql -dir migrations -seq create_users_table

### Apply Migration

migrate -source file://migrations -database "mysql://root:@tcp(localhost:3306)/luthier_saas_db" up
