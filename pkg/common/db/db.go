package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"os/exec"
	"time"
)

func Init(url string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})

	if err != nil {
		log.Fatalln(err)
	}
	waitDbIsUp(db)

	return db
}

func waitDbIsUp(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalln("Failed to get generic database object from gorm DB:", err)
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
		log.Fatalln("Failed to connect to the database:", err)
	}
}

func RunMigrations(dbUrl string) {
	cmd := exec.Command("migrate", "-path", "pkg/common/db/migrations", "-database", dbUrl+"?sslmode=disable", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func GetUrl() string {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}
func WaitDbIsUp() {

}
