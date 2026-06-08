package storage

import (
	"bytes"
	"errors"
	"path/filepath"
	"testing"

	"password_vault/internal/vault"
)

// newTestStore spins up a SQLiteStore backed by a throwaway temp-file DB.
// t.TempDir() is auto-removed when the test finishes, and a real file (not
// :memory:) avoids database/sql connection-pool surprises.
func newTestStore(t *testing.T) *SQLiteStore {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test.db")
	store, err := NewSQLiteStore(dbPath)
	if err != nil {
		t.Fatalf("NewSQLiteStore: %v", err)
	}
	return store
}

func TestGetSaltOnFreshDBReturnsErrNoSalt(t *testing.T) {
	store := newTestStore(t)

	// A brand-new vault has no salt row yet; the store must translate
	// sql.ErrNoRows into the vault.ErrNoSalt sentinel that NewService expects.
	_, err := store.GetSalt()
	if !errors.Is(err, vault.ErrNoSalt) {
		t.Fatalf("expected vault.ErrNoSalt on fresh DB, got %v", err)
	}
}

func TestSaveAndGetSalt(t *testing.T) {
	store := newTestStore(t)
	salt := []byte("0123456789abcdef")

	if err := store.SaveSalt(salt); err != nil {
		t.Fatalf("SaveSalt: %v", err)
	}
	got, err := store.GetSalt()
	if err != nil {
		t.Fatalf("GetSalt: %v", err)
	}
	if !bytes.Equal(got, salt) {
		t.Fatalf("salt round-trip mismatch: got %x want %x", got, salt)
	}
}

func TestAddAndGetEntry(t *testing.T) {
	store := newTestStore(t)

	// The store deals only in bytes, so fake "ciphertext" is fine here.
	in := vault.EncryptedEntry{
		Username: []byte("enc-user"),
		Password: []byte("enc-pass"),
		URL:      []byte("enc-url"),
		Notes:    []byte("enc-notes"),
	}
	if err := store.AddEntry(in); err != nil {
		t.Fatalf("AddEntry: %v", err)
	}

	// First inserted row gets id 1 from AUTOINCREMENT.
	got, err := store.GetEntry(1)
	if err != nil {
		t.Fatalf("GetEntry: %v", err)
	}
	if !bytes.Equal(got.Username, in.Username) ||
		!bytes.Equal(got.Password, in.Password) ||
		!bytes.Equal(got.URL, in.URL) ||
		!bytes.Equal(got.Notes, in.Notes) {
		t.Fatalf("entry round-trip mismatch: got %+v want %+v", got, in)
	}
}

func TestListEntries(t *testing.T) {
	store := newTestStore(t)

	for i := 0; i < 3; i++ {
		entry := vault.EncryptedEntry{
			Username: []byte("u"),
			Password: []byte("p"),
			URL:      []byte("url"),
			Notes:    []byte("n"),
		}
		if err := store.AddEntry(entry); err != nil {
			t.Fatalf("AddEntry %d: %v", i, err)
		}
	}

	entries, err := store.ListEntries()
	if err != nil {
		t.Fatalf("ListEntries: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
}

func TestDeleteEntry(t *testing.T) {
	store := newTestStore(t)

	if err := store.AddEntry(vault.EncryptedEntry{
		Username: []byte("u"), Password: []byte("p"), URL: []byte("url"), Notes: []byte("n"),
	}); err != nil {
		t.Fatalf("AddEntry: %v", err)
	}
	if err := store.DeleteEntry(1); err != nil {
		t.Fatalf("DeleteEntry: %v", err)
	}

	entries, err := store.ListEntries()
	if err != nil {
		t.Fatalf("ListEntries: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries after delete, got %d", len(entries))
	}
}
