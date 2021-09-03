package main

import (
  cli      "github.com/urfave/cli/v2"
  context  "context"
  internal "github.com/gpenaud/needys-api-need/internal"
  log      "github.com/sirupsen/logrus"
  os       "os"
  signal   "os/signal"
  syscall  "syscall"
)

// -------------------------------------------------------------------------- //
// 1. Application Initialization
// -------------------------------------------------------------------------- //

var mainLog *log.Entry
var a        internal.Application

func init() {
  mainLog = log.WithFields(log.Fields{
    "_file": "cmd/needys-api-need-server/main.go",
    "_type": "system",
  })

  registerConfiguration(&a)
  registerVersion(&a)

  a.Initialize()
}

// -------------------------------------------------------------------------- //
// 2. Application Configuration
// -------------------------------------------------------------------------- //

var PossibleOptionValues = map[string][]string{
  "environment": {"development", "integration", "production"},
  "verbosity": {"error", "warning", "info", "debug"},
  "log-format": {"unset", "text", "json"},
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str { return true }
	}

	return false
}

func registerConfiguration(a *internal.Application) {
  a.Config = &internal.Configuration{}

  app := &cli.App{
    Action: func(c *cli.Context) error {
      return nil
    },
    Flags: []cli.Flag{
      &cli.StringFlag{Name: "environment", Aliases: []string{"e"}, Value: "development", Usage: "The current environment `ENV`", Destination: &a.Config.Environment, EnvVars: []string{"NEEDYS_API_NEED_ENVIRONMENT"}},
      &cli.StringFlag{Name: "verbosity", Aliases: []string{"v"}, Value: "info", Usage: "Verbosity `LEVEL` for log-level", Destination: &a.Config.Verbosity, EnvVars: []string{"NEEDYS_API_NEED_VERBOSITY"}},
      &cli.StringFlag{Name: "log-format", Aliases: []string{"l"}, Value: "unset", Usage: "Log formatter to use `FORMAT`", Destination: &a.Config.LogFormat, EnvVars: []string{"NEEDYS_API_NEED_LOG_FORMAT"}},
      &cli.BoolFlag  {Name: "log-healthcheck", Value: false, Usage: "Log healthcheck queries", Destination: &a.Config.LogHealthcheck, EnvVars: []string{"NEEDYS_API_NEED_LOG_HEALTHCHECK"}},
      &cli.StringFlag{Name: "server.host", Value: "127.0.0.1", Usage: "API server host `HOST`", Destination: &a.Config.Server.Host, EnvVars: []string{"NEEDYS_API_NEED_SERVER_HOST"}},
      &cli.StringFlag{Name: "server.port", Value: "8010", Usage: "API server port `PORT`", Destination: &a.Config.Server.Port, EnvVars: []string{"NEEDYS_API_NEED_SERVER_PORT"}},
      &cli.StringFlag{Name: "database.host", Value: "127.0.0.1", Usage: "Database host `HOST`", Destination: &a.Config.Database.Host, EnvVars: []string{"NEEDYS_API_NEED_DATABASE_HOST"}},
      &cli.StringFlag{Name: "database.port", Value: "3306", Usage: "Database port `PORT`", Destination: &a.Config.Database.Port, EnvVars: []string{"NEEDYS_API_NEED_DATABASE_PORT"}},
      &cli.StringFlag{Name: "database.username", Value: "needys", Usage: "Database user name `USERNAME`", Destination: &a.Config.Database.Username, EnvVars: []string{"NEEDYS_API_NEED_DATABASE_USERNAME"}},
      &cli.StringFlag{Name: "database.password", Value: "needys", Usage: "Database user password `PASSWORD`", Destination: &a.Config.Database.Password, EnvVars: []string{"NEEDYS_API_NEED_DATABASE_PASSWORD"}},
      &cli.StringFlag{Name: "database.name", Value: "needys", Usage: "Database name `NAME`", Destination: &a.Config.Database.Name, EnvVars: []string{"NEEDYS_API_NEED_DATABASE_NAME"}},
      &cli.BoolFlag  {Name: "database.initialize", Value: false, Usage: "Initialize database", Destination: &a.Config.Database.Initialize, EnvVars: []string{"NEEDYS_API_NEED_DATABASE_INITIALIZE"}},
      &cli.StringFlag{Name: "messaging.host", Value: "127.0.0.1", Usage: "Messaging host `HOST`", Destination: &a.Config.Messaging.Host, EnvVars: []string{"NEEDYS_API_NEED_MESSAGING_HOST"}},
      &cli.StringFlag{Name: "messaging.port", Value: "5672", Usage: "Messaging port `PORT`", Destination: &a.Config.Messaging.Port, EnvVars: []string{"NEEDYS_API_NEED_MESSAGING_PORT"}},
      &cli.StringFlag{Name: "messaging.username", Value: "guest", Usage: "Messaging Username `USERNAME`", Destination: &a.Config.Messaging.Username, EnvVars: []string{"NEEDYS_API_NEED_MESSAGING_USERNAME"}},
      &cli.StringFlag{Name: "messaging.password", Value: "guest", Usage: "Messaging password `PASSWORD`", Destination: &a.Config.Messaging.Password, EnvVars: []string{"NEEDYS_API_NEED_MESSAGING_PASSWORD"}},
      &cli.StringFlag{Name: "messaging.vhost", Value: "needys", Usage: "Messaging vhost `VHOST`", Destination: &a.Config.Messaging.Vhost, EnvVars: []string{"NEEDYS_API_NEED_MESSAGING_VHOST"}},
    },
  }

  err := app.Run(os.Args)
  if err != nil {
    log.Fatal(err)
  }

  // application general configuration
  if (! contains(PossibleOptionValues["environment"], a.Config.Environment)) {
    mainLog.WithFields(log.Fields{
      "environment": a.Config.Environment,
    }).Fatal("Wrong value for option environment (should be \"development\", \"integration\" or \"production\")")
  }

  if (! contains(PossibleOptionValues["verbosity"], a.Config.Verbosity)) {
    mainLog.WithFields(log.Fields{
      "verbosity": a.Config.Verbosity,
    }).Fatal("Wrong value for option verbosity (should be \"fatal\", \"error\", \"warning\", \"info\" or \"debug\")")
  }

  if (! contains(PossibleOptionValues["log-format"], a.Config.LogFormat)) {
    mainLog.WithFields(log.Fields{
      "log-format": a.Config.LogFormat,
    }).Fatal("Wrong value for option log-format (should be \"unset\", \"text\" or \"json\")")
  }
}

// -------------------------------------------------------------------------- //
// 3. Application Version
// -------------------------------------------------------------------------- //

var BuildTime = "unset"
var Commit 		= "unset"
var Release 	= "unset"

func registerVersion(a *internal.Application) {
  a.Version = &internal.Version{BuildTime, Commit, Release}
}

// -------------------------------------------------------------------------- //
// 4. Main function
// -------------------------------------------------------------------------- //

func main() {
  c := make(chan os.Signal, 1) // creation of a channel of type os.Signal
	signal.Notify(c, os.Interrupt, syscall.SIGKILL, syscall.SIGTERM) // add 2 signals to the channel
	ctx, cancel := context.WithCancel(context.Background())

  go func() {
		oscall := <-c

    mainLog.WithFields(log.Fields{
      "signal": oscall,
    }).Warn("received a system call")

		cancel() // linked with ctx, cancel
	}()

  a.Run(ctx)
}
