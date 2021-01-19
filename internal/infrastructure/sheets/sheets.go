package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"habitsSync/internal/domain"
	"io/ioutil"
	"log"
	"os"
)

type repository struct {
	client *sheets.Service
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

	srv, err := sheets.NewService(context.Background(), option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return &repository{client: srv}, nil
}

// TODO: This belongs to the Drive, not to Sheets. Can we use what we already have?
func (r *repository) FindDocument(name string) (string, error) {
	return "", nil
}

// TODO: This belongs to the Drive, not to Sheets. Can we use what we already have?
func (r *repository) SheetExists(id string, name string) (bool, error) {
	return false, nil
}

func (r *repository) CreateSheet(id string, name string) error {
	sheetAlreadyExists := func(s *sheets.Spreadsheet, name string) bool {
		for _, sh := range s.Sheets {
			if sh.Properties.Title == name {
				return true
			}
		}
		return false
	}
	createSheet := func(name string) error {
		req := sheets.Request{
			AddSheet: &sheets.AddSheetRequest{
				Properties: &sheets.SheetProperties{
					Title: name,
				},
			},
		}

		rbb := &sheets.BatchUpdateSpreadsheetRequest{
			Requests: []*sheets.Request{&req},
		}

		_, err := r.client.Spreadsheets.BatchUpdate(id, rbb).Do()
		return err
	}

	s, err := r.client.Spreadsheets.Get(id).Do()
	if err != nil {
		return err
	}

	if sheetAlreadyExists(s, name) {
		return nil
	}

	return createSheet(name)
}

func (r *repository) UpdateSheet(id string, name string, stats []domain.Stat) error {
	rows := make([][]interface{}, 0)
	rows = append(rows, []interface{}{"ID", "Name", "Count"})
	for _, st := range stats {
		rows = append(rows, []interface{}{st.ID, st.Name, st.Count})
	}

	rb := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "USER_ENTERED",
	}
	rb.Data = append(rb.Data, &sheets.ValueRange{
		Range:  fmt.Sprintf("%v!A2", name),
		Values: rows,
	})

	_, err := r.client.Spreadsheets.Values.BatchUpdate(id, rb).Do()
	if err != nil {
		return err
	}
	return nil
}

func getConfig(credentialsPath string) (*oauth2.Config, error) {
	b, err := ioutil.ReadFile(credentialsPath)
	if err != nil {
		return nil, err
	}
	return google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope, drive.DriveReadonlyScope)
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
