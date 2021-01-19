package application

import (
	"errors"
	"fmt"
	"habitsSync/internal/domain"
	"io"
	"time"
)

type SyncService struct {
	driveRepo    domain.DriveRepository
	sheetsRepo   domain.SheetsRepository
	dbFile       domain.FileRepository
	storageMaker domain.StorageMaker
	output       io.Writer
}

func NewSyncService(
	r domain.DriveRepository,
	s domain.SheetsRepository,
	db domain.FileRepository,
	sm domain.StorageMaker,
	out io.Writer,
) *SyncService {
	return &SyncService{
		driveRepo:    r,
		sheetsRepo:   s,
		dbFile:       db,
		storageMaker: sm,
		output:       out,
	}
}

type SyncCMD struct {
	Prefix      string
	From        time.Time
	To          time.Time
	Spreadsheet string
	SheetName   string
}

func (s *SyncService) Handle(cmd SyncCMD) error {
	if err := validateInput(cmd); err != nil {
		return err
	}

	res, err := s.driveRepo.ListByPrefix(cmd.Prefix)
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

	spreadsheetID, err := s.findSpreadsheet(cmd.Spreadsheet)
	if err != nil {
		return err
	}

	if _, err = fmt.Fprintf(s.output, "Importing %v habits...\n", len(habits)); err != nil {
		return err
	}

	if err := s.sheetsRepo.CreateSheet(spreadsheetID, cmd.SheetName); err != nil {
		return err
	}
	if err := s.sheetsRepo.UpdateSheet(spreadsheetID, cmd.SheetName, habits); err != nil {
		return err
	}

	if _, err = fmt.Fprint(s.output, "Habits imported successfully\n"); err != nil {
		return err
	}

	return nil
}

func (s *SyncService) findSpreadsheet(spreadsheet string) (string, error) {
	res, err := s.driveRepo.ListByPrefix(spreadsheet)
	if err != nil {
		return "", err
	}
	if len(res) == 0 {
		return "", fmt.Errorf("spreadsheet not found under the name of '%v'", spreadsheet)
	} else if len(res) > 1 {
		// TODO: We're searching by prefix so we could check if there is an exact match instead of aborting
		return "", fmt.Errorf("multiple spreadsheets found with same name: %v. Aborting", spreadsheet)
	}

	return res[0].ID, nil
}

func (s *SyncService) openOrDownload(res domain.ListResult) (domain.Storage, error) {
	if s.dbFile.Exists(res.Name) {
		if _, err := fmt.Fprintf(s.output, "Newest backup file already downloaded: '%v'\n", res.Name); err != nil {
			return nil, err
		}
		return s.storageMaker.Make(res.Name)
	}

	if _, err := fmt.Fprintf(s.output, "Downloading newest backup: '%v'\n", res.Name); err != nil {
		return nil, err
	}
	db, err := s.driveRepo.Download(res.ID)
	if err != nil {
		return nil, err
	}
	if err := s.dbFile.Store(res.Name, db); err != nil {
		return nil, err
	}

	return s.storageMaker.Make(res.Name)
}

func validateInput(cmd SyncCMD) error {
	if cmd.From.IsZero() || cmd.To.IsZero() {
		return errors.New("from/to dates cannot be empty")
	}
	if cmd.Prefix == "" {
		return errors.New("prefix cannot be empty")
	}

	if cmd.Spreadsheet == "" {
		return errors.New("spreadsheet cannot be empty")
	}

	if cmd.SheetName == "" {
		return errors.New("sheet name cannot be empty")
	}

	return nil
}
