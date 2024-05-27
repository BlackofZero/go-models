package mysql

import (
	"io"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func initMySQLConfig(dsn string) mysql.Config {
	return mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         128,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}
}

func initGormConfig(logLevel logger.LogLevel, logFile io.Writer) *gorm.Config {
	return &gorm.Config{
		Logger: logger.New(
			log.New(logFile, "\n", log.LstdFlags),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: false,
				Colorful:                  true,
			},
		),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt:     true,
		CreateBatchSize: 100,
		QueryFields:     true,
	}
}
