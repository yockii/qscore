package database

import (
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"xorm.io/xorm"
	"xorm.io/xorm/log"
	"xorm.io/xorm/names"

	"github.com/yockii/qscore/pkg/config"
	"github.com/yockii/qscore/pkg/logger"
)

const (
	mysqlConnStringFmt = "%s:%s@tcp(%s:%d)/%s"
	pgConnStringFmt    = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
)

var DB *xorm.Engine

func InitDB(dbType, host, user, password, dbName string, port int) {
	var err error
	DB, err = initDB(dbType, host, user, password, dbName, port)
	if err != nil {
		logger.Fatalf("数据库连接失败! %v", err)
	}
	if err = DB.Ping(); err != nil {
		logger.Fatalf("数据库连接失败! %v", err)
	}
	if config.IsSet("database.prefix") {
		DB.SetTableMapper(names.NewPrefixMapper(names.SnakeMapper{}, config.GetString("database.prefix")))
	}

	if config.GetBool("database.showSql") {
		DB.ShowSQL(true)
	}
	if config.IsSet("log.level") {
		switch strings.ToLower(config.GetString("log.level")) {
		case "error":
			DB.SetLogLevel(log.LOG_ERR)
		case "warn":
			DB.SetLogLevel(log.LOG_WARNING)
		case "info":
			DB.SetLogLevel(log.LOG_INFO)
		case "debug":
			DB.SetLogLevel(log.LOG_DEBUG)
		default:
			DB.SetLogLevel(log.LOG_OFF)
		}
	}

}

func initDB(dbType string, host string, user string, password string, dbName string, port int) (*xorm.Engine, error) {
	if dbType == "mysql" {
		return xorm.NewEngine("mysql", fmt.Sprintf(
			mysqlConnStringFmt,
			user,
			password,
			host,
			port,
			dbName,
		))
	} else if dbType == "pg" || dbType == "postgres" {
		return xorm.NewEngine("postgres", fmt.Sprintf(
			pgConnStringFmt,
			host,
			port,
			user,
			password,
			dbName,
		))
	} else {
		logger.Errorf("数据库初始化失败, 不支持的数据库类型! type=%s, host=%s, user=%s, pwd=%s, db=%s, port=%d", dbType, host, user, password, dbName, port)
		return nil, errors.New("数据库初始化失败, 不支持的数据库类型")
	}
}

func CloseDB() {
	_ = DB.Close()
}
