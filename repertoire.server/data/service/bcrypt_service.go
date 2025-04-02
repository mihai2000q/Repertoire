package service

import "golang.org/x/crypto/bcrypt"

type BCryptService interface {
	Hash(str string) (string, error)
}

type bCryptService struct {
}

func NewBCryptService() BCryptService {
	return new(bCryptService)
}

func (b bCryptService) Hash(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
	return string(hash), err
}
