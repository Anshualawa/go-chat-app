package config

import (
	"log"

	"github.com/go-redis/redis/v7"
)

var RDB *redis.Client

func ConnectRedis() {
	// Ensure Config id Loaded
	if ConfigData.RedisAddr == "" {
		LoadConfig() // Load config if not already loaded
	}

	// Initialize Redis Client
	RDB = redis.NewClient(&redis.Options{
		Addr: ConfigData.RedisAddr,
	})

	// Ping Redis to test connection
	_, err := RDB.Ping().Result()
	if err != nil {
		log.Fatal("❌ Redis connection failed :", err)
	}
	log.Println("✅ Connect to Redis")
}
