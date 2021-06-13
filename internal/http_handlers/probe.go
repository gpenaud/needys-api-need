package http_handlers

import (
  "fmt"
  "net/http"
  // local imports
  "github.com/streadway/amqp"
  "github.com/gpenaud/needys-api-need/internal/config"
  "github.com/gpenaud/needys-api-need/pkg/log"
  "github.com/gpenaud/needys-api-need/pkg/mysql"
)

func HealthHandler(w http.ResponseWriter, _ *http.Request) {
  log.InfoLogger.Println("health")
  w.WriteHeader(http.StatusOK)
}

// -------------------------------------------------------------------------- //

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
  query := "SELECT null"
  http_status := http.StatusOK

  db := mysql.DbConn()
  _, err := db.Query(query)

  defer db.Close()

  if err != nil {
    http_status = http.StatusInternalServerError
    log.WarningLogger.Println("Mysql server is not available")
  }

  rabbitmq_connection_parameters := fmt.Sprintf(
    "amqp://%s:%s@%s:%s/",
    config.Cfg.Rabbitmq.Username,
    config.Cfg.Rabbitmq.Password,
    config.Cfg.Rabbitmq.Host,
    config.Cfg.Rabbitmq.Port,
  )

  conn, err := amqp.Dial(rabbitmq_connection_parameters)
  
  if err != nil {
    http_status = http.StatusInternalServerError
    log.WarningLogger.Println("Rabbitmq server is not available")
  } else {
    conn.Close()
  }

  w.WriteHeader(http_status)
  return
}