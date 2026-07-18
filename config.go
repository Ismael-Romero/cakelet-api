package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Debug       bool              `json:"debug"`
	Server      HttpServerConfig  `json:"server"`
	MediaSource MediaSourceConfig `json:"mediaSource"`
	Database    DatabaseConfig    `json:"dataSource"`
}

type MediaSourceConfig struct {
	DirMedia string `json:"dirMedia"`
	DirLogs  string `json:"dirLogs"`
}

type HttpServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type DatabaseConfig struct {
	Engine   string `json:"engine"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

func NewConfig() (*Config, error) {
	file, err := os.Open("./config.json")
	if err != nil {
		return nil, fmt.Errorf("could not open the configuration file: %w", err)
	}
	defer file.Close()

	var cfg Config

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("Error parsing the configuration JSON\n: %w", err)
	}

	return &cfg, nil
}

func TestConfig(cfg *Config) {
	fmt.Println("--- Checking the configuration file read ---")
	fmt.Println("Debug mode:", cfg.Debug)
	fmt.Println("Database engine:", cfg.Database.Engine)
	fmt.Println("Server host:", cfg.Server.Host)
	fmt.Println("Server port:", cfg.Server.Port)
}
