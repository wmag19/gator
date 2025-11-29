package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DB_Url   string `json:"db_url"`
	Username string `json:"current_user_name"`
}

const configfilename string = ".gatorconfig.json"

func getConfigFileLocation() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + configfilename, nil
}

func Read() (Config, error) {
	userConfigFileLocation, err := getConfigFileLocation()
	if err != nil {
		return Config{}, err
	}
	dat, err := os.ReadFile(userConfigFileLocation)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	if err := json.Unmarshal(dat, &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFileLocation()
	if err != nil {
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)

	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config) SetUser(userName string) error {
	cfg.Username = userName
	return write(*cfg)
}
