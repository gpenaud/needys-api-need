package internal

import (
  fmt     "fmt"
  http    "net/http"
  json    "encoding/json"
  log     "github.com/sirupsen/logrus"
  mux     "github.com/gorilla/mux"
  need    "github.com/gpenaud/needys-api-need/internal/need"
  runtime "runtime"
  strconv "strconv"
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

func (a *Application) InitializeDB(w http.ResponseWriter, _ *http.Request) {
  initialized, err := a.InitializeDatabase()

  if (initialized) {
    payload := map[string]bool{"initialized": initialized}
    respondWithJSON(w, http.StatusOK, payload)
  } else {
    handlerLog.Info(err)
    respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Database is not initializable - Error: %s", err.Error()))
  }
}

// -------------------------------------------------------------------------- //

func (a *Application) createNeed(w http.ResponseWriter, r *http.Request) {
  need := need.Need{}

  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&need)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "The payload is invalid")
    return
  }
  defer r.Body.Close()

  err = need.CreateNeed(a.DB)

  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  } else {
    respondWithJSON(w, http.StatusOK, need)
  }
}

func (a *Application) getNeed(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  id, err := strconv.Atoi(vars["id"])
  if err != nil {
    respondWithError(w, http.StatusBadRequest, fmt.Sprintf("The need with Id %d is invalid", id))
    return
  }

  need := need.Need{Id: id}
  err = need.GetNeed(a.DB)

  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  } else {
    respondWithJSON(w, http.StatusOK, need)
  }
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

func (a *Application) updateNeed(w http.ResponseWriter, r *http.Request) {
  need := need.Need{}

  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&need)
  if err != nil {
    respondWithError(w, http.StatusBadRequest, "The payload is invalid")
    return
  }
  defer r.Body.Close()

  err = need.UpdateNeed(a.DB)

  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  } else {
    respondWithJSON(w, http.StatusOK, need)
  }
}

func (a *Application) deleteNeed(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)

  need := need.Need{Name: vars["name"]}
  err  := need.DeleteNeed(a.DB)

  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
  } else {
    respondWithJSON(w, http.StatusOK, need)
  }
}
