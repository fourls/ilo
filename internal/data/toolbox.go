package data

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hairyhenderson/go-which"
)

type ToolInfo struct {
	Name string
	Path string
}

type Toolbox struct {
	path  string
	tools map[string]ToolInfo
}

type jsonToolbox struct {
	Tools map[string]ToolInfo
}

func (t *Toolbox) Load() error {
	bytes, err := os.ReadFile(t.path)
	if err != nil {
		return err
	}

	var data jsonToolbox
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return err
	}

	t.tools = data.Tools
	return nil
}

func (t *Toolbox) Save() error {
	data := jsonToolbox{
		Tools: t.tools,
	}

	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(t.path, bytes, os.ModePerm)
}

func (t *Toolbox) Get(name string) (ToolInfo, bool) {
	val, ok := t.tools[name]
	return val, ok
}

func (t *Toolbox) AddAuto(name string) (ToolInfo, error) {
	programName := name

	if runtime.GOOS == "windows" && filepath.Ext(name) == "" {
		programName = name + ".exe"
	}

	path := which.Which(programName)
	if path == "" {
		return ToolInfo{}, fmt.Errorf("add tool '%s': could not find on PATH", name)
	}

	return t.AddManual(name, path), nil
}

func (t *Toolbox) AddManual(name string, path string) ToolInfo {
	t.tools[name] = ToolInfo{
		Name: name,
		Path: path,
	}
	return t.tools[name]
}

func (t *Toolbox) Remove(name string) bool {
	_, exists := t.tools[name]
	if exists {
		delete(t.tools, name)
	}
	return exists
}

func getConfigPath() (string, error) {
	var userConfig, err = os.UserConfigDir()
	if err != nil {
		return "", err
	}

	var configPath = filepath.Join(userConfig, "ilo")
	return configPath, os.MkdirAll(configPath, os.ModePerm)
}

func getConfigFilePath() (string, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return "", err
	}

	return filepath.Join(configPath, "toolbox.json"), nil
}

func NewToolbox(path string) (Toolbox, error) {
	var toolbox = Toolbox{path: path, tools: make(map[string]ToolInfo)}

	err := toolbox.Load()
	if errors.Is(err, os.ErrNotExist) {
		// the file not existing is fine, we can just create it
		return toolbox, toolbox.Save()
	} else {
		return toolbox, err
	}
}

func NewProdToolbox() (*Toolbox, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	toolbox, err := NewToolbox(configFilePath)
	if err != nil {
		return nil, err
	} else {
		return &toolbox, nil
	}
}
