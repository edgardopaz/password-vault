package vault

type Store interface {
	// TODO: Implememnt the store interface to store the encrypted entries
	NewService(store Store, masterPassword string) *Service
	AddEntry(entry Entry) error
	GetEntry(id int) (Entry, error)
	ListEntries() ([]Entry, error)
	UpdateEntry(entry Entry) error
	DeleteEntry(id int) error
}

type Service struct {
	store Store
	key []byte
}

NewService(store Store, masterPassword string) *Service {
	salt, err := generateSalt()
	if err != nil {
		return nil, err
	}
	key := deriveKey(masterPassword, salt)
	return &Service{store: store, key: key}, nil
}