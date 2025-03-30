package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type DispensingMode string

const (
	DispensingPublic  DispensingMode = "public"
	DispensingPrivate DispensingMode = "private"
	DispensingNone    DispensingMode = "none"
)

// Config struktura
type Config struct {
	Tokens []struct {
		Name      string   `yaml:"name"`
		Token     string   `yaml:"token"`
		Disallow  []string `yaml:"disallow"`
		Isolation string   `yaml:"isolation"`
	} `yaml:"tokens"`
	S3 struct {
		AccessKey string `yaml:"access_key"`
		SecretKey string `yaml:"secret_key"`
		Bucket    string `yaml:"bucket"`
		Host      string `yaml:"host"`
		Secure    bool   `yaml:"secure"`
	} `yaml:"s3"`
	Dispensing       DispensingMode `yaml:"dispensing"`
	DispensingPrefix []struct {
		Prefix string         `yaml:"prefix"`
		Mode   DispensingMode `yaml:"mode"`
	} `yaml:"dispensingPrefix"`
	Plugins []string `yaml:"plugins"`
}

var AppConfig Config

func LoadConfig() {
	// Získáme cestu k souboru z proměnné prostředí, pokud není nastavena, použijeme výchozí hodnotu
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Failed to open config file: %v", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}
}
