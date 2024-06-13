package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

var (
	_config *Config
	once    sync.Once
)

type Config struct {
	Application Application `yaml:"application"`
}

type Application struct {
	ContextPath string `yaml:"contextPath"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Mode        string `yaml:"mode"`
}

func Initialize(path string) {
	once.Do(func() {
		_config = &Config{}
	})
	if path == "" {
		getLocalConfig()
	}
}

func GetConfig() *Config {
	if _config == nil {
		return &Config{}
	}
	return _config
}
func getLocalConfig() {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(currentFile)
	cfgPath := filepath.Join(basePath, "local/config.yaml")
	file, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(file, _config)
	if err != nil {
		panic(err)
	}
}
