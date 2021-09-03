package internal

import (
  amqp          "github.com/streadway/amqp"
  context       "context"
  fmt           "fmt"
  healthcheck   "github.com/hellofresh/health-go/v4"
  http          "net/http"
  log           "github.com/sirupsen/logrus"
  _             "github.com/lib/pq"
  mux           "github.com/gorilla/mux"
  databasecheck "github.com/hellofresh/health-go/v4/checks/mysql"
  httpcheck     "github.com/hellofresh/health-go/v4/checks/http"
  sql           "database/sql"
  time          "time"
)

// -------------------------------------------------------------------------- //
// 1. Application Logging
// -------------------------------------------------------------------------- //

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
  if a.Config.LogFormat != "unset" {
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
// 2. Application Declarative Configuration
// -------------------------------------------------------------------------- //

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
    Initialize bool
  }
  Messaging struct {
    Port string
    Host string
    Username string
    Password string
    Vhost string
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
// 3. Backends Initialization (MariaDB / RabbitMQ)
// -------------------------------------------------------------------------- //

func (a* Application) initializeMessagingConnection() {
  amqp_connection_parameters := fmt.Sprintf(
    "amqp://%s:%s@%s:%s/%s",
    a.Config.Messaging.Username,
    a.Config.Messaging.Password,
    a.Config.Messaging.Host,
    a.Config.Messaging.Port,
    a.Config.Messaging.Vhost,
  )

  var err error

  a.AMQP, err = amqp.Dial(amqp_connection_parameters)
  if err != nil {
    applicationLog.Fatal(fmt.Sprintf("Failed to open RabbitMQ connection: %s", amqp_connection_parameters))
  }
}

func (a* Application) initializeDatabaseConnection() {
  dbDriver          := "mysql"
  dbCharset         := "charset=utf8mb4&collation=utf8mb4_unicode_ci"
  dbConnectionQuery := a.Config.Database.Username + ":" + a.Config.Database.Password + "@tcp(" + a.Config.Database.Host + ":" + a.Config.Database.Port + ")/" + a.Config.Database.Name + "?" + dbCharset

  var err error

  a.DB, err = sql.Open(dbDriver, dbConnectionQuery)
  if err != nil {
    applicationLog.Fatal(fmt.Sprintf("Database server is not available: %s", dbConnectionQuery))
  }

  a.DB.SetMaxIdleConns(3)
  a.DB.SetMaxOpenConns(10)
  a.DB.SetConnMaxLifetime(3600 * time.Second)
}

// -------------------------------------------------------------------------- //
// 4. Router setup
// -------------------------------------------------------------------------- //

func (a *Application) initializeRoutes() {
  // application need-related routes
  a.Router.HandleFunc("/needs", a.getNeeds).Methods("GET")
  a.Router.HandleFunc("/need/{id:[0-9]+}", a.getNeed).Methods("GET")
  a.Router.HandleFunc("/need", a.updateNeed).Methods("POST")
  a.Router.HandleFunc("/need/{id:[0-9]+}", a.updateNeed).Methods("PUT")
  a.Router.HandleFunc("/need/{name:[a-zA-Z]+}", a.deleteNeed).Methods("DELETE")
  // application maintenance routes
  a.Router.HandleFunc("/initialize_db", a.InitializeDB).Methods("GET")
}

// -------------------------------------------------------------------------- //
// 5. Application Setup
// -------------------------------------------------------------------------- //

func (a *Application) Initialize() {
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

  // ---------------------------------------------------------------------------
  // 5.1. manage healthchecks and healthchecks server
  // ---------------------------------------------------------------------------

  // health checks
  health, _ := healthcheck.New()

  health.Register(healthcheck.Config{
		Name:  "health-check",
		Check: func(context.Context) error { return nil },
	})

  // live checks
  live, _ := healthcheck.New()

  live.Register(healthcheck.Config{
		Name:      "mysql-check",
		Timeout:   time.Second * 5,
		SkipOnErr: false,
		Check: databasecheck.New(databasecheck.Config{
			DSN: fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", a.Config.Database.Username, a.Config.Database.Password, a.Config.Database.Host, a.Config.Database.Port, a.Config.Database.Name),
		}),
	})

  live.Register(healthcheck.Config{
		Name:      "rabbit-aliveness-check",
		Timeout:   time.Second * 5,
		SkipOnErr: false,
		Check: httpcheck.New(httpcheck.Config{
			URL: fmt.Sprintf(
        "http://%s:%s@%s:%s/api/aliveness-test/%s", a.Config.Messaging.Username, a.Config.Messaging.Password, a.Config.Messaging.Host, "15672", a.Config.Messaging.Vhost),
		}),
	})

  http.Handle("/health", health.Handler())
  http.Handle("/live", live.Handler())

  // expose healthcheck endpoints
  healthcheckServer := &http.Server{
    Addr:    "0.0.0.0:8090",
    Handler: nil,
  }

  go func() {
    healthcheckServer.ListenAndServe()
  }()

  // ---------------------------------------------------------------------------
  // 5.2. manage healthchecks and healthchecks server
  // ---------------------------------------------------------------------------

  httpServer := &http.Server{
		Addr:    server_address,
		Handler: a.Router,
	}

  go func() {
    // we keep this log on standard format
    log.Info(server_message)
    httpServer.ListenAndServe()
  }()

  // -------------------------------------------------
  // 5.3. initialize database if specified in configuration
  // -------------------------------------------------

  if (a.Config.Database.Initialize) {
    initialized, err := a.InitializeDatabase()

    if (initialized) {
      applicationLog.Info("database initialisation succeeded")
    } else {
      applicationLog.WithFields(log.Fields{
        "error": err,
      }).Fatal("database initialisation failed")
    }
  }

  // ---------------------------------------------------------------------------
  // 5.4. manage server shutdown
  // ---------------------------------------------------------------------------

  <-ctx.Done()
  applicationLog.Info("application server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

  var err error

	if err = httpServer.Shutdown(ctxShutdown); err != nil {
    applicationLog.WithFields(log.Fields{
      "error": err,
    }).Fatal("application server shutdown failed")
	}

  applicationLog.Info("application server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
