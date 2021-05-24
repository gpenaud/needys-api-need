package log

import (
  "io"
  "log"
  "os"
)

var (
  InfoLogger    *log.Logger
  WarningLogger *log.Logger
  ErrorLogger   *log.Logger
)

func init() {
  file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
  if err != nil {
    log.Fatal(err)
  }

  InfoLogger = log.New(io.MultiWriter(file ,os.Stdout), "INFO: ", log.Ldate|log.Ltime)
  WarningLogger = log.New(io.MultiWriter(file ,os.Stdout), "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
  ErrorLogger = log.New(io.MultiWriter(file ,os.Stdout), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
