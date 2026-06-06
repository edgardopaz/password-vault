package vault

import (
	"bytes"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	key := deriveKey("correct-horse", []byte("0123456789abcdef")) // 16-byte salt
	plaintext := []byte("hunter2")

	ciphertext, err := encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	got, err := decrypt(ciphertext, key)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if string(got) != string(plaintext) {
		t.Fatalf("round-trip mismatch: got %q want %q", got, plaintext)
	}
}

func TestDecryptWrongKeyFails(t *testing.T) {
	salt := []byte("0123456789abcdef")
	rightKey := deriveKey("correct-horse", salt)
	wrongKey := deriveKey("wrong-password", salt)
	plaintext := []byte("hunter2")

	ciphertext, err := encrypt(plaintext, rightKey)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	// AES-GCM authenticates the ciphertext, so decrypting with a key derived
	// from the wrong master password must fail rather than return garbage.
	if _, err := decrypt(ciphertext, wrongKey); err == nil {
		t.Fatal("expected decrypt with wrong key to fail, but it succeeded")
	}
}

func TestEncryptDoesNotLeakPlaintext(t *testing.T) {
	key := deriveKey("correct-horse", []byte("0123456789abcdef"))
	plaintext := []byte("hunter2")

	ciphertext, err := encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	if bytes.Contains(ciphertext, plaintext) {
		t.Fatalf("ciphertext contains plaintext %q: encryption did not happen", plaintext)
	}
}

func TestDeriveKeyDeterministic(t *testing.T) {
	salt := []byte("0123456789abcdef")
	otherSalt := []byte("fedcba9876543210")

	// Same password + same salt must always produce the same key,
	// otherwise data encrypted yesterday could never be decrypted today.
	k1 := deriveKey("correct-horse", salt)
	k2 := deriveKey("correct-horse", salt)
	if !bytes.Equal(k1, k2) {
		t.Fatal("deriveKey is not deterministic for the same password and salt")
	}

	// Same password + different salt must produce a different key,
	// which is the whole reason the salt exists.
	k3 := deriveKey("correct-horse", otherSalt)
	if bytes.Equal(k1, k3) {
		t.Fatal("deriveKey produced the same key for different salts")
	}
}

