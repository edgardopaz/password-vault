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
	encryptedEntry, err := s.encryptEntry(entry)
	if err != nil {
		return err
	}
	return s.store.AddEntry(encryptedEntry)
}

func (s *Service) GetEntry(id int) (Entry, error) {
	// decrypt the entry using the key
	encryptedEntry, err := s.store.GetEntry(id)
	if err != nil {
		return Entry{}, err
	}
	
	decryptedUsername, err := decrypt(encryptedEntry.Username, s.key)
	if err != nil {
		return Entry{}, err
	}
	decryptedPassword, err := decrypt(encryptedEntry.Password, s.key)
	if err != nil {
		return Entry{}, err
	}
	decryptedURL, err := decrypt(encryptedEntry.URL, s.key)
	if err != nil {
		return Entry{}, err
	}
	decryptedNotes, err := decrypt(encryptedEntry.Notes, s.key)
	if err != nil {
		return Entry{}, err
	}
	return Entry{ID: id, Username: string(decryptedUsername), Password: string(decryptedPassword), URL: string(decryptedURL), Notes: string(decryptedNotes)}, nil
}

func (s *Service) ListEntries() ([]EntryMetadata, error) {
	encryptedEntries, err := s.store.ListEntries()
	if err != nil {
		return []EntryMetadata{}, err
	}
	entries := make([]EntryMetadata, len(encryptedEntries))
	for i, encryptedEntry := range encryptedEntries {	
		decryptedUsername, err := decrypt(encryptedEntry.Username, s.key)
		if err != nil {
			return []EntryMetadata{}, err
		}
		decryptedURL, err := decrypt(encryptedEntry.URL, s.key)
		if err != nil {
			return []EntryMetadata{}, err
		}
		entries[i] = EntryMetadata{ID: encryptedEntry.ID, Username: string(decryptedUsername), URL: string(decryptedURL)}
	}
	return entries, nil
}
func (s *Service) UpdateEntry(entry Entry) error {
	encryptedEntry, err := s.encryptEntry(entry)
	if err != nil {
		return err
	}
	return s.store.UpdateEntry(encryptedEntry)
}

func (s* Service) DeleteEntry(id int) error {
	return s.store.DeleteEntry(id)
}

func (s* Service) encryptEntry(entry Entry) (EncryptedEntry, error) {

	encryptedUsername, err := encrypt([]byte(entry.Username), s.key)
	if err != nil {
		return EncryptedEntry{}, err
	}
	encryptedURL, err := encrypt([]byte(entry.URL), s.key)
	if err != nil {
		return EncryptedEntry{}, err
	}
	encryptedNotes, err := encrypt([]byte(entry.Notes), s.key)
	if err != nil {
		return EncryptedEntry{}, err
	}
	encryptedPassword, err := encrypt([]byte(entry.Password), s.key)
	if err != nil {
		return EncryptedEntry{}, err
	}
	return EncryptedEntry{ID: entry.ID, Username: encryptedUsername, Password: encryptedPassword, URL: encryptedURL, Notes: encryptedNotes}, nil
}