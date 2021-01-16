package auth

import (
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"io"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/api/drive/v3"
)

type ReadWriter struct {
}

func NewReadWriter() *ReadWriter {
	return &ReadWriter{}
}

func (i *ReadWriter) Read(p []byte) (n int, err error) {
	buf := make([]byte, 0)
	if n, err = fmt.Fscan(os.Stdin, &buf); err != nil {
		return n, fmt.Errorf("unable to read authorization code %v", err)
	}
	return copy(p, buf), io.EOF
}

func (i *ReadWriter) Write(url []byte) (n int, err error) {
	msg := fmt.Sprintf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", string(url))
	return fmt.Fprint(os.Stdout, msg)
}

type authRepository struct {
	credentialsPath string
	tokenPath       string
}

func NewAuthRepository(credentialsPath, tokenPath string) *authRepository {
	return &authRepository{
		credentialsPath: credentialsPath,
		tokenPath:       tokenPath,
	}
}

func (ar *authRepository) GenerateAuthURL() (string, error) {
	config, err := getConfig(ar.credentialsPath)
	if err != nil {
		return "", err
	}
	return config.AuthCodeURL("state-token", oauth2.AccessTypeOffline), nil
}

func (ar *authRepository) SaveAuthToken(authCode string) error {
	config, err := getConfig(ar.credentialsPath)
	if err != nil {
		return err
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return fmt.Errorf("unable to retrieve token from web %v", err)
	}
	return saveToken(ar.tokenPath, tok)
}

func getConfig(credentialsPath string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
}

func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %v", err)
	}
	defer func() { _ = f.Close() }()
	return json.NewEncoder(f).Encode(token)
}
