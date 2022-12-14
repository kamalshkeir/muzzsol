package tests

import (
	"testing"

	"github.com/kamalshkeir/muzzsol/services"
)

// Test Encryption Service
func EncryptDecryptTest(t *testing.T) {
	expected := "secret"
	// encrypt
	enc,err := services.Encrypt(expected)
	if err != nil {
		t.Fatal(err)
	}
	// decrypt
	dec,err := services.Decrypt(enc)
	if err != nil {
		t.Fatal(err)
	}
	if dec != expected {
		t.Errorf("expected decrypted string to be %s but got %s",expected,dec)
	}
}

// Test Hashing Service
func HashingUnhashingTest(t *testing.T) {
	expected := "secret"
	// hash
	hash,err := services.GenerateHash(expected)
	if err != nil {
		t.Fatal(err)
	}
	// unhash
	match,err := services.ComparePasswordToHash(expected,hash)
	if err != nil || !match {
		t.Error("hash doesn't match")
	}
}

// ...