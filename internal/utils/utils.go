package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/julianstephens/go-utils/cliutil"
)

const (
	DEFAULT_TIMEOUT         = 30 * time.Second
	DEFAULT_KEY_SIZE        = 32 // 256 bits
	DEFAULT_KEY_FILE        = "key.bin"
	DEFAULT_AUTH_CACHE_FILE = "auth.json"
)

const (
	lower     = "abcdefghijklmnopqrstuvwxyz"
	upper     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits    = "0123456789"
	allChars  = lower + upper + digits
	minLength = 8
)

func GenerateID() string {
	return uuid.New().String()
}

func PrintUser(id, email, name, createdAt, updatedAt string, roles []string) {
	cliutil.PrintInfo("User Details:")
	cliutil.PrintInfo("-------------")
	if id != "" && id != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  ID:        %s", id))
	}
	if email != "" && email != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  Email:     %s", email))
	}
	if name != "" && name != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  Name:      %s", name))
	}
	if createdAt != "" && createdAt != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  Created At: %s", createdAt))
	}
	if updatedAt != "" && updatedAt != "N/A" {
		cliutil.PrintInfo(fmt.Sprintf("  Updated At: %s", updatedAt))
	}
	if len(roles) > 0 {
		cliutil.PrintInfo(fmt.Sprintf("  Roles:     %s", strings.Join(roles, ", ")))
	} else {
		cliutil.PrintInfo("  Roles:     None")
	}
}

// GenerateTempPassword generates a temporary password that is at least 8 characters long
// and contains at least one uppercase letter, one lowercase letter, and one number.
func GenerateTempPassword() (string, error) {
	if minLength < 3 {
		return "", fmt.Errorf("password length too short")
	}

	password := make([]byte, minLength)

	// Ensure at least one character from each required set
	sets := []string{lower, upper, digits}
	for i, set := range sets {
		char, err := randomCharFrom(set)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Fill the rest with random characters from all sets
	for i := 3; i < minLength; i++ {
		char, err := randomCharFrom(allChars)
		if err != nil {
			return "", err
		}
		password[i] = char
	}

	// Shuffle password for randomness
	shuffle(password)

	return string(password), nil
}

func randomCharFrom(set string) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
	if err != nil {
		return 0, err
	}
	return set[n.Int64()], nil
}

func shuffle(data []byte) {
	for i := range data {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(len(data))))
		data[i], data[j.Int64()] = data[j.Int64()], data[i]
	}
}
