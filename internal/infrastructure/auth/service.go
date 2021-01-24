package auth

import (
	"fmt"
	"io"
	"io/ioutil"
)

type authRepo interface {
	GenerateAuthURL() (string, error)
	SaveAuthToken(string) error
}

// TODO: Move this service to the Service package
type Service struct {
	rw       io.ReadWriter
	authRepo authRepo
}

func NewService(rw io.ReadWriter, authRepo authRepo) *Service {
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

	authCode, err := ioutil.ReadAll(s.rw)
	if err != nil {
		return fmt.Errorf("unable to read authorization code %v", err)
	}

	return s.authRepo.SaveAuthToken(string(authCode))
}
