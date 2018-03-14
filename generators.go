package main

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateMail get an email
func GenerateMail(username string) string {
	return username + "@spdrun.com"
}

// GenerateRandomName creates a random name
func GenerateRandomName() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return "spdr" + strconv.Itoa(r.Intn(999999))
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

// GenerateRandomPassword Creates a random 12 char password
func GenerateRandomPassword() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, 12)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}

// GenerateID generates a random id
func GenerateID() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(999999)
}
