package application_test

import (
	"bytes"
	"errors"
	"habitsSync/internal/application"
	"habitsSync/internal/domain"
	"io"
	"testing"
	"time"
)

func TestSyncService_Handle(t *testing.T) {
	type fields struct {
		driveRepo    domain.DriveRepository
		sheetsRepo   domain.SheetsRepository
		dbFile       domain.FileRepository
		storageMaker domain.StorageMaker
		output       io.Writer
	}
	type args struct {
		cmd application.SyncCMD
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "Error when date From is not defined",
			fields: fields{},
			args: args{
				cmd: application.SyncCMD{
					// From is undefined
					To: time.Now(),
				},
			},
			wantErr: true,
		},
		{
			name:   "Error when date To is not defined",
			fields: fields{},
			args: args{
				cmd: application.SyncCMD{
					From: time.Now(),
					// To is undefined
				},
			},
			wantErr: true,
		},
		{
			name:   "Error when Prefix is empty",
			fields: fields{},
			args: args{
				cmd: application.SyncCMD{
					From:   time.Now(),
					To:     time.Now(),
					Prefix: "",
				},
			},
			wantErr: true,
		},
		{
			name:   "Error when Spreadsheet is empty",
			fields: fields{},
			args: args{
				cmd: application.SyncCMD{
					From:        time.Now(),
					To:          time.Now(),
					Prefix:      "prefix",
					Spreadsheet: "",
				},
			},
			wantErr: true,
		},
		{
			name:   "Error when SheetName is empty",
			fields: fields{},
			args: args{
				cmd: application.SyncCMD{
					From:        time.Now(),
					To:          time.Now(),
					Prefix:      "prefix",
					Spreadsheet: "spreadsheet",
					SheetName:   "",
				},
			},
			wantErr: true,
		},
		{
			name: "Error when driveRepo errors",
			fields: fields{
				driveRepo: fakeDriveRepo{
					err: errors.New("fake error"),
				},
			},
			args:    args{cmd: validCMD()},
			wantErr: true,
		},
		{
			name: "Error when driveRepo returns no results",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: make([]domain.ListResult, 0),
				},
			},
			args:    args{cmd: validCMD()},
			wantErr: true,
		},
		{
			name: "Error when driveRepo cannot download file",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: []domain.ListResult{
						{
							ID:   "1",
							Name: "document 1",
						},
					},
					errDownload: errors.New("fake download error"),
				},
				dbFile: fakeDBfile{
					exists: false,
				},
				output: bytes.NewBuffer(nil),
			},
			args:    args{cmd: validCMD()},
			wantErr: true,
		},
		{
			name: "Error when saving database",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: []domain.ListResult{
						{
							ID:   "1",
							Name: "document 1",
						},
					},
					payloadDownload: []byte("fake response"),
				},
				dbFile: fakeDBfile{
					exists: false,
					err:    errors.New("fake error saving to disk"),
				},
				output: bytes.NewBuffer(nil),
			},
			args:    args{cmd: validCMD()},
			wantErr: true,
		},
		{
			name: "Error when storage maker fails",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: []domain.ListResult{
						{
							ID:   "1",
							Name: "document 1",
						},
					},
					payloadDownload: []byte("fake response"),
				},
				dbFile: fakeDBfile{
					exists: false,
				},
				output: bytes.NewBuffer(nil),
				storageMaker: fakeStorageMaker{
					err: errors.New("fake storage maker error"),
				},
			},
			args:    args{cmd: validCMD()},
			wantErr: true,
		},
		{
			name: "Error when AllHabits fails",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: []domain.ListResult{
						{
							ID:   "1",
							Name: "document 1",
						},
					},
					payloadDownload: []byte("fake response"),
				},
				dbFile: fakeDBfile{
					exists: false,
				},
				output: bytes.NewBuffer(nil),
				storageMaker: fakeStorageMaker{
					storage: fakeStorage{
						err: errors.New("fake AllHabits error"),
					},
				},
			},
			args:    args{cmd: validCMD()},
			wantErr: true,
		},
	}

	// Test cases definition end here

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := application.NewSyncService(
				tt.fields.driveRepo,
				tt.fields.sheetsRepo,
				tt.fields.dbFile,
				tt.fields.storageMaker,
				tt.fields.output,
			)
			if err := s.Handle(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func validCMD() application.SyncCMD {
	return application.SyncCMD{
		From:        time.Now(),
		To:          time.Now(),
		Prefix:      "prefix",
		Spreadsheet: "spreadsheet",
		SheetName:   "sheet name",
	}
}

type fakeDBfile struct {
	exists bool
	err    error
}

func (f fakeDBfile) Exists(name string) bool {
	return f.exists
}

func (f fakeDBfile) Store(name string, db []byte) error {
	return f.err
}

type fakeDriveRepo struct {
	listResult      []domain.ListResult
	err             error
	errDownload     error
	payloadDownload []byte
}

func (f fakeDriveRepo) ListByPrefix(contains string) ([]domain.ListResult, error) {
	return f.listResult, f.err
}

func (f fakeDriveRepo) Download(id string) ([]byte, error) {
	return f.payloadDownload, f.errDownload
}

type fakeStorageMaker struct {
	storage domain.Storage
	err     error
}

func (f fakeStorageMaker) Make(name string) (domain.Storage, error) {
	return f.storage, f.err
}

type fakeStorage struct {
	stats []domain.Stat
	err   error
}

func (f fakeStorage) AllHabits(from, to time.Time) ([]domain.Stat, error) {
	return f.stats, f.err
}
