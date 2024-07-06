package config

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
	DB       int    `yaml:"db"`
	Password string `yaml:"password"`
}

func InitializeConfig(path string) error {
	var err error
	once.Do(func() {
		cfgManager := newConfigManager()
		err = cfgManager.readPath(path).initConfigWithViper().watchConfig().error()
	})
	return err
}

func GetConfig() *Config {
	return cfgManager.getConfig()
}

func (cfg *Config) GetHTTPConfig(key string) HttpServer {
	return cfg.Application.Http[key]
}
