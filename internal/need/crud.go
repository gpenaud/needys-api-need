package need

import (
  _   "github.com/go-sql-driver/mysql"
  log "github.com/sirupsen/logrus"
  sql "database/sql"
)

var needLog *log.Entry

func init() {
  needLog = log.WithFields(log.Fields{
    "_file": "internal/need/crud.go",
    "_type": "data",
  })
}

type Need struct {
  Id       int
  Name     string
  Priority string
}

func (n *Need) DeleteNeed(db *sql.DB) error {
  needLog.WithFields(log.Fields{
    "type": "database query",
    "parameter_name": n.Name,
  }).Debug("DELETE FROM need WHERE name={name}")

  _, err := db.Exec("DELETE FROM need WHERE name=$1", n.Name)

  return err
}

func (n *Need) InsertNeed(db *sql.DB) (err error) {
  needLog.WithFields(log.Fields{
    "type": "database query",
    "parameter_name": n.Name,
    "parameter_priority": n.Priority,
  }).Debug("INSERT INTO need (name, priority) VALUES ({name}, {priority})")

  _, err = db.Exec("INSERT INTO need (name, priority) VALUES ('%s', '%s')", n.Name, n.Priority)

  return err
}

func (n *Need) GetNeeds(db *sql.DB) ([]Need, error) {
  need  := Need{}
  needs := []Need{}

  selDB, err := db.Query("SELECT * FROM need ORDER BY id ASC")

  if err != nil {
    return needs, err
  }

  for selDB.Next() {
    var id int
    var name, priority string

    err = selDB.Scan(&id, &name, &priority)
    if err != nil {
      return needs, err
    }

    need.Id = id
    need.Name = name
    need.Priority = priority

    needs = append(needs, need)
  }

  return needs, err
}
