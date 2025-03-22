package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	// 🚀 Hardcoded Database Credentials
	dbUser := "root"
	dbPass := "system"
	dbHost := "127.0.0.1"
	dbPort := "3306"
	dbName := "social_chat"

	// ✅ Creating DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	fmt.Println("🔍 Connecting to DB:", dsn)

	// 🚀 Connecting to MySQL
	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Database Connection Failed:", err)
	}

	// 🔄 Checking Connection
	if err = DB.Ping(); err != nil {
		log.Fatal("❌ Database Ping Failed:", err)
	}

	log.Println("✅ Connected to MySQL Database")
}
