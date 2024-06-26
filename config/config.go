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
	Database    Database    `yaml:"database"`
}

type Application struct {
	Http   map[string]HttpServer `yaml:"http"`
	Logger Logger                `yaml:"logger"`
}

type Logger struct {
	Level    string `yaml:"level"`
	StdType  string `yaml:"stdType"`
	FilePath string `yaml:"filePath"`
}

type HttpServer struct {
	Host        string
	Port        int
	Mode        string
	ContextPath string
}

type Database struct {
	Mysql map[string]Mysql `yaml:"mysql"`
	Redis map[string]Redis `yaml:"redis"`
}

type Mysql struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
	Dbname          string `yaml:"dbname"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	ConnMaxLifetime int64  `yaml:"connMaxLifetime"`
}

type Redis struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       int    `yaml:"DB"`
	Password string `yaml:"password"`
}

func InitializeConfig(path string) {
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

func (cfg *Config) GetHttpConfig(key string) HttpServer {
	return cfg.Application.Http[key]
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
