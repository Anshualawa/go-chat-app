package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// Config struct to hold database and redis config
type Config struct {
	DBUser    string `json:"DB_USER"`
	DBPass    string `json:"DB_PASS"`
	DBHost    string `json:"DB_HOST"`
	DBPort    string `json:"DB_PORT"`
	DBName    string `json:"DB_NAME"`
	RedisAddr string `json:"REDIS_ADDR"`
}

var (
	DB         *sql.DB
	ConfigData Config
)

// LoadConfig reads config.json and parses it
func LoadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to open config file:", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&ConfigData)
	if err != nil {
		log.Fatal("Failed to parse config file:", err)
	}
	fmt.Println("Config Loaded Successfully:", ConfigData)
}

func ConnectDB() {
	// Load configration
	LoadConfig()

	// ✅ Creating DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		ConfigData.DBUser, ConfigData.DBPass,
		ConfigData.DBHost, ConfigData.DBPort,
		ConfigData.DBName)

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
