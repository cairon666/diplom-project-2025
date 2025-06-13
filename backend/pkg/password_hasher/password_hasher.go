package password_hasher

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher struct {
}

func NewPasswordHasher() *PasswordHasher {
	return &PasswordHasher{}
}

func (h *PasswordHasher) Hash(password string) (string, string, error) {
	salt := uuid.New().String()
	bytes, err := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)

	return string(bytes), salt, err
}
func (h *PasswordHasher) Compare(password, salt, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password+salt))

	return err == nil
}
