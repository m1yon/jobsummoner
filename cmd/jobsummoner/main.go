package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/m1yon/jobsummoner/internal/database"
	_ "modernc.org/sqlite"
)

func main() {
	godotenv.Load()
	dbConnection := os.Getenv("DB_CONNECTION")
	db, err := sql.Open("sqlite", dbConnection)

	if err != nil {
		fmt.Println("err:", err.Error())
		return
	}

	dbQueries := database.New(db)
	err = dbQueries.CreateCompany(context.Background(), database.CreateCompanyParams{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: "test", Url: "google.com"})

	if err != nil {
		fmt.Println("err:", dbConnection, err.Error())
		return
	}

}
