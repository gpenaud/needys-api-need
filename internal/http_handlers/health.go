package http_handlers

import (
  "encoding/json"
  "net/http"
)

func HealthyHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "Method is not supported.", http.StatusNotFound)
    return
  }

  healthy := map[string]bool{"is_healthy": true}

  json_healthy, err := json.Marshal(healthy)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(json_healthy)
}

// -------------------------------------------------------------------------- //

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Error(w, "Method is not supported.", http.StatusNotFound)
    return
  }

  ready := map[string]bool{"is_ready": true}

  json_ready, err := json.Marshal(ready)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(json_ready)
}
