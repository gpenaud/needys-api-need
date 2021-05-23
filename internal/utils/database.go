package utils

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  // local imports
  "github.com/gpenaud/needys-api-need/internal/models"
)

func DbConn() (db *sql.DB) {
  dbDriver   := "mysql"

  dbHost     := models.Cfg.Mysql.Host
  dbPort     := models.Cfg.Mysql.Port
  dbUser     := models.Cfg.Mysql.Username
  dbPassword := models.Cfg.Mysql.Password
  dbName     := models.Cfg.Mysql.Dbname
  dbCharset  := "charset=utf8mb4&collation=utf8mb4_unicode_ci"

  db, err := sql.Open(dbDriver, dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?" + dbCharset)

  if err != nil {
    panic(err.Error())
  }

  return db
}

func ShowAll() *[]models.Need {
  db := DbConn()
  selDB, err := db.Query("SELECT * FROM need ORDER BY id ASC")

  if err != nil {
    panic(err.Error())
  }

  need := models.Need{}
  needs := []models.Need{}

  for selDB.Next() {
    var id int
    var name, priority string

    err = selDB.Scan(&id, &name, &priority)
    if err != nil {
        panic(err.Error())
    }

    need.Id = id
    need.Name = name
    need.Priority = priority

    needs = append(needs, need)
  }

  defer db.Close()

  return &needs
}
