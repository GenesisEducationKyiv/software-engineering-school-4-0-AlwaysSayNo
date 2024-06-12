package db

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file" // blank import needed for migration purposes

	migrPostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init is used to establish connection with database.
func Init(url string) *gorm.DB {
	db, err := gorm.Open(gormPostgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	waitDBIsUp(db)

	return db
}

// waitDBIsUp is a function that prevents application container to fail
// if docker-compose is used and db hasn't started yet.
func waitDBIsUp(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("failed to get generic database object from gorm DB:", err)
	}

	// Wait for the database to be ready
	for i := 0; i < 10; i++ {
		err = sqlDB.Ping()
		if err == nil {
			log.Println("Database connection is successful")
			break
		}

		log.Printf("Waiting for database to be ready (%d) (%v)", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatalln("failed to connect to the database:", err)
	}
}

// RunMigrations is used to run migrations against database using golang migrate.
func RunMigrations(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get SQL DB from Gorm DB: %v", err)
	}

	driver, err := migrPostgres.WithInstance(sqlDB, &migrPostgres.Config{})
	if err != nil {
		log.Fatalf("failed to create migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://pkg/common/db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("failed to run migrations: %v", err)
	}

	log.Println("Migrations ran successfully")
}

// GetDatabaseURL is used to prepare a database url.
func GetDatabaseURL() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}
