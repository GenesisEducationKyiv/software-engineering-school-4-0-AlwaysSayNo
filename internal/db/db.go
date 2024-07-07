package db

import (
	"fmt"
	"genesis-currency-api/internal/db/config"
	"log"
	"net/url"

	_ "github.com/golang-migrate/migrate/v4/source/file" // blank import needed for migration purposes

	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Init is used to establish connection with database.
func Init(url string) *gorm.DB {
	db, err := gorm.Open(gormPostgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	return db
}

// GetDatabaseURL is used to prepare a database url.
func GetDatabaseURL(cnf config.DatabaseConfig) string {
	dbURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cnf.DBUser, cnf.DBPassword),
		Host:   fmt.Sprintf("%s:%s", cnf.DBHost, cnf.DBPort),
		Path:   cnf.DBName,
	}

	query := dbURL.Query()
	query.Add("sslmode", "disable")
	dbURL.RawQuery = query.Encode()

	return dbURL.String()
}
