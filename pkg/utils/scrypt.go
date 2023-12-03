package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

type passwordParams struct {
	R  int
	N  int
	s  int
	sp int
}

var nilString = ""

var saltLength = 23
var DefaultPasswordParams = passwordParams{R: 8, N: 14, s: 43}

var MisMatchedError = fmt.Errorf("pass: provided password does not match the actual password")

func CreatePassword(password string) (string, error) {
	salt, err := saltBytes(saltLength)
	if err != nil {
		return nilString, err
	}

	dk, err := scrypt.Key([]byte(password), salt, 1<<uint(DefaultPasswordParams.N), DefaultPasswordParams.R, 1, 16)
	if err != nil {
		return nilString, err
	}

	hash := encodeHash(salt) + encodeHash(dk)
	hashToReturn := fmt.Sprintf("$disto$%d$%d$%s", DefaultPasswordParams.N, DefaultPasswordParams.R, hash)
	return hashToReturn, nil
}

func ComparePassword(password string, hash string) error {
	var values []string = strings.Split(hash, "$")
	N, memErr := strconv.Atoi(values[2])
	if memErr != nil {
		return memErr
	}
	R, RErr := strconv.Atoi(values[3])
	if RErr != nil {
		return RErr
	}

	var params passwordParams = passwordParams{N: N, R: R, s: DefaultPasswordParams.s}
	salt, saltErr := decodeHash(values[4][:params.s])
	if saltErr != nil {
		return saltErr
	}

	dk, dkErr := scrypt.Key([]byte(password), salt, 1<<uint(params.N), params.R, 1, 16)
	if dkErr != nil {
		return dkErr
	}

	actualDk, err := decodeHash(values[4][params.s:])
	if err != nil {
		return err
	}

	if subtle.ConstantTimeCompare(dk, actualDk) == 1 {
		return nil
	}
	return MisMatchedError
}

func saltBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func encodeHash(hash []byte) string {
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(hash)
}

func decodeHash(hash string) ([]byte, error) {
	return base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(hash)
}
