package domain

import (
	"time"
)

type StorageMaker interface {
	Make(name string) (Storage, error)
}

type Storage interface {
	AllHabits(from, to time.Time) ([]Stat, error)
}
