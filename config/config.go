package config

import (
	"encoding/json"
	"fmt"
	"github.com/prasmussen/gdrive/util"
	"io/ioutil"
)

// Client ID and secrect for installed applications
const (
	ClientId     = "367116221053-7n0vf5akeru7on6o2fjinrecpdoe99eg.apps.googleusercontent.com"
	ClientSecret = "1qsNodXNaWq1mQuBjUjmvhoO"
)

type Config struct {
	ClientId     string
	ClientSecret string
}

func defaultConfig() *Config {
	return &Config{
		ClientId:     ClientId,
		ClientSecret: ClientSecret,
	}
}

func promptUser() *Config {
	return &Config{
		ClientId:     util.Prompt("Enter Client Id: "),
		ClientSecret: util.Prompt("Enter Client Secret: "),
	}
}

func load(fname string) (*Config, error) {
	data, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	return config, json.Unmarshal(data, config)
}

func save(fname string, config *Config) error {
	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	if err = util.Mkdir(fname); err != nil {
		return err
	}
	return ioutil.WriteFile(fname, data, 0600)
}

func Load(fname string, advancedUser bool) *Config {
	config, err := load(fname)
	if err != nil {
		// Unable to read existing config, lets start from scracth
		// Get config from user input for advanced users, or just use default settings
		if advancedUser {
			config = promptUser()
		} else {
			config = defaultConfig()
		}

		// Save config to file
		err := save(fname, config)
		if err != nil {
			fmt.Printf("Failed to save config (%s)\n", err)
		}
	}
	return config
}
