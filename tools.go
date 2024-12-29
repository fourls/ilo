package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type ToolInfo struct {
	Name string
	Path string
}

type Toolbox struct {
	Tools map[string]ToolInfo
}

func getConfigPath() (string, error) {
	var userConfig, err = os.UserConfigDir()
	if err != nil {
		return "", err
	}

	var configPath = filepath.Join(userConfig, "ilo")
	return configPath, os.MkdirAll(configPath, os.ModePerm)
}

func GetToolbox() (*Toolbox, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return &Toolbox{}, err
	}

	var configFilePath = filepath.Join(configPath, "toolbox.json")

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return &Toolbox{}, err
	}

	var toolbox Toolbox
	err = json.Unmarshal(data, &toolbox)
	if err != nil {
		return &Toolbox{}, err
	}

	return &toolbox, nil
}
