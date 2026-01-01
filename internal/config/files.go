package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

func Read() Config {
	filePath, err := getConfigFilePath()
	if err != nil {
		log.Fatal(err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	err = os.WriteFile(filePath, data, 0666)
	if err != nil {
		return err
	}
	return nil
}

func getConfigFilePath() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	configPath := filepath.Join(homedir, configFileName)
	return configPath, nil
}
