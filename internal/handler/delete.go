package handler

import (
  "fmt"
  "net/http"
  "encoding/json"
  // "database/sql"
  _ "github.com/go-sql-driver/mysql"
  // local imports
  "github.com/gpenaud/needys-api-need/internal/utils"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")

  query := fmt.Sprintf("DELETE FROM need WHERE name = '%s'", name)
  fmt.Println(query)

  db := utils.DbConn()
  _, err := db.Query(query)

  if err != nil {
    panic(err.Error())
  }

  json_healthy, err := json.MarshalIndent(*utils.ShowAll(), "", "  ")
  if err != nil {
    panic(err.Error())
  }

  defer db.Close()

  w.Header().Set("Content-Type", "application/json")
  w.Write(json_healthy)
}
