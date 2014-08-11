package main

import "log"
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
type PasswordHandler func(*Password, chan *Password)

// Hash the given Password using bcrypt.
func Hash(password *Password, results chan *Password) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password.Plain), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("ERROR Hash: %s\n", err.Error())
		results <- nil
	}
	results <- &Password{Hash: string(hash)}
}

// Compare the given Password's plaintext representation to its hash.
func Compare(password *Password, results chan *Password) {
	match := (bcrypt.CompareHashAndPassword(
		[]byte(password.Hash), []byte(password.Plain)) == nil)
	results <- &Password{Match: match}
}
