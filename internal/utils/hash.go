package utils

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func Hash(s string) (string, error) {
	hash, err := argon2id.CreateHash(s, argon2id.DefaultParams)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return hash, nil
}

func Verify(s, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(s, hash)

	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return match, nil
}
