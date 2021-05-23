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

func InsertHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Method POST is required", http.StatusNotFound)
    return
  }

  r.ParseForm()

  name     := r.FormValue("name")
  priority := r.FormValue("priority")

  query := fmt.Sprintf("INSERT INTO need (name, priority) VALUES ('%s', '%s')", name, priority)
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

  utils.SendAmqpMessages(json_healthy)
}
