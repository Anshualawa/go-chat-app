package config

import (
	"log"
	"os"

	"github.com/go-redis/redis/v7"
)

var RDB *redis.Client

func ConnectRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"),
	})

	_, err := RDB.Ping().Result()
	if err != nil {
		log.Fatal("❌ Redis connection failed :", err)
	}
	log.Println("✅ Connect to Redis")
}
