package main

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// setConfig sets the config global object using the config file if present
func setConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	file := filepath.Join(homeDir, ".config", "chgo", "config.yaml")
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	return err
}

// createConfig creats the config in `.config/chgo/config.yaml` location
// if not present or overrides the previous config
func createConfig(config Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	configPath := filepath.Join(homeDir, ".config", "chgo")
	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return nil
	}

	file, err := os.Create(filepath.Join(configPath, "config.yaml"))
	if err != nil {
		return nil
	}
	defer file.Close()

	_, err = file.Write(data)

	return err
}
