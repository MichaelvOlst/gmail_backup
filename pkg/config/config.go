package config

import (
	"errors"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/gommon/log"
)

// Config wraps the configuration structs for the various application parts
type Config struct {
	Server struct {
		Host   string `default:"127.0.0.1"`
		Port   string `default:"8080"`
		Secret string
	}
	Database struct {
		Filename string `default:"gmail_backup.db"`
	}
	Google struct {
		File string `default:"credentials.json"`
	}
}

// Load loads the config file
func Load(file string) error {

	if file == "" {
		// log.Warn("Using default config")
		return nil
	}

	absFile, _ := filepath.Abs(file)
	_, err := os.Stat(absFile)
	fileNotExists := os.IsNotExist(err)

	if fileNotExists {
		return errors.New("Error reading configuration. File " + file + " does not exist.")
	}

	log.Printf("Configuration file: %s", absFile)

	// read file into env values
	err = godotenv.Load(absFile)
	if err != nil {
		return err
	}

	return nil
}

// Parse handles the config
func Parse() (*Config, error) {
	var cfg Config

	// with config file loaded into env values, we can now parse env into our config struct
	err := envconfig.Process("app", &cfg)
	if err != nil {
		return nil, err
	}

	// // alias sqlite to sqlite3
	// if cfg.Database.Driver == "sqlite" {
	// 	cfg.Database.Driver = "sqlite3"
	// }

	// // use absolute path to sqlite3 database
	// if cfg.Database.Driver == "sqlite3" {
	// 	cfg.Database.Name, _ = filepath.Abs(cfg.Database.Name)
	// }

	// if secret key is empty, use a randomly generated one
	if cfg.Server.Secret == "" {
		cfg.Server.Secret = randomString(40)
	}

	return &cfg, nil
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}

	return string(bytes)
}
