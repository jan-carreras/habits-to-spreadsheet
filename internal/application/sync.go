package application

import (
	"fmt"
	"habitsSync/internal/domain"
	"io"
	"time"
)

type HabitsGetter interface {
	GetAll(cmd domain.GetAllCMD) ([]domain.Habit, error)
}

type SpreadsheetUpdater interface {
	Update(cmd domain.UpdateCMD) error
}

type SyncService struct {
	habitsGetter       HabitsGetter
	spreadsheetUpdater SpreadsheetUpdater
	output             io.Writer
}

func NewSyncService(
	h HabitsGetter,
	su SpreadsheetUpdater,
	out io.Writer,
) *SyncService {
	return &SyncService{
		habitsGetter:       h,
		spreadsheetUpdater: su,
		output:             out,
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
	habits, err := s.habitsGetter.GetAll(domain.GetAllCMD{
		From: cmd.From,
		To:   cmd.To,
	})
	if err != nil {
		return err
	}

	if _, err = fmt.Fprintf(s.output, "Importing %v habits...\n", len(habits)); err != nil {
		return err
	}

	updateCMD := domain.UpdateCMD{
		Spreadsheet: cmd.Spreadsheet,
		SheetName:   cmd.SheetName,
		Habits:      habits,
	}
	if err := s.spreadsheetUpdater.Update(updateCMD); err != nil {
		return err
	}

	if _, err = fmt.Fprint(s.output, "Habits imported successfully\n"); err != nil {
		return err
	}

	return nil
}
