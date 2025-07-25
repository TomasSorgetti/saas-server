package persistance

import (
	"database/sql"
	"log"
	"luthierSaas/internal/infrastructure/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Database struct {
    DB *sql.DB
}

func NewDatabase(cfg *config.Config) (*Database, error) {
    log.Printf("Attempting to connect to MySQL with URL: %s", cfg.DatabaseURL)
    db, err := sql.Open("mysql", cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("Failed to connect to MySQL: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("Failed to ping MySQL: %v", err)
    }

    driver, err := mysql.WithInstance(db, &mysql.Config{})
    if err != nil {
        log.Fatalf("Failed to create migration driver: %v", err)
    }
    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "mysql",
        driver,
    )
    if err != nil {
        log.Fatalf("Failed to initialize migrations: %v", err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Failed to apply migrations: %v", err)
    }
    log.Println("Successfully applied migrations")

    log.Println("Successfully connected to MySQL")
    return &Database{DB: db}, nil
}

func (d *Database) Close() error {
    return d.DB.Close()
}