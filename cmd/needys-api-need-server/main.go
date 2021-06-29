package main

import (
  cmdline  "github.com/galdor/go-cmdline"
  context  "context"
  internal "github.com/gpenaud/needys-api-need/internal"
  log      "github.com/sirupsen/logrus"
  os       "os"
  signal   "os/signal"
  syscall  "syscall"
)

var mainLog *log.Entry
var a        internal.Application

func init() {
  mainLog = log.WithFields(log.Fields{
    "_file": "cmd/needys-api-need-server/main.go",
    "_type": "system",
  })

  registerCliConfiguration(&a)
  registerVersion(&a)

  a.Initialize()
}

var PossibleOptionValues = map[string][]string{
  "environment": {"development", "integration", "production"},
  "verbosity": {"error", "warning", "info", "debug"},
  "log-format": {"text", "json"},
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str { return true }
	}

	return false
}

func registerCliConfiguration(a *internal.Application) {
  cmdline := cmdline.New()

  a.Config = &internal.Configuration{}

  // application configuration flags
  cmdline.AddOption("e", "environment", "ENVIRONMENT", "the current environment (development, integration, production)")
  cmdline.SetOptionDefault("environment", "production")

  cmdline.AddOption("v", "verbosity", "LEVEL", "verbosity for log-level (error, warning, info, debug)")
  cmdline.SetOptionDefault("verbosity", "info")

  cmdline.AddOption("l", "log-format", "FORMAT", "log format (text, json)")
  cmdline.SetOptionDefault("log-format", "unset")

  cmdline.AddFlag("", "log-healthcheck", "log healthcheck queries")

  // application server configuration flags
  cmdline.AddOption("", "server.host", "HOST", "host of application")
  cmdline.SetOptionDefault("server.host", "localhost")

  cmdline.AddOption("", "server.port", "PORT", "port of application")
  cmdline.SetOptionDefault("server.port", "8010")

  // db configuration flags
  cmdline.AddOption("", "database.host", "HOST", "host of the MySQL server")
  cmdline.SetOptionDefault("database.host", "127.0.0.1")

  cmdline.AddOption("", "database.port", "PORT", "port of the MySQL server")
  cmdline.SetOptionDefault("database.port", "3306")

  cmdline.AddOption("", "database.username", "USERNAME", "username of MySQL server")
  cmdline.SetOptionDefault("database.username", "needys")

  cmdline.AddOption("", "database.password", "PASSWORD", "password of the MySQL user")
  cmdline.SetOptionDefault("database.password", "needys")

  cmdline.AddOption("", "database.name", "DB_NAME", "the database name")
  cmdline.SetOptionDefault("database.name", "needys")

  // rabbitmq configuration flags
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

  // application general configuration
  if (! contains(PossibleOptionValues["environment"], cmdline.OptionValue("environment"))) {
    mainLog.WithFields(log.Fields{
      "environment": a.Config.Environment,
    }).Fatal("Wrong value for option environment")
  } else {
    a.Config.Environment = cmdline.OptionValue("environment")
  }

  if (! contains(PossibleOptionValues["verbosity"], cmdline.OptionValue("verbosity"))) {
    mainLog.WithFields(log.Fields{
      "verbosity": a.Config.Verbosity,
    }).Fatal("Wrong value for option verbosity")
  } else {
    a.Config.Verbosity = cmdline.OptionValue("verbosity")
  }

  if (cmdline.OptionValue("log-format") != "unset") {
    if (! contains(PossibleOptionValues["log-format"], cmdline.OptionValue("log-format"))) {
      mainLog.WithFields(log.Fields{
        "log-format": a.Config.LogFormat,
      }).Fatal("Wrong value for option log-format")
    } else {
      a.Config.LogFormat = cmdline.OptionValue("log-format")
    }
  }

  a.Config.LogHealthcheck = cmdline.IsOptionSet("log-healthcheck")

  // a server configuration values
  a.Config.Server.Host = cmdline.OptionValue("server.host")
  a.Config.Server.Port = cmdline.OptionValue("server.port")

  // database configuration value
  a.Config.Database.Host     = cmdline.OptionValue("database.host")
  a.Config.Database.Port     = cmdline.OptionValue("database.port")
  a.Config.Database.Name     = cmdline.OptionValue("database.name")
  a.Config.Database.Username = cmdline.OptionValue("database.username")
  a.Config.Database.Password = cmdline.OptionValue("database.password")

  // rabitmq configuration value
  a.Config.Rabbitmq.Host     = cmdline.OptionValue("rabbitmq.host")
  a.Config.Rabbitmq.Port     = cmdline.OptionValue("rabbitmq.port")
  a.Config.Rabbitmq.Username = cmdline.OptionValue("rabbitmq.username")
  a.Config.Rabbitmq.Password = cmdline.OptionValue("rabbitmq.password")
}

var BuildTime = "unset"
var Commit 		= "unset"
var Release 	= "unset"

func registerVersion(a *internal.Application) {
  a.Version = &internal.Version{BuildTime, Commit, Release}
}

func main() {
  c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())

  go func() {
		oscall := <-c

    mainLog.WithFields(log.Fields{
      "signal": oscall,
    }).Warn("received a system call")

    a.DB.Close()
    a.AMQP.Close()

		cancel()
	}()

  a.Run(ctx)
}
