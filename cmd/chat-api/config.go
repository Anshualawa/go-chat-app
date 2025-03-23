package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var (
	Cfg *Config
)

// Config holds all the configration parameters of the service.
type Config struct {
	Address                string `json:"address"`
	Port                   uint16 `json:"port"`
	ShutdownTimeoutSeconds uint8  `json:"shutdownTimeoutSeconds"`
	MaxReqPerMinute        uint32 `json:"maxReqPerMinute"`
	ReqTimeoutSeconds      uint8  `json:"reqTimeoutSeconds"`
}

// LoadConfigS reads config.json and parses it
func LoadConfigS() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to open config file:", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Cfg)
	if err != nil {
		log.Fatal("Failed to parse config file:", err)
	}
	fmt.Println("Config Loaded Successfully:", Cfg)
}
