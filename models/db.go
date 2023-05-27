package models

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql" // indirect
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
)

const (
	DRIVER_MY_SQL      = "mysql"
	DRIVER_POSTGRE_SQL = "postgres"
	DRIVER_SQLITE3     = "sqlite3"
	DRIVER_SQL_SERVER  = "mssql"
)

var DB *gorm.DB

// ConnectDB 连接数据库
func ConnectDB() {
	driver := os.Getenv("DB_DRIVER")
	fmt.Printf("数据库类型：%s\n", driver)
	var dialector gorm.Dialector
	switch driver {
	case DRIVER_MY_SQL:
		dialector = ConnectDbMySQL(
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_DATABASE"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_CHARSET"),
		)
	case DRIVER_POSTGRE_SQL:
		dialector = ConnectDbPostgreSQL(
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_DATABASE"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
		)
	case DRIVER_SQLITE3:
		dialector = ConnectDbSqlite3(
			os.Getenv("DB_HOST"),
		)
	case DRIVER_SQL_SERVER:
		dialector = ConnectDbMySQL(
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_DATABASE"),
			os.Getenv("DB_USERNAME"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_CHARSET"),
		)
	default:
		log.Fatalf("models.ConnectDB driver err: %s", driver)
	}
	// 日志记录
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,   // 慢 SQL 阈值
			LogLevel:      logger.Silent, // Log level
			Colorful:      false,         // 禁用彩色打印
		},
	)

	db, err := gorm.Open(dialector, &gorm.Config{
		// 跳过默认事务 为了确保数据一致性，GORM 会在事务里执行写入操作（创建、更新、删除）。如果没有这方面的要求，您可以在初始化时禁用它。
		// SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "t_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
		Logger: newLogger,
	})

	if err != nil {
		log.Fatalln("db err：", err)
		return
	}

	maxConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	openConnections, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNECTIONS"))

	db.Use(
		dbresolver.Register(dbresolver.Config{}).
			// SetConnMaxIdleTime(time.Hour).
			// SetConnMaxLifetime(24 * time.Hour).
			SetMaxIdleConns(openConnections).
			SetMaxOpenConns(maxConnections))
	DB = db
}

// ConnectDbMySQL 初始化Mysql db
func ConnectDbMySQL(host, port, database, user, pass, charset string) gorm.Dialector {
	dns := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		user,
		pass,
		host,
		port,
		database,
		charset,
	)

	return mysql.Open(dns)
}

// ConnectDbPostgreSQL 连接dbpostgresql数据库
func ConnectDbPostgreSQL(host, port, database, user, pass string) gorm.Dialector {
	dns := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s",
		host,
		port,
		user,
		database,
		pass,
	)

	return postgres.Open(dns)
}

// ConnectDbSqlite3 连接sqlite3数据库
func ConnectDbSqlite3(host string) gorm.Dialector {
	dns := fmt.Sprintf(
		"host=%s",
		host,
	)
	return sqlite.Open(dns)
}

// ConnectDbSqlServer 连接 sql server数据库
func ConnectDbSqlServer(host, port, database, user, pass string) gorm.Dialector {
	dns := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		user,
		pass,
		host,
		port,
		database,
	)
	return sqlserver.Open(dns)
}
