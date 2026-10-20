package helper

import (
	"crypto/rand"
	"math/big"
)

const temporaryPasswordLength = 16

const (
	passwordUpper = "ABCDEFGHJKLMNPQRSTUVWXYZ"
	passwordLower = "abcdefghijkmnopqrstuvwxyz"
	passwordDigit = "23456789"
)

// GenerateTemporaryPassword returns a random initial password for a new employee.
// The value is not derived from the employee name and is not written to the process environment.
func GenerateTemporaryPassword() (string, error) {
	alphabet := passwordUpper + passwordLower + passwordDigit
	chars := make([]byte, temporaryPasswordLength)
	var err error
	chars[0], err = randomChar(passwordUpper)
	if err != nil {
		return "", err
	}
	chars[1], err = randomChar(passwordLower)
	if err != nil {
		return "", err
	}
	chars[2], err = randomChar(passwordDigit)
	if err != nil {
		return "", err
	}
	for i := 3; i < len(chars); i++ {
		chars[i], err = randomChar(alphabet)
		if err != nil {
			return "", err
		}
	}
	if err := shuffleBytes(chars); err != nil {
		return "", err
	}
	return string(chars), nil
}

func randomChar(alphabet string) (byte, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
	if err != nil {
		return 0, err
	}
	return alphabet[n.Int64()], nil
}

func shuffleBytes(chars []byte) error {
	for i := len(chars) - 1; i > 0; i-- {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		j := int(n.Int64())
		chars[i], chars[j] = chars[j], chars[i]
	}
	return nil
}
