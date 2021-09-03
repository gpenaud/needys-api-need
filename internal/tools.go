package internal

// -----------------------------------------------------------------------------
// 1. Database initialization (if specified in configuration)
// -----------------------------------------------------------------------------

const dbReset = `DROP TABLE IF EXISTS need`
const dbInit = `
  CREATE TABLE need (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    priority VARCHAR(100)
  ) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
  `
const dbInsertFirst  = `INSERT INTO need (name, priority) VALUES ('adventure', 'low');`
const dbInsertSecond = `INSERT INTO need (name, priority) VALUES ('rest', 'high');`

func (a *Application) InitializeDatabase() (bool, error) {
  var err error

  if _, err = a.DB.Exec(dbReset); err != nil {
    return false, err
  }

  if _, err = a.DB.Exec(dbInit); err != nil { 
    return false, err
  }

  if _, err = a.DB.Exec(dbInsertFirst); err != nil {
    return false, err
  }

  if _, err = a.DB.Exec(dbInsertSecond); err != nil {
    return false, err
  }

  return true, err
}
