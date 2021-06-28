package internal

import (
  amqp    "github.com/streadway/amqp"
  context "context"
  fmt     "fmt"
  http    "net/http"
  log     "github.com/sirupsen/logrus"
  _       "github.com/lib/pq"
  mux     "github.com/gorilla/mux"
  sql     "database/sql"
  time    "time"
)

var applicationLog *log.Entry

func init() {
  applicationLog = log.WithFields(log.Fields{
    "_file": "internal/application.go",
    "_type": "system",
  })
}

type Configuration struct {
  Environment    string
  Verbosity      string
  LogFormat      string
  LogHealthcheck bool
  Healthcheck struct {
    Timeout  int
  }
  Server struct {
    Port string
    Host string
  }
  Database struct {
    Port string
    Host string
    Username string
    Password string
    Name string
  }
  Rabbitmq struct {
    Port string
    Host string
    Username string
    Password string
  }
}

type Version struct {
  BuildTime string
  Commit    string
  Release   string
}

type Application struct {
  AMQP    *amqp.Connection
  DB      *sql.DB
  Config  *Configuration
  Router  *mux.Router
  Version *Version
}

func (a *Application) IsDatabaseReachable() (reachable bool, err error) {
  _, err = a.DB.Query("SELECT null")
  return (err == nil), err
}

func (a *Application) IsMessagingReachable() (reachable bool, err error) {
  channel, err := a.AMQP.Channel()
  defer channel.Close()
  return (err == nil), err
}

func (a *Application) Initialize() {
  applicationLog.WithFields(log.Fields{
    "database_host": a.Config.Database.Host,
    "database_port": a.Config.Database.Port,
    "database_username": a.Config.Database.Username,
    "database_name":   a.Config.Database.Name,
  }).Info("trying to connect to database")

  a.Router = mux.NewRouter()

  a.initializeDatabaseConnection()
  a.initializeAMQPConnection()
  a.initializeLogger()
  a.initializeRoutes()

  applicationLog.Info("application is initialized")
}

func (a* Application) initializeAMQPConnection() {
  amqp_connection_parameters := fmt.Sprintf(
    "amqp://%s:%s@%s:%s/",
    a.Config.Rabbitmq.Username,
    a.Config.Rabbitmq.Password,
    a.Config.Rabbitmq.Host,
    a.Config.Rabbitmq.Port,
  )

  var err error

  a.AMQP, err = amqp.Dial(amqp_connection_parameters)
  if err != nil {
    applicationLog.Error("Failed to open RabbitMQ connection")
  }
}

func (a* Application) initializeDatabaseConnection() {
  dbDriver   := "mysql"

  dbHost     := a.Config.Database.Host
  dbPort     := a.Config.Database.Port
  dbUser     := a.Config.Database.Username
  dbPassword := a.Config.Database.Password
  dbName     := a.Config.Database.Name
  dbCharset  := "charset=utf8mb4&collation=utf8mb4_unicode_ci"

  var err error

  a.DB, err = sql.Open(dbDriver, dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?" + dbCharset)
  if err != nil {
    applicationLog.Error("Database server is not available")
  }

  a.DB.SetMaxIdleConns(3)
  a.DB.SetMaxOpenConns(10)
  a.DB.SetConnMaxLifetime(3600 * time.Second)
}

func (a *Application) initializeLogger() {
  switch a.Config.Verbosity {
  case "fatal":
    log.SetLevel(log.FatalLevel)
  case "error":
    log.SetLevel(log.ErrorLevel)
  case "warning":
    log.SetLevel(log.WarnLevel)
  case "info":
    log.SetLevel(log.InfoLevel)
  case "debug":
    log.SetLevel(log.DebugLevel)
    log.SetReportCaller(false)
  default:
    log.WithFields(
      log.Fields{"verbosity": a.Config.Verbosity},
    ).Fatal("Unkown verbosity level")
  }

  switch a.Config.Environment {
  case "development":
    log.SetFormatter(&log.TextFormatter{})
  case "integration":
    log.SetFormatter(&log.JSONFormatter{})
  case "production":
    log.SetFormatter(&log.JSONFormatter{})
  default:
    log.WithFields(
      log.Fields{"environment": a.Config.Environment},
    ).Fatal("Unkown environment type")
  }

  if a.Config.LogFormat != "unset" {
    switch a.Config.LogFormat {
    case "text":
      log.SetFormatter(&log.TextFormatter{})
    case "json":
      log.SetFormatter(&log.JSONFormatter{})
    default:
      log.WithFields(
        log.Fields{"log_format": a.Config.LogFormat},
      ).Fatal("Unkown log format")
    }
  }
}

func (a *Application) initializeRoutes() {
  // application need-related routes
  a.Router.HandleFunc("/needs", a.getNeeds).Methods("GET")
  a.Router.HandleFunc("/need", a.createNeed).Methods("POST")
  // a.Router.HandleFunc("/need/{id:[0-9]+}", a.getNeed).Methods("GET")
  a.Router.HandleFunc("/need/{id:[0-9]+}", a.updateNeed).Methods("PUT")
  a.Router.HandleFunc("/need/{id:[0-9]+}", a.deleteNeed).Methods("DELETE")
  // application probes routes
  a.Router.HandleFunc("/health", a.isHealthy).Methods("GET")
  a.Router.HandleFunc("/ready", a.isReady).Methods("GET")
  // application maintenance routes
  a.Router.HandleFunc("/initialize_db", a.InitializeDB).Methods("GET")
}

func (a *Application) Run(ctx context.Context) {
  server_address :=
    fmt.Sprintf("%s:%s", a.Config.Server.Host, a.Config.Server.Port)

  server_message :=
    fmt.Sprintf(
  `

START INFOS
-----------
Listening needys-api-need on %s:%s...

BUILD INFOS
-----------
time: %s
release: %s
commit: %s

`,
      a.Config.Server.Host,
      a.Config.Server.Port,
      a.Version.BuildTime,
      a.Version.Release,
      a.Version.Commit,
    )

  httpServer := &http.Server{
		Addr:    server_address,
		Handler: a.Router,
	}

  go func() {
    // we keep this log on standard format
    log.Info(server_message)
    applicationLog.Fatal(httpServer.ListenAndServe())
  }()

  <-ctx.Done()
  applicationLog.Info("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

  var err error

	if err = httpServer.Shutdown(ctxShutDown); err != nil {
    applicationLog.WithFields(log.Fields{
      "error": err,
    }).Fatal("server shutdown failed")
	}

  applicationLog.Info("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
