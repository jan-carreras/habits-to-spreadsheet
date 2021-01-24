package domain_test

import (
	"habitsSync/internal/domain"
	"time"
)

type fakeDriveRepo struct {
	listResult      []domain.File
	err             error
	errDownload     error
	payloadDownload []byte
}

func (f fakeDriveRepo) ListByPrefix(contains string) ([]domain.File, error) {
	return f.listResult, f.err
}

func (f fakeDriveRepo) Download(id string) ([]byte, error) {
	return f.payloadDownload, f.errDownload
}

type fakeFileRepo struct {
	exists     bool
	errOnStore error
}

func (f fakeFileRepo) Exists(name string) bool {
	return f.exists
}

func (f fakeFileRepo) Store(name string, db []byte) error {
	return f.errOnStore
}

type fakeStorageMaker struct {
	storage domain.Storage
	err     error
}

func (f fakeStorageMaker) Make(name string) (domain.Storage, error) {
	return f.storage, f.err
}

type fakeStorage struct {
	stats []domain.Habit
	err   error
}

func (f fakeStorage) AllHabits(from, to time.Time) ([]domain.Habit, error) {
	return f.stats, f.err
}

type fakeSheetRepo struct {
	createErr error
	updateErr error
}

func (f *fakeSheetRepo) CreateSheet(id string, name string) error {
	return f.createErr
}

func (f *fakeSheetRepo) UpdateSheet(id string, name string, stats []domain.Habit) error {
	return f.updateErr
}
