package password

import (
	"crypto/subtle"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Hash returns a bcrypt hash of the plaintext password.
func Hash(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// isBcryptHash reports whether stored looks like a bcrypt hash.
func isBcryptHash(stored string) bool {
	return strings.HasPrefix(stored, "$2a$") ||
		strings.HasPrefix(stored, "$2b$") ||
		strings.HasPrefix(stored, "$2y$")
}

// Check verifies plain against a stored bcrypt hash or legacy plaintext value.
// needsRehash is true when the stored value was plaintext and matched (caller should upgrade).
func Check(stored, plain string) (ok bool, needsRehash bool) {
	if isBcryptHash(stored) {
		err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(plain))
		return err == nil, false
	}

	if subtle.ConstantTimeCompare([]byte(stored), []byte(plain)) == 1 {
		return true, true
	}
	return false, false
}
