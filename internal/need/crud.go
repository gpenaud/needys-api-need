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

func (n *Need) CreateNeed(db *sql.DB) (err error) {
  needLog.WithFields(log.Fields{
    "type": "database query",
    "parameter_name": n.Name,
    "parameter_priority": n.Priority,
  }).Debug("INSERT INTO need (name, priority) VALUES ({name}, {priority})")

  r, err := db.Query("INSERT INTO need (name, priority) VALUES (?, ?)", n.Name, n.Priority)
  r.Close()

  return err
}

func (n *Need) GetNeed(db *sql.DB) (error) {
  needLog.WithFields(log.Fields{
    "type": "database query",
    "parameter_id": n.Id,
  }).Debug("SELECT name, priority FROM need WHERE id={id}")

  needLog.Debug(n.Id)

  err := db.QueryRow("SELECT name, priority FROM need WHERE id=?", n.Id).Scan(&n.Name, &n.Priority)
  return err
}

func (n *Need) UpdateNeed(db *sql.DB) (err error) {
  needLog.WithFields(log.Fields{
    "type": "database query",
    "parameter_priority": n.Priority,
    "parameter_name": n.Name,
  }).Debug("UPDATE need SET priority = '{priority}' WHERE name = '{name}'")

  r, err := db.Query("UPDATE need SET priority = ? WHERE name = ?", n.Priority, n.Name)
  r.Close()

  return err
}

func (n *Need) DeleteNeed(db *sql.DB) error {
  needLog.WithFields(log.Fields{
    "type": "database query",
    "parameter_name": n.Name,
  }).Debug("DELETE FROM need WHERE name={name}")

  r, err := db.Query("DELETE FROM need WHERE name=?", n.Name)
  r.Close()

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
