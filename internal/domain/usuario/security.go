package usuario

import (
	"golang.org/x/crypto/bcrypt"
)

func HashSenha(senha string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	return string(bytes), err
}

func VerificarSenha(senha, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha))
}
