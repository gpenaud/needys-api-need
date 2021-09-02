package internal

// -----------------------------------------------------------------------------
// 1. Database initialization (if specified in configuration)
// -----------------------------------------------------------------------------

const dbInitQuery = `
  CREATE TABLE IF NOT EXISTS need (
    id INTEGER PRIMARY KEY NOT NULL AUTO_INCREMENT,
    name VARCHAR(100),
    priority VARCHAR(100)
  ) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
  `

func (a *Application) InitializeDatabase() (bool, error) {
  initialized := false
  var err error

  if _, err := a.DB.Exec(dbInitQuery); err == nil {
    initialized = true
  }

  return initialized, err
}
