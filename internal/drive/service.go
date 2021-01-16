package drive

import (
	"errors"
	"fmt"
	"io"
	"time"
)

type driveRepository interface {
	ListByPrefix(contains string) ([]listResult, error)
	Download(id string) ([]byte, error)
}

type fileRepository interface {
	Exists(name string) bool
	Store(name string, db []byte) error
}

type storageMaker interface {
	Make(name string) (*storage, error)
}

type Service struct {
	repository   driveRepository
	dbFile       fileRepository
	storageMaker storageMaker
	output       io.Writer
}

func NewService(r driveRepository, dbFile *dbFile, storageMaker storageMaker, out io.Writer) *Service {
	return &Service{
		repository:   r,
		dbFile:       dbFile,
		storageMaker: storageMaker,
		output:       out,
	}
}

type CMD struct {
	Prefix string
	From   time.Time
	To     time.Time
}

func (s *Service) Handle(cmd CMD) error {
	if cmd.From.IsZero() || cmd.To.IsZero() {
		return errors.New("from/to dates cannot be empty")
	}

	res, err := s.repository.ListByPrefix(cmd.Prefix)
	if err != nil {
		return err
	}

	if len(res) == 0 {
		return fmt.Errorf("no backup found with prefix '%v'", cmd.Prefix)
	}

	fmt.Fprintf(s.output, "Found %d backup files\n", len(res))

	storage, err := s.openOrDownload(res[0])
	if err != nil {
		return err
	}

	habits, err := storage.AllHabits(cmd.From, cmd.To)
	if err != nil {
		return err
	}
	for _, h := range habits {
		fmt.Println(h)
	}

	// OpenDB and start importing info to GoogleSheets

	// Read from sqlite3 file and map certain habits in between dates. Count the events that have happened,
	// and update a google sheet with the result

	return nil
}

func (s *Service) openOrDownload(res listResult) (*storage, error) {
	if s.dbFile.Exists(res.name) {
		fmt.Fprintf(s.output, "File already downloaded: '%v'\n", res.name)
		return s.storageMaker.Make(res.name)
	}

	fmt.Fprintf(s.output, "Downloading newest: '%v'\n", res.name)
	db, err := s.repository.Download(res.id)
	if err != nil {
		return nil, err
	}
	if err := s.dbFile.Store(res.name, db); err != nil {
		return nil, err
	}

	return s.storageMaker.Make(res.name)
}
