package main

import "code.google.com/p/go.crypto/bcrypt"

// Password represents a password sent and received by the client.
// When sent by client, can contain Hash and Plain attributes.
// When sent by the server to the client, can contain the Hash
// and Match attributes.
type Password struct {
	Hash  string `json:"hash,omitempty"`
	Plain string `json:"plain,omitempty"`
	Match bool   `json:"match,omitempty"`
}

// PasswordHandler is the signature for Hash and Compare functions.
type PasswordHandler func(*Password) (*Password, error)

// Hash the given Password using bcrypt.
func Hash(password *Password) (*Password, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password.Plain), bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return &Password{Hash: string(hash)}, nil
}

// Compare the given Password's plaintext representation to its hash.
func Compare(password *Password) (*Password, error) {
	match := (bcrypt.CompareHashAndPassword(
		[]byte(password.Hash), []byte(password.Plain)) == nil)
	return &Password{Match: match}, nil
}
