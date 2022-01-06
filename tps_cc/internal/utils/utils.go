package utils

import (
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func GenerateID() string {
	uuidStr := uuid.NewString()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	return uuidStr
}

func GenerateRamdomPath(base string) string {
	path := base
	exists := true
	for exists {
		path = base + "/" + RandSeq(8)
		exists = false
		if _, err := os.Stat(path); err == nil {
			exists = true
		}
	}

	os.MkdirAll(path, os.ModePerm)
	return path
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
