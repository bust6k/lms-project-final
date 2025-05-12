package security

import "golang.org/x/crypto/bcrypt"

var DefaultCOst int = 16

func Hash(Password string) ([]byte, error) {
	bytePassword := []byte(Password)

	hashed, err := bcrypt.GenerateFromPassword(bytePassword, int(DefaultCOst))
	if err != nil {
		return nil, err
	}
	return hashed, nil
}

func Compare(password string, hashedPass []byte) error {
	return bcrypt.CompareHashAndPassword(hashedPass, []byte(password))
}
