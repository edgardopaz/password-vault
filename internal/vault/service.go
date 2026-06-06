package vault

import "errors"

type Store interface {
	// TODO: Implememnt the store interface using EncryptedEntry + salt
	AddEntry(entry EncryptedEntry) error
	GetEntry(id int) (EncryptedEntry, error)
	ListEntries() ([]EncryptedEntry, error)
	UpdateEntry(entry EncryptedEntry) error
	DeleteEntry(id int) error
	GetSalt() ([]byte, error)
	SaveSalt(salt []byte) error
}

type Service struct {
	store Store
	key []byte
}

var ErrNoSalt = errors.New("no salt found")

func NewService(store Store, masterPassword string) (*Service, error) {
	salt, err := store.GetSalt()

	switch {
	case errors.Is(err,ErrNoSalt):

		salt, err = generateSalt()
		if err != nil {
			return nil, err
		}

		if err := store.SaveSalt(salt); err != nil {
			return nil, err
		}
	case err != nil:
		return nil, err
	}

	key := deriveKey(masterPassword, salt)
	return &Service{store: store, key: key}, nil
}

func (s *Service) AddEntry(entry Entry) error {
	// encrypt the entry using the key
	encryptedUsername, err := encrypt([]byte(entry.Username), s.key)
	if err != nil {
		return err
	}
	encryptedURL, err := encrypt([]byte(entry.URL), s.key)
	if err != nil {
		return err
	}
	encryptedNotes, err := encrypt([]byte(entry.Notes), s.key)
	if err != nil {
		return err
	}
	encryptedPassword, err := encrypt([]byte(entry.Password), s.key)
	if err != nil {
		return err
	}
	encryptedEntry := EncryptedEntry{
		ID: entry.ID,
		Username: encryptedUsername,
		Password: encryptedPassword,
		URL: encryptedURL,
		Notes: encryptedNotes,
	}
	return s.store.AddEntry(encryptedEntry)
}
func (s *Service) GetEntry(id int) (Entry, error) {
	return Entry{}, nil	
}
func (s *Service) ListEntries() ([]Entry, error) {
	return []Entry{}, nil
}
func (s *Service) UpdateEntry(entry Entry) error {
	return nil
}
func (s *Service) DeleteEntry(id int) error {
	return nil
}