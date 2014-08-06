package main

import "code.google.com/p/go.crypto/bcrypt"

func Hash(password *Password) (*Password, error) {
	hash, err := bcrypt.GenerateFromPassword(
		[]byte(password.Plain), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Password{Hash: string(hash)}, nil
}

func Compare(password *Password) (*Password, error) {
	match := (bcrypt.CompareHashAndPassword(
		[]byte(password.Hash), []byte(password.Plain)) == nil)
	return &Password{Match: match}, nil
}
