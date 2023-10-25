package database

import (
	"fmt"
	"github.com/yockii/qscore/pkg/config"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	mysqlConnStringFmt = "%s:%s@tcp(%s:%d)/%s"
	pgConnStringFmt    = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
)

func initDatabaseDefault() {
	config.DefaultInstance.SetDefault("database.driver", "mysql")
	config.DefaultInstance.SetDefault("database.host", "localhost")
	config.DefaultInstance.SetDefault("database.user", "root")
	config.DefaultInstance.SetDefault("database.password", "root")
	config.DefaultInstance.SetDefault("database.db", "householder")
	config.DefaultInstance.SetDefault("database.port", 3306)
	config.DefaultInstance.SetDefault("database.prefix", "t_")
	config.DefaultInstance.SetDefault("database.showSql", false)
}

var DB *gorm.DB

func Initial() {
	initDatabaseDefault()
	InitDB(
		config.GetString("database.driver"),
		config.GetString("database.host"),
		config.GetString("database.user"),
		config.GetString("database.password"),
		config.GetString("database.db"),
		config.GetInt("database.port"),
		config.GetString("database.prefix"),
		config.GetString("logger.level"),
	)
}

func InitDB(dbType, host, user, password, dbName string, port int, prefix string, logLevel string) {
	var err error
	if dbType == "mysql" {
		DB, err = gorm.Open(mysql.Open(fmt.Sprintf(mysqlConnStringFmt,
			user, password, host, port, dbName)), &gorm.Config{})
	} else if dbType == "pg" || dbType == "postgres" {
		DB, err = gorm.Open(postgres.Open(fmt.Sprintf(pgConnStringFmt,
			host, port, user, password, dbName)), &gorm.Config{})
	} else if dbType == "sqlite" {
		DB, err = gorm.Open(sqlite.Open(host), &gorm.Config{})
	} else {
		logrus.Fatalf("不支持的数据库: %s", dbType)
	}
	if err != nil {
		logrus.Fatalf("数据库连接失败! %v", err)
	}
	DB.Config.NamingStrategy = schema.NamingStrategy{
		TablePrefix:   prefix,
		SingularTable: true,
	}
	DB.Config.SkipDefaultTransaction = true

	slowThreshold := time.Second
	level := logger.Silent
	switch strings.ToLower(logLevel) {
	case "error":
		level = logger.Error
	case "warn":
		level = logger.Warn
	case "info":
		level = logger.Info
		slowThreshold = 0
	}

	newLogger := logger.New(logrus.StandardLogger(), logger.Config{
		SlowThreshold:             slowThreshold,
		LogLevel:                  level,
		IgnoreRecordNotFoundError: true,
	})
	DB.Config.Logger = newLogger

	//if logLevel != "" {
	//	switch strings.ToLower(logLevel) {
	//	case "error":
	//		DB.Config.Logger.LogMode(logger.Error)
	//	case "warn":
	//		DB.Config.Logger.LogMode(logger.Warn)
	//	case "info":
	//		DB.Config.Logger.LogMode(logger.Info)
	//	default:
	//		DB.Config.Logger.LogMode(logger.Silent)
	//	}
	//}
}

func Close() {
	db, _ := DB.DB()
	if db != nil {
		db.Close()
	}
}

// AutoMigrate 自动迁移
func AutoMigrate(models ...interface{}) error {
	return DB.AutoMigrate(models...)
	// 不能直接使用db的迁移，要增加表注释, 如何处理?兼容不同数据库
}
