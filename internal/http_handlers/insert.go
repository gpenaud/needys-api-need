package http_handlers

import (
  "net/http"
  "encoding/json"
  _ "github.com/go-sql-driver/mysql"
  // local imports
  "github.com/gpenaud/needys-api-need/internal/mysql"
  "github.com/gpenaud/needys-api-need/internal/rabbitmq_sender"
)

func InsertHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Method POST is required", http.StatusNotFound)
    return
  }

  r.ParseForm()

  name     := r.FormValue("name")
  priority := r.FormValue("priority")

  mysql.InsertNeed(&name, &priority)

  json_healthy, err := json.MarshalIndent(*mysql.GetNeeds(), "", "  ")
  if err != nil {
    panic(err.Error())
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(json_healthy)

  rabbitmq_sender.SendAmqpMessages(json_healthy)
}
