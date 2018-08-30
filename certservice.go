package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

// adjust this to change expiration date
var expiration = time.Duration(24 * time.Hour)

// simulated sleep duration
var sleepDuration = time.Duration(10 * time.Second)

// RequestCert simulates requesting a cert and returning it
func RequestCert(domain string) (string, time.Time, error) {

	time.Sleep(sleepDuration)

	h := sha256.New()
	h.Write([]byte(domain))

	var cert = fmt.Sprintf("%x", h.Sum(nil))

	var expires = time.Now().Add(expiration)

	return cert, expires, nil
}
