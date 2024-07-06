package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	cfgManager *configManager
	once       sync.Once
)

var (
	defaultConfigName = "config"
	defaultConfigType = "yaml"
)

type configManager struct {
	cfgName string
	cfgType string
	cfgPath string

	cfg *Config

	lock sync.Mutex

	err error
}

func newConfigManager() *configManager {
	cfgManager = &configManager{
		cfgType: defaultConfigType,
		cfgName: defaultConfigName,
		cfgPath: getLocalConfigPath(),
	}
	return cfgManager
}

func (m *configManager) readPath(path string) *configManager {
	if path != "" {
		dir, file := filepath.Split(path)
		ext := filepath.Ext(file)
		m.cfgName = strings.TrimSuffix(file, ext)
		m.cfgType = strings.TrimPrefix(ext, ".")
		m.cfgPath = dir
	}
	return m
}

func (m *configManager) initConfigWithViper() *configManager {
	viper.SetConfigName(m.cfgName)
	viper.SetConfigType(m.cfgType)
	viper.AddConfigPath(m.cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		m.err = err
		return m
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		m.err = err
		return m
	}
	m.cfg = &config

	return m
}

func (m *configManager) watchConfig() *configManager {
	if m.err != nil {
		return m
	}
	viper.OnConfigChange(func(e fsnotify.Event) {
		switch e.Op {
		case fsnotify.Write:
			m.reloadCfg()
		}
	})
	viper.WatchConfig()
	return m
}

func (m *configManager) error() error {
	return m.err

}

func getLocalConfigPath() string {
	_, currentFile, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(currentFile)
	cfgPath := filepath.Join(basePath, "local")
	return cfgPath
}

func (m *configManager) reloadCfg() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.initConfigWithViper()
	fmt.Printf("%#v", m.cfg)
}

func (m *configManager) getConfig() *Config {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.cfg

}
