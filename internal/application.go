package internal

import (
  amqp    "github.com/streadway/amqp"
  context "context"
  errors  "errors"
  fmt     "fmt"
  http    "net/http"
  log     "github.com/sirupsen/logrus"
  _       "github.com/lib/pq"
  mux     "github.com/gorilla/mux"
  sql     "database/sql"
  time    "time"
)

// -------------------------------------------------------------------------- //
// log

var applicationLog *log.Entry

var LogLevels = map[string]log.Level{
  "fatal":   log.FatalLevel,
  "error":   log.ErrorLevel,
  "warning": log.WarnLevel,
  "info":    log.InfoLevel,
  "debug":   log.DebugLevel,
}

var LogFormatters = map[string]log.Formatter{
  "development": &log.TextFormatter{},
  "integration": &log.JSONFormatter{},
  "production":  &log.JSONFormatter{},
  "text":        &log.TextFormatter{},
  "json":        &log.JSONFormatter{},
}

func (a *Application) initializeLogger() {
  // configure log verbosity
  log.SetLevel(LogLevels[a.Config.Verbosity])

  if (a.Config.Verbosity == "debug") {
    log.SetReportCaller(false)
  }

  // if log format is specified, configure it, else we base our choice on the environment
  if a.Config.LogFormat != "" {
    log.SetFormatter(LogFormatters[a.Config.LogFormat])
  } else {
    log.SetFormatter(LogFormatters[a.Config.Environment])
  }
}

func init() {
  applicationLog = log.WithFields(log.Fields{
    "_file": "internal/application.go",
    "_type": "system",
  })
}

// -------------------------------------------------------------------------- //
// configuration

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

// -------------------------------------------------------------------------- //
// backends

func (a* Application) initializeMessagingConnection() {
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

func (a *Application) IsMessagingReachable() (reachable bool, err error) {
  if (a.AMQP == nil) {
    return false, errors.New("Messaging connection is not set-up")
  }

  channel, err := a.AMQP.Channel()
  if (err != nil) {
    return false, errors.New(fmt.Sprintf("Messaging is not available - Error: %s", err.Error()))
  }

  defer channel.Close()

  return true, nil
}

func (a* Application) initializeDatabaseConnection() {
  dbDriver   := "mysql"
  dbCharset  := "charset=utf8mb4&collation=utf8mb4_unicode_ci"

  var err error

  a.DB, err = sql.Open(dbDriver, a.Config.Database.Username + ":" + a.Config.Database.Password + "@tcp(" + a.Config.Database.Host + ":" + a.Config.Database.Port + ")/" + a.Config.Database.Name + "?" + dbCharset)
  if err != nil {
    applicationLog.Error("Database server is not available")
  }

  a.DB.SetMaxIdleConns(3)
  a.DB.SetMaxOpenConns(10)
  a.DB.SetConnMaxLifetime(3600 * time.Second)
}

func (a *Application) IsDatabaseReachable() (reachable bool, err error) {
  if (a.DB == nil) {
    return false, errors.New("Database connection is not set-up")
  }

  _, err = a.DB.Query("SELECT null")
  if (err != nil) {
    return false, errors.New(fmt.Sprintf("Database is not available - Error: %s", err.Error()))
  }

  return true, nil
}

// -------------------------------------------------------------------------- //
// router

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

// -------------------------------------------------------------------------- //
// application

func (a *Application) Initialize() {
  applicationLog.WithFields(log.Fields{
    "database_host": a.Config.Database.Host,
    "database_port": a.Config.Database.Port,
    "database_username": a.Config.Database.Username,
    "database_name":   a.Config.Database.Name,
  }).Info("trying to connect to database")

  a.Router = mux.NewRouter()

  a.initializeDatabaseConnection()
  a.initializeMessagingConnection()
  a.initializeLogger()
  a.initializeRoutes()

  applicationLog.Info("application is initialized")
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
