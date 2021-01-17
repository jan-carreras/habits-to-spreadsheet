package application

import (
	"errors"
	"fmt"
	"habitsSync/internal/domain"
	"io"
	"time"
)

type SyncService struct {
	repository   domain.DriveRepository
	dbFile       domain.FileRepository
	storageMaker domain.StorageMaker
	output       io.Writer
}

func NewSyncService(
	r domain.DriveRepository,
	dbFile domain.FileRepository,
	storageMaker domain.StorageMaker,
	out io.Writer,
) *SyncService {
	return &SyncService{
		repository:   r,
		dbFile:       dbFile,
		storageMaker: storageMaker,
		output:       out,
	}
}

type SyncCMD struct {
	Prefix string
	From   time.Time
	To     time.Time
}

func (s *SyncService) Handle(cmd SyncCMD) error {
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

	if _, err = fmt.Fprintf(s.output, "Found %d backup files\n", len(res)); err != nil {
		return err
	}

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

	// TODO: OpenDB and start importing info to GoogleSheets
	// TODO: Read from sqlite3 file and map certain habits in between dates. Count the events that have happened,
	// 		 and update a google sheet with the result

	return nil
}

func (s *SyncService) openOrDownload(res domain.ListResult) (domain.Storage, error) {
	if s.dbFile.Exists(res.Name) {
		if _, err := fmt.Fprintf(s.output, "File already downloaded: '%v'\n", res.Name); err != nil {
			return nil, err
		}
		return s.storageMaker.Make(res.Name)
	}

	if _, err := fmt.Fprintf(s.output, "Downloading newest: '%v'\n", res.Name); err != nil {
		return nil, err
	}
	db, err := s.repository.Download(res.ID)
	if err != nil {
		return nil, err
	}
	if err := s.dbFile.Store(res.Name, db); err != nil {
		return nil, err
	}

	return s.storageMaker.Make(res.Name)
}
