package handler

import (
  "net/http"
  "encoding/json"
  // local imports
  "github.com/gpenaud/needys-api-need/internal/utils"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Method GET is required", http.StatusNotFound)
    return
  }

  json_healthy, err := json.MarshalIndent(*utils.ShowAll(), "", "  ")
  if err != nil {
    panic(err.Error())
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(json_healthy)
}
