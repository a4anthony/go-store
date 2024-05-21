package config

import (
	"database/sql"
	"fmt"
	"github.com/a4anthony/go-store/internal/database"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *database.Queries

//var DBFilter *handlers.QueriesFilter

//func NewFilter(db database.DBTX) *handlers.QueriesFilter {
//	return &handlers.QueriesFilter{Db: db}
//}

func ConnectDb() {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatalf("DB_URL environment variable not set")
	}

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		fmt.Println(err)
		panic("Failed to connect database!")
	}

	fmt.Println("Connected to database!")
	db := database.New(conn)
	DB = db
	//DBFilter = NewFilter(conn)
}
