package main

import (
  "fmt"
  "github.com/galdor/go-cmdline"
  "log"
  "net/http"
  "os"
  "time"
  // local imports
  "github.com/gpenaud/needys-api-need/internal/handler"
  "github.com/gpenaud/needys-api-need/internal/models"
)

func defaultHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case http.MethodGet:
    handler.ListHandler(w, r)
  case http.MethodPost:
    handler.InsertHandler(w, r)
  case http.MethodPut:
      // Update an existing record.
  case http.MethodDelete:
    handler.DeleteHandler(w, r)
  default:
    http.Error(w, "404 not found.", http.StatusNotFound)
  }

  return
}

func main() {
  fmt.Printf("starting needys-api-need at port 8010\n")

  cmdline := cmdline.New()

  cmdline.AddOption("", "server.host", "HOST", "host of application")
  cmdline.SetOptionDefault("server.host", "localhost")
  cmdline.AddOption("", "server.port", "PORT", "port of application")
  cmdline.SetOptionDefault("server.port", "8010")

  cmdline.AddOption("", "mysql.host", "HOST", "host of the MySQL server")
  cmdline.SetOptionDefault("mysql.host", "127.0.0.1")
  cmdline.AddOption("", "mysql.port", "PORT", "port of the MySQL server")
  cmdline.SetOptionDefault("mysql.port", "3306")
  cmdline.AddOption("", "mysql.username", "USERNAME", "username of MySQL server")
  cmdline.SetOptionDefault("mysql.username", "needys")
  cmdline.AddOption("", "mysql.password", "PASSWORD", "password of the MySQL user")
  cmdline.SetOptionDefault("mysql.password", "needys")
  cmdline.AddOption("", "mysql.dbname", "DB_NAME", "the database name")
  cmdline.SetOptionDefault("mysql.dbname", "needys")

  cmdline.AddOption("", "rabbitmq.host", "HOST", "host of the rabbitMQ server")
  cmdline.SetOptionDefault("rabbitmq.host", "127.0.0.1")
  cmdline.AddOption("", "rabbitmq.port", "PORT", "port of the rabbitMQ server")
  cmdline.SetOptionDefault("rabbitmq.port", "5672")
  cmdline.AddOption("", "rabbitmq.username", "USERNAME", "username of rabbitMQ server")
  cmdline.SetOptionDefault("rabbitmq.username", "guest")
  cmdline.AddOption("", "rabbitmq.password", "PASSWORD", "password of the rabbitMQ user")
  cmdline.SetOptionDefault("rabbitmq.password", "guest")

  cmdline.AddFlag("v", "verbose", "log more information")
  cmdline.Parse(os.Args)

  models.Cfg.Server.Host = cmdline.OptionValue("server.host")
  models.Cfg.Server.Port = cmdline.OptionValue("server.port")

  models.Cfg.Mysql.Host = cmdline.OptionValue("mysql.host")
  models.Cfg.Mysql.Port = cmdline.OptionValue("mysql.port")
  models.Cfg.Mysql.Username = cmdline.OptionValue("mysql.username")
  models.Cfg.Mysql.Password = cmdline.OptionValue("mysql.password")
  models.Cfg.Mysql.Dbname = cmdline.OptionValue("mysql.dbname")

  models.Cfg.Rabbitmq.Host = cmdline.OptionValue("rabbitmq.host")
  models.Cfg.Rabbitmq.Port = cmdline.OptionValue("rabbitmq.port")
  models.Cfg.Rabbitmq.Username = cmdline.OptionValue("rabbitmq.username")
  models.Cfg.Rabbitmq.Password = cmdline.OptionValue("rabbitmq.password")

  http.HandleFunc("/", defaultHandler)

  server_address := fmt.Sprintf("%s:%s", models.Cfg.Server.Host, models.Cfg.Server.Port)
  server := &http.Server{
    Addr:           server_address,
    ReadTimeout:    10 * time.Second,
    WriteTimeout:   10 * time.Second,
    MaxHeaderBytes: 1 << 20,
  }

  log.Fatal(server.ListenAndServe())
}
