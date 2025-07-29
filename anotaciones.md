# Anotaciones

## Migrations

go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

### Create Migration

migrate create -ext sql -dir migrations -seq create_users_table

### Apply Migration

migrate -source /app/migrations -database "mysql://root:@tcp(localhost:3306)/luthier_saas_db" up

## Docker

docker exec -it backend-backend-1 sh
docker exec -it backend-redis-1 sh
docker exec -it backend-mysql-1 sh
