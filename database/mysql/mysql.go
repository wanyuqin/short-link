package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"short-link/config"
	"time"
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
		db, err := newDb(mysqlCfg)
		if err != nil {
			panic(err)
		}
		dbMap[key] = db
	}
}

func newDb(mysqlCfg config.Mysql) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysqlCfg.User, mysqlCfg.Password, mysqlCfg.Host, mysqlCfg.Port, mysqlCfg.Dbname)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	_db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	connPool := &sql.DB{}
	connPool.SetMaxOpenConns(defaultMaxOpenConns)
	connPool.SetMaxIdleConns(defaultMaxIdleConns)
	connPool.SetConnMaxLifetime(defaultConnMaxLifetime)
	if mysqlCfg.MaxIdleConns > 0 {
		connPool.SetMaxIdleConns(mysqlCfg.MaxIdleConns)
	}
	if mysqlCfg.MaxOpenConns > 0 {
		connPool.SetMaxOpenConns(mysqlCfg.MaxOpenConns)
	}
	if mysqlCfg.ConnMaxLifetime > 0 {
		connPool.SetConnMaxLifetime(time.Duration(mysqlCfg.ConnMaxLifetime) * time.Second)
	}
	_db.ConnPool = connPool
	return _db, nil

}
