package domain

import (
	"errors"
	"fmt"
)

type Spreadsheet struct {
	driveRepo  DriveRepository
	sheetsRepo SheetsRepository
}

func NewSpreadsheet(
	d DriveRepository,
	s SheetsRepository) *Spreadsheet {
	return &Spreadsheet{
		driveRepo:  d,
		sheetsRepo: s,
	}
}

type UpdateCMD struct {
	Spreadsheet string
	SheetName   string
	Habits      []Habit
}

func (c *UpdateCMD) Validate() error {
	if c.Spreadsheet == "" {
		return errors.New("spreadsheet cannot be empty")
	}
	if c.SheetName == "" {
		return errors.New("sheet name cannot be empty")
	}
	return nil
}

func (s *Spreadsheet) Update(cmd UpdateCMD) error {
	if err := cmd.Validate(); err != nil {
		return err
	}

	if len(cmd.Habits) == 0 {
		return nil // Nothing to update
	}

	spreadsheetID, err := s.findSpreadsheet(cmd.Spreadsheet)
	if err != nil {
		return err
	}

	if err := s.sheetsRepo.CreateSheet(spreadsheetID, cmd.SheetName); err != nil {
		return err
	}
	if err := s.sheetsRepo.UpdateSheet(spreadsheetID, cmd.SheetName, cmd.Habits); err != nil {
		return err
	}

	return nil
}

func (s *Spreadsheet) findSpreadsheet(spreadsheet string) (string, error) {
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
