package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	DB_Url   string `json:"db_url"`
	Username string `json:"current_user_name"`
}

const configfilename string = ".gatorconfig.json"

func getConfigFileLocation() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Can't find config location")
		return ""
	}
	return homeDir + "/" + configfilename
}

func Read() (Config, error) {
	userConfigFileLocation := getConfigFileLocation()
	dat, err := os.ReadFile(userConfigFileLocation)
	if err != nil {
		return Config{}, errors.New("can't read config file")
	}
	cfg := Config{}
	if err := json.Unmarshal(dat, &cfg); err != nil {
		return Config{}, errors.New("can't read config file")
	}
	return cfg, nil
}

// func (cfg *Config) SetUser(userName string) error {
// 	cfg.Username = userName
// 	return write(*cfg)
// }

func write(cfg Config) error {
	return nil
}

func (c Config) SetUser(username string) error {
	c.Username = username
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	// fmt.Println(string(data))
	path := getConfigFileLocation()
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

//Rewrite using json.NewEncoder then encoder.Encode
//JSON marshall vs encode
