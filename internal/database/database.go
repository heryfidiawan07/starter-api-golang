package database

import (
	"fmt"
	"log"

	"starter-api-golang/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	gormCfg := &gorm.Config{}
	if cfg.App.Env == "development" {
		gormCfg.Logger = logger.Default.LogMode(logger.Info)
	} else {
		gormCfg.Logger = logger.Default.LogMode(logger.Error)
	}

	var db *gorm.DB
	var err error

	switch cfg.DB.Driver {
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
		db, err = gorm.Open(mysql.Open(dsn), gormCfg)

	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			cfg.DB.Host, cfg.DB.User, cfg.DB.Pass, cfg.DB.Name, cfg.DB.Port, cfg.DB.SSLMode)
		db, err = gorm.Open(postgres.Open(dsn), gormCfg)

	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
			cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
		db, err = gorm.Open(sqlserver.Open(dsn), gormCfg)

	case "sqlite":
		db, err = gorm.Open(sqlite.Open(cfg.DB.Name+".db"), gormCfg)

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.DB.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Printf("Connected to %s database", cfg.DB.Driver)
	return db, nil
}
