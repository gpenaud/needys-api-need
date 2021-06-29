package internal

import (
  fmt     "fmt"
  http    "net/http"
  json    "encoding/json"
  log     "github.com/sirupsen/logrus"
  need    "github.com/gpenaud/needys-api-need/internal/need"
  runtime "runtime"
  strings "strings"
)

var handlerLog *log.Entry

func init() {
  log.SetReportCaller(true)
  if pc, file, line, ok := runtime.Caller(1); ok {
    file = file[strings.LastIndex(file, "/")+1:]
		funcName := runtime.FuncForPC(pc).Name()
    handlerLog = log.WithFields(log.Fields{
      "_src": fmt.Sprintf("%s:%s:%d", file, funcName, line),
      "_type": "router",
    })
	}
}

// -------------------------------------------------------------------------- //
// Common functions for handlers

func respondHTTPCodeOnly(w http.ResponseWriter, code int) {
  w.WriteHeader(code)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
  handlerLog.Error(message)
  respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  response, _ := json.Marshal(payload)
  handlerLog.Debug(fmt.Sprintf("JSON response: %s", response))

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}

// -------------------------------------------------------------------------- //
// Maintenance handlers

const dbInitQuery = `
  CREATE TABLE IF NOT EXISTS need (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    priority VARCHAR(100)
  ) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
  `

func (a *Application) InitializeDB(w http.ResponseWriter, _ *http.Request) {
  if is_reachable, err := a.IsDatabaseReachable(); !is_reachable {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  } else {
    if _, err = a.DB.Exec(dbInitQuery); err == nil {
      payload := map[string]bool{
        "initialized": true,
      }
      respondWithJSON(w, http.StatusOK, payload)
    } else {
      handlerLog.Info(err)
      respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database is not initializable - Error: %s", err.Error()))
    }
  }
}

// -------------------------------------------------------------------------- //

func (a *Application) isHealthy(w http.ResponseWriter, _ *http.Request) {
  w.WriteHeader(http.StatusOK)
}

func (a *Application) isReady(w http.ResponseWriter, r *http.Request) {
  http_status := http.StatusOK

  if is_reachable, err := a.IsDatabaseReachable(); !is_reachable {
    http_status = http.StatusInternalServerError
    log.Debug(err.Error())

    // if (err.Error() == "Database connection is not set-up") {
    //   handlerLog.Info(fmt.Sprintf("Self-healing: trying to re-etablish Database server connection after error: %s"), err.Error())
    //   a.initializeDatabaseConnection()
    // } else {
    //   handlerLog.Error(err.Error())
    // }
  }

  if is_reachable, err := a.IsMessagingReachable(); !is_reachable {
    http_status = http.StatusInternalServerError
    log.Debug(err.Error())

    if (err.Error() == "Messaging connection is not set-up") {
      handlerLog.Info("Self-healing: trying to re-etablish Messaging server connection")
      a.initializeMessagingConnection()
    } else {
      handlerLog.Error(err.Error())
    }
  }

  w.WriteHeader(http_status)
  return
}

func (a *Application) getNeeds(w http.ResponseWriter, r *http.Request) {
  need       := need.Need{}
  needs, err := need.GetNeeds(a.DB)

  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  } else {
    respondWithJSON(w, http.StatusOK, needs)
  }
}

func (a *Application) createNeed(w http.ResponseWriter, r *http.Request) {
  r.ParseForm()

  name     := r.FormValue("name")
  priority := r.FormValue("priority")

  need := need.Need{Name: name, Priority: priority}
  err := need.InsertNeed(a.DB)

  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  } else {
    respondWithJSON(w, http.StatusOK, need)
  }
}

func (a *Application) updateNeed(w http.ResponseWriter, r *http.Request) {
  fmt.Println(w, r)
}

func (a *Application) deleteNeed(w http.ResponseWriter, r *http.Request) {
  name := r.URL.Query().Get("name")
  need := need.Need{Name: name}

  err := need.DeleteNeed(a.DB)

  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  } else {
    respondWithJSON(w, http.StatusOK, need)
  }
}
