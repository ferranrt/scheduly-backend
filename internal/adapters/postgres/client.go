package postgres

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GormPostgreSQLConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func getDSN(cfg GormPostgreSQLConfig) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.DBName,
		cfg.SSLMode,
	)
}

func getGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Info),
		PrepareStmt:            true,
		AllowGlobalUpdate:      false,
		SkipDefaultTransaction: true,
	}
}

func NewGormPostgreSQL(cfg GormPostgreSQLConfig) (*gorm.DB, error) {
	dsn := getDSN(cfg)
	gormCfg := getGormConfig()
	db, err := gorm.Open(postgres.Open(dsn), gormCfg)
	if err != nil {
		return nil, err
	}

	return db, nil
}
