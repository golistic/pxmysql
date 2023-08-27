// Copyright (c) 2022, Geert JM Vanderkelen

package xmysql

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
)

type authn struct {
	username  string
	password  string
	schema    string
	challenge []byte
}

// authSHA256Data prepares authentication data to be sent with the AuthenticateContinue
// message using SHA256. Username and scrambled password are returned as hex.
// See: https://dev.mysql.com/doc/internals/en/x-protocol-authentication-authentication.html.
func authSHA256Data(an authn) ([]byte, error) {
	if len(an.challenge) != authChallengeLen {
		return nil, fmt.Errorf("authentication challenge must be 20 bytes (was %d)", len(an.challenge))
	}

	var scramble string
	if an.password != "" {
		// hex(sha256(password) XOR sha256(challenge + sha256(sha256(password))))
		h1 := sha256.Sum256([]byte(an.password))
		hh1 := sha256.Sum256(h1[:])

		hr := sha256.New()
		hr.Write(hh1[:])
		hr.Write(an.challenge)
		h2 := hr.Sum(nil)

		for i := range h2 {
			h1[i] ^= h2[i]
		}
		scramble = fmt.Sprintf("%x", h1)
	}

	return []byte(fmt.Sprintf("%s\x00%s\x00%s", an.schema, an.username, scramble)), nil
}

// authMYSQL41Data prepares authentication data to be sent with the AuthenticateContinue
// message using SHA1 (also known as mysql_native_password). Username and scrambled password
// are returned as hex.
// See: https://dev.mysql.com/doc/internals/en/x-protocol-authentication-authentication.html.
func authMySQL41Data(an authn) ([]byte, error) {
	if len(an.challenge) != authChallengeLen {
		return nil, fmt.Errorf("authentication challenge must be 20 bytes (was %d)", len(an.challenge))
	}

	var scramble string
	if an.password != "" {
		// hex(sha1(password) XOR sha1(challenge + sha1(sha1(password))))
		h1 := sha1.Sum([]byte(an.password))
		hh1 := sha1.Sum(h1[:])

		hr := sha1.New()
		hr.Write(an.challenge)
		hr.Write(hh1[:])
		h2 := hr.Sum(nil)

		for i := range h1 {
			h1[i] ^= h2[i]
		}

		scramble = fmt.Sprintf("*%x", h1)
	}

	return []byte(fmt.Sprintf("%s\x00%s\x00%s", an.schema, an.username, scramble)), nil
}

// authMySQLPlain prepares authentication data to be sent in plain text. This is only
// supported when connection is encrypted (TLS)
// See: https://dev.mysql.com/doc/internals/en/x-protocol-authentication-authentication.html.
func authMySQLPlain(an authn) ([]byte, error) {
	return []byte(fmt.Sprintf("%s\x00%s\x00%s", an.schema, an.username, an.password)), nil
}
