package drive

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"
)

type repository struct {
	client *drive.Service
}

func NewRepository(credentialsPaths, tokenPath string) (*repository, error) {
	config, err := getConfig(credentialsPaths)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %v", err)
	}
	tok, err := tokenFromFile(tokenPath)
	if err != nil {
		return nil, err
	}
	client := config.Client(context.Background(), tok)

	srv, err := drive.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	return &repository{
		client: srv,
	}, nil
}

func getConfig(credentialsPath string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
}

type listResult struct {
	id   string
	name string
}

func (r *repository) ListByPrefix(contains string) ([]listResult, error) {
	if strings.Contains(contains, "'") {
		return nil, errors.New("prefix contains unsupported single quote character")
	}

	rsp, err := r.client.Files.List().
		Q(fmt.Sprintf("name contains '%v'", contains)).
		PageSize(30).
		Fields("nextPageToken, files(id, name)").
		Do()

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve files: %v", err)
	}

	lr := make([]listResult, 0)
	for _, r := range rsp.Files {
		lr = append(lr, listResult{
			id:   r.Id,
			name: r.Name,
		})
	}

	return lr, nil
}

func (r *repository) Download(id string) ([]byte, error) {
	rsp, err := r.client.Files.Get(id).Download()
	if err != nil {
		return nil, err
	}
	defer func() { _ = rsp.Body.Close() }()
	return ioutil.ReadAll(rsp.Body)
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}
