package auth

import (
	"fmt"
	"io"
)

type AuthRepo interface {
	GenerateAuthURL() (string, error)
	SaveAuthToken(string) error
}

type Service struct {
	rw       io.ReadWriter
	authRepo AuthRepo
}

func NewService(rw io.ReadWriter, authRepo AuthRepo) *Service {
	return &Service{
		rw:       rw,
		authRepo: authRepo,
	}
}

func (s *Service) Handle() error {
	url, err := s.authRepo.GenerateAuthURL()
	if err != nil {
		return err
	}

	if _, err := s.rw.Write([]byte(url)); err != nil {
		return err
	}

	var authCode []byte
	if _, err := s.rw.Read(authCode); err != nil {
		return fmt.Errorf("unable to read authorization code %v", err)
	}

	return s.authRepo.SaveAuthToken(string(authCode))
}
