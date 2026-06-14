package storage

import "database/sql"
// the CREATE TABLE statements for the databse and a function to run them on startup

const createEntriesTable = `
CREATE TABLE IF NOT EXISTS entries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username BLOB NOT NULL,
    password BLOB NOT NULL,
    url BLOB NOT NULL,
    notes BLOB NOT NULL
)
`

const createMetadataTable = `CREATE TABLE IF NOT EXISTS metadata (
    id INTEGER PRIMARY KEY CHECK (id = 1),
    salt BLOB NOT NULL,
    verifier BLOB NOT NULL
)
`

func migrate(db *sql.DB) error {
    _, err := db.Exec(createEntriesTable)
    if err != nil {
        return err
    }
    _, err = db.Exec(createMetadataTable)
    if err != nil {
        return err
    }
    return nil
}