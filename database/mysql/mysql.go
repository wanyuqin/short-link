package mysql

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"short-link/config"
	"short-link/utils/valuex"
)

var (
	defaultMaxOpenConns    = 100
	defaultMaxIdleConns    = 10
	defaultConnMaxLifetime = time.Hour
	defaultDbName          = "default"

	dbMap = make(map[string]*gorm.DB)
)

func NewDBClient(ctx context.Context, key ...string) *gorm.DB {
	dbName := defaultDbName
	if dbMap == nil || len(dbMap) == 0 {
		return nil
	}
	if len(key) > 0 {
		dbName = key[0]
	}
	if db, ok := dbMap[dbName]; ok {
		return db.WithContext(ctx)
	}
	return nil
}

func InitializeDBClient() {
	cfg := config.GetConfig()
	for key, mysqlCfg := range cfg.Database.Mysql {
		db, err := newDB(&mysqlCfg)
		if err != nil {
			panic(err)
		}
		dbMap[key] = db
	}
}

func newDB(mysqlCfg *config.Mysql) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlCfg.User, mysqlCfg.Password, mysqlCfg.Host, mysqlCfg.Port, mysqlCfg.Dbname)

	newLogger := createLogger()
	_db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := _db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(valuex.GetOrDefault(mysqlCfg.MaxOpenConns, defaultMaxOpenConns))
	sqlDB.SetMaxIdleConns(valuex.GetOrDefault(mysqlCfg.MaxIdleConns, defaultMaxIdleConns))
	sqlDB.SetConnMaxLifetime(valuex.GetDurationOrDefault(mysqlCfg.ConnMaxLifetime, defaultConnMaxLifetime))
	return _db, nil
}

func createLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,        // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
}
