package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFilename = ".gatorconfig.json"

// Read config from file
func Read() (*Config, error) {
	filename, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := Config{}
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Update current user and write config
func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	return c.write()
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home_dir + "/" + configFilename, nil
}

func (c *Config) write() error {
	filename, err := getConfigFilePath()
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(c); err != nil {
		return err
	}

	return nil
}
