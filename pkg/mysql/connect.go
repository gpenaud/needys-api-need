package mysql

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  // local imports
  "github.com/gpenaud/needys-api-need/internal/config"
)

func DbConn() (db *sql.DB) {
  dbDriver   := "mysql"

  dbHost     := config.Cfg.Mysql.Host
  dbPort     := config.Cfg.Mysql.Port
  dbUser     := config.Cfg.Mysql.Username
  dbPassword := config.Cfg.Mysql.Password
  dbName     := config.Cfg.Mysql.Dbname
  dbCharset  := "charset=utf8mb4&collation=utf8mb4_unicode_ci"

  db, err := sql.Open(dbDriver, dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?" + dbCharset)

  if err != nil {
    panic(err.Error())
  }

  return db
}
