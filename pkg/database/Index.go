package database

import (
	"fmt"
	"strings"

	_ "gitee.com/chunanyong/dm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/godror/godror"

	_ "modernc.org/sqlite"

	_ "github.com/lib/pq"
	logger "github.com/sirupsen/logrus"

	"gitee.com/chunanyong/zorm"

	"github.com/yockii/qscore/pkg/config"
)

var MainDB *zorm.DBDao

func init() {
	dbType := config.GetString("database.type")
	switch strings.ToLower(dbType) {
	case "oracle":
		initOracle()
	case "dm":
		initDm()
	case "pg", "postgres", "pgsql":
		initPostgres()
	case "kingbase":
		initKingbase()
	case "mysql":
		initMysql()
	case "sqlite":
		initSqlite()
	default:
		logger.Fatal("暂未开通配置的数据库类型")
	}
}

func initSqlite() {
	var err error
	MainDB, err = zorm.NewDBDao(&zorm.DataSourceConfig{
		DSN:        config.GetString("database.address"),
		DriverName: "sqlite",
		Dialect:    "sqlite",
	})
	if err != nil {
		logger.Fatal("数据库创建失败! %v", err)
	}
}

func initKingbase() {
	var err error
	sslMode := "disabled"
	if config.IsSet("database.sslMode") {
		sslMode = config.GetString("database.sslMode")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.GetString("database.address"),
		config.GetInt("database.port"),
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.dbName"),
		sslMode,
	)
	dbConfig := &zorm.DataSourceConfig{
		DSN:           dsn,
		DriverName:    "postgres", // kingbase
		Dialect:       "postgres", // kingbase
		SlowSQLMillis: config.GetInt("database.slowSqlMillis"),
		MaxOpenConns:  config.GetInt("database.maxConn"),
		MaxIdleConns:  config.GetInt("database.maxIdle"),
	}

	MainDB, err = zorm.NewDBDao(dbConfig)
	if err != nil {
		logger.Fatal("数据库创建失败! %v", err)
	}
}

func initOracle() {
	var err error
	dsn := fmt.Sprintf(`user="%s" password="%s" connectString="%s:%d/%s"`,
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.address"),
		config.GetInt("database.port"),
		config.GetString("database.dbName"),
	)
	dbConfig := &zorm.DataSourceConfig{
		DSN:           dsn,
		DriverName:    "godror",
		Dialect:       "oracle",
		SlowSQLMillis: config.GetInt("database.slowSqlMillis"),
		MaxOpenConns:  config.GetInt("database.maxConn"),
		MaxIdleConns:  config.GetInt("database.maxIdle"),
	}

	MainDB, err = zorm.NewDBDao(dbConfig)
	if err != nil {
		logger.Fatal("数据库创建失败! %v", err)
	}
}

func initDm() {
	var err error
	dsn := fmt.Sprintf(
		"dm://%s:%s@%s:%d",
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.address"),
		config.GetInt("database.port"),
	)
	dbConfig := &zorm.DataSourceConfig{
		DSN:           dsn,
		DriverName:    "dm",
		Dialect:       "dm",
		SlowSQLMillis: config.GetInt("database.slowSqlMillis"),
		MaxOpenConns:  config.GetInt("database.maxConn"),
		MaxIdleConns:  config.GetInt("database.maxIdle"),
	}

	MainDB, err = zorm.NewDBDao(dbConfig)
	if err != nil {
		logger.Fatal("数据库创建失败! %v", err)
	}
}

func initPostgres() {
	var err error
	sslMode := "disabled"
	if config.IsSet("database.sslMode") {
		sslMode = config.GetString("database.sslMode")
	}
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.GetString("database.address"),
		config.GetInt("database.port"),
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.dbName"),
		sslMode,
	)
	dbConfig := &zorm.DataSourceConfig{
		DSN:           dsn,
		DriverName:    "postgres",
		Dialect:       "postgresql",
		SlowSQLMillis: config.GetInt("database.slowSqlMillis"),
		MaxOpenConns:  config.GetInt("database.maxConn"),
		MaxIdleConns:  config.GetInt("database.maxIdle"),
	}

	MainDB, err = zorm.NewDBDao(dbConfig)
	if err != nil {
		logger.Fatal("数据库创建失败! %v", err)
	}
}

func initMysql() {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		config.GetString("database.username"),
		config.GetString("database.password"),
		config.GetString("database.address"),
		config.GetInt("database.port"),
		config.GetString("database.dbName"),
	)
	if config.IsSet("database.charset") {
		dsn += "&charset=" + config.GetString("database.charset")
	}

	dbConfig := &zorm.DataSourceConfig{
		DSN:           dsn,
		DriverName:    "mysql",
		Dialect:       "mysql",
		SlowSQLMillis: config.GetInt("database.slowSqlMillis"),
		MaxOpenConns:  config.GetInt("database.maxConn"),
		MaxIdleConns:  config.GetInt("database.maxIdle"),
	}

	MainDB, err = zorm.NewDBDao(dbConfig)
	if err != nil {
		logger.Fatal("数据库创建失败! %v", err)
	}
}
