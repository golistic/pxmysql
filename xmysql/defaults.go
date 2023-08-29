// Copyright (c) 2023, Geert JM Vanderkelen

package xmysql

const authChallengeLen = 20

const (
	AuthMethodPlain        AuthMethodType = "PLAIN"
	AuthMethodAuto         AuthMethodType = "AUTO"
	AuthMethodSHA256Memory AuthMethodType = "SHA256_MEMORY"
	AuthMethodMySQL41      AuthMethodType = "MYSQL41"
)

const DefaultPort = "33060"
const DefaultHost = "127.0.0.1"

type AuthMethodType string

type AuthMethodTypes []AuthMethodType

func (a AuthMethodTypes) Has(m AuthMethodType) bool {
	for _, v := range a {
		if v == m {
			return true
		}
	}
	return false
}

var defaultAuthMethods = []AuthMethodType{AuthMethodMySQL41, AuthMethodSHA256Memory}

var supportedAuthMethods = AuthMethodTypes{AuthMethodSHA256Memory, AuthMethodMySQL41, AuthMethodPlain, AuthMethodAuto}

func DefaultAuthMethods() []AuthMethodType {
	return defaultAuthMethods
}

func SupportedAuthMethods() AuthMethodTypes {
	return supportedAuthMethods
}
