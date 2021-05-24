package mysql

import (
  "fmt"
  _ "github.com/go-sql-driver/mysql"
  // local imports
  // "github.com/gpenaud/needys-api-need/internal/models"
  "github.com/gpenaud/needys-api-need/internal/objects"
  utils_mysql "github.com/gpenaud/needys-api-need/pkg/mysql"
  "github.com/gpenaud/needys-api-need/pkg/log"
)

func DeleteNeedsByName(name *string) {
  query := fmt.Sprintf("DELETE FROM need WHERE name = '%s'", *name)
  log.InfoLogger.Println(query)

  db := utils_mysql.DbConn()
  _, err := db.Query(query)

  if err != nil {
    panic(err.Error())
  }

  defer db.Close()
}

func InsertNeed(name *string, priority *string) {
  query := fmt.Sprintf("INSERT INTO need (name, priority) VALUES ('%s', '%s')", *name, *priority)
  log.InfoLogger.Println(query)

  db := utils_mysql.DbConn()
  _, err := db.Query(query)

  if err != nil {
    panic(err.Error())
  }

  defer db.Close()
}

func GetNeeds() *[]objects.Need {
  db := utils_mysql.DbConn()
  selDB, err := db.Query("SELECT * FROM need ORDER BY id ASC")

  if err != nil {
    log.ErrorLogger.Fatalln(err.Error())
  }

  need := objects.Need{}
  needs := []objects.Need{}

  for selDB.Next() {
    var id int
    var name, priority string

    err = selDB.Scan(&id, &name, &priority)
    if err != nil {
      log.ErrorLogger.Fatalln(err.Error())
    }

    need.Id = id
    need.Name = name
    need.Priority = priority

    needs = append(needs, need)
  }

  defer db.Close()

  return &needs
}
