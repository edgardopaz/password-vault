package vault

type Entry struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	URL      string `json:"url" db:"url"`
	Notes    string `json:"notes" db:"notes"`
}

type EncryptedEntry struct {
	ID       int    `json:"id" db:"id"`
	Username []byte `json:"username" db:"username"`
	Password []byte `json:"password" db:"password"`
	URL      []byte `json:"url" db:"url"`
	Notes    []byte `json:"notes" db:"notes"`
}

type EncryptionMetadata struct {
	ID   int    `json:"id" db:"id"`
	Salt []byte `json:"salt" db:"salt"`
}