package storage

// sqlite store struct that opens the DB, and the store method implementations

import (
	"database/sql"
	_ "modernc.org/sqlite"
	"password_vault/internal/vault"
	"errors"
)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	
	if err != nil {
		return nil, err
	}
	// ping the database to ensure it is connected
	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := migrate(db); err != nil {
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

var _ vault.Store = (*SQLiteStore)(nil)

func (s *SQLiteStore) AddEntry(entry vault.EncryptedEntry) error {
	_, err := s.db.Exec("INSERT INTO entries (username, password, url, notes) VALUES (?, ?, ?, ?)", entry.Username, entry.Password, entry.URL, entry.Notes)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) GetEntry(id int) (vault.EncryptedEntry, error) {
	var entry vault.EncryptedEntry
	err := s.db.QueryRow("SELECT id, username, password, url, notes FROM entries WHERE id = ?", id).Scan(&entry.ID, &entry.Username, &entry.Password, &entry.URL, &entry.Notes)
	if err != nil {
		return vault.EncryptedEntry{}, err
	}
	return entry, nil
}

func (s *SQLiteStore) ListEntries() ([]vault.EncryptedEntry, error) {
	rows, err := s.db.Query("SELECT id, username, password, url, notes FROM entries")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var entries []vault.EncryptedEntry
	for rows.Next() {
		var entry vault.EncryptedEntry
		err := rows.Scan(&entry.ID, &entry.Username, &entry.Password, &entry.URL, &entry.Notes)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (s *SQLiteStore) UpdateEntry(entry vault.EncryptedEntry) error {
	_, err := s.db.Exec("UPDATE entries SET username = ?, password = ?, url = ?, notes = ? WHERE id = ?", entry.Username, entry.Password, entry.URL, entry.Notes, entry.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) DeleteEntry(id int) error {
	_, err := s.db.Exec("DELETE FROM entries WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (s *SQLiteStore) GetSalt() ([]byte, error) {
	var salt []byte
	err := s.db.QueryRow("SELECT salt FROM metadata WHERE id = 1").Scan(&salt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, vault.ErrNoSalt
	}
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func (s *SQLiteStore) SaveSalt(salt []byte) error {
	_, err := s.db.Exec("INSERT OR REPLACE INTO metadata (id, salt) VALUES (1, ?)", salt)
	if err != nil {
		return err
	}
	return nil
}