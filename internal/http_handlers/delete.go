package http_handlers

import (
  "net/http"
  "encoding/json"
  // "database/sql"
  _ "github.com/go-sql-driver/mysql"
  // local imports
  "github.com/gpenaud/needys-api-need/internal/mysql"
)

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")

  mysql.DeleteNeedsByName(&name)

  json_healthy, err := json.MarshalIndent(*mysql.GetNeeds(), "", "  ")
  if err != nil {
    panic(err.Error())
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(json_healthy)
}
