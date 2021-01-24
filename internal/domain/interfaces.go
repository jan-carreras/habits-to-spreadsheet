package domain

import (
	"time"
)

type DriveRepository interface {
	ListByPrefix(contains string) ([]File, error)
	Download(id string) ([]byte, error)
}

type FileRepository interface {
	Exists(name string) bool
	Store(name string, db []byte) error
}

type StorageMaker interface {
	Make(name string) (Storage, error)
}

type Storage interface {
	AllHabits(from, to time.Time) ([]Habit, error)
}

type SheetsRepository interface {
	CreateSheet(id string, name string) error
	UpdateSheet(id string, name string, stats []Habit) error
}
