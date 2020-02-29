package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/enr/go-commons/environment"
	"github.com/enr/go-files/files"
	"gopkg.in/yaml.v2"
)

const (
	// ConfigFileBaseName is the default base name of config file in user home.
	ConfigFileBaseName = ".resyncrc"
)

// Configuration model.
type Configuration struct {
	Folders map[string]Folder
}

// Folder model.
type Folder struct {
	LocalPath    string `yaml:"local_path"`
	ExternalPath string `yaml:"external_path"`
}

func configurationFilePath() string {
	home, err := environment.UserHome()
	if err != nil {
		ui.Errorf("Error retrieving user home: %v\n", err)
		os.Exit(1)
	}
	configurationFile := filepath.Join(home, ConfigFileBaseName)
	ui.Confidentialf("Using configuration file %s", configurationFile)
	if !files.Exists(configurationFile) {
		ui.Errorf("Configuration file %s not found. Exit", configurationFile)
		os.Exit(configurationFileNotFoundCode)
	}
	return configurationFile
}

func loadConfigurationFromFile(configurationPath string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(configurationPath)
	if err != nil {
		ui.Errorf("Error reading %s: %v", configurationPath, err)
		return Configuration{}, err
	}
	return loadConfigurationFromBytes(bytes)
}

func loadConfigurationFromBytes(bytes []byte) (Configuration, error) {
	var conf Configuration
	err := yaml.Unmarshal(bytes, &conf)
	if err != nil {
		ui.Errorf("Error reading configuration: %v", err)
		return conf, err
	}
	return conf, nil
}
