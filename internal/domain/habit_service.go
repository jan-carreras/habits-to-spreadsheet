package domain

import (
	"errors"
	"fmt"
	"time"
)

type Habits struct {
	fileRepo     FileRepository
	storageMaker StorageMaker
	driveRepo    DriveRepository
}

func NewHabits(
	f FileRepository,
	s StorageMaker,
	d DriveRepository) *Habits {
	return &Habits{
		fileRepo:     f,
		storageMaker: s,
		driveRepo:    d,
	}
}

type GetAllCMD struct {
	Prefix string
	From   time.Time
	To     time.Time
}

func (c *GetAllCMD) Validate() error {
	if c.Prefix == "" {
		return errors.New("prefix cannot be empty")
	}
	if c.From.IsZero() || c.To.IsZero() {
		return errors.New("from and to cannot be zero")
	}
	return nil
}

func (h *Habits) GetAll(cmd GetAllCMD) ([]Habit, error) {
	if err := cmd.Validate(); err != nil {
		return nil, err
	}

	files, err := h.driveRepo.ListByPrefix(cmd.Prefix)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no backup found with prefix '%v'", cmd.Prefix)
	}

	file := files[0]

	if !h.fileRepo.Exists(file.Name) {
		if err := h.download(file); err != nil {
			return nil, err
		}
	}
	storage, err := h.storageMaker.Make(file.Name)
	if err != nil {
		return nil, err
	}

	return storage.AllHabits(cmd.From, cmd.To)
}

func (h *Habits) download(res File) error {
	db, err := h.driveRepo.Download(res.ID)
	if err != nil {
		return err
	}
	if err := h.fileRepo.Store(res.Name, db); err != nil {
		return err
	}
	return nil
}
