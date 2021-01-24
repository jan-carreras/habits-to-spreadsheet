package domain_test

import (
	"errors"
	"habitsSync/internal/domain"
	"reflect"
	"testing"
	"time"
)

func TestGetAllCMD_Validate(t *testing.T) {
	type fields struct {
		Prefix string
		From   time.Time
		To     time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "all args are ok",
			fields: fields{
				Prefix: "prefix",
				From:   time.Now(),
				To:     time.Now(),
			},
			wantErr: false,
		},
		{
			name: "fail on empty prefix",
			fields: fields{
				Prefix: "",
			},
			wantErr: true,
		},
		{
			name: "fail on empty start date",
			fields: fields{
				Prefix: "prefix",
				// From: empty
				To: time.Now(),
			},
			wantErr: true,
		},
		{
			name: "fail on empty end date",
			fields: fields{
				Prefix: "prefix",
				From:   time.Now(),
				// To: empty
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &domain.GetAllCMD{
				Prefix: tt.fields.Prefix,
				From:   tt.fields.From,
				To:     tt.fields.To,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHabits_GetAll(t *testing.T) {
	type fields struct {
		fileRepo     domain.FileRepository
		storageMaker domain.StorageMaker
		driveRepo    domain.DriveRepository
	}
	type args struct {
		cmd domain.GetAllCMD
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Habit
		wantErr bool
	}{
		{
			name: "fail when invalid args are passed",
			fields: fields{
				driveRepo: fakeDriveRepo{
					err: errors.New("fake list error"),
				},
			},
			args: args{
				cmd: domain.GetAllCMD{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail when Prefix not found in storage",
			fields: fields{
				driveRepo: fakeDriveRepo{
					err: errors.New("fake list error"),
				},
			},
			args: args{
				cmd: validCMD(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail when storage returns no results",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: make([]domain.File, 0),
				},
			},
			args: args{
				cmd: validCMD(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail when DB cannot be downloaded",
			fields: fields{
				fileRepo: fakeFileRepo{
					errOnStore: errors.New("fake file error"),
				},
				driveRepo: fakeDriveRepo{
					listResult: []domain.File{
						{ID: "1", Name: "file"},
					},
					errDownload: errors.New("fake error download"),
				},
			},
			args: args{
				cmd: validCMD(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail when DB cannot be stored",
			fields: fields{
				fileRepo: fakeFileRepo{
					errOnStore: errors.New("fake file error"),
				},
				driveRepo: fakeDriveRepo{
					listResult: []domain.File{
						{ID: "1", Name: "file"},
					},
				},
			},
			args: args{
				cmd: validCMD(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail to open the storage file",
			fields: fields{
				fileRepo: fakeFileRepo{},
				storageMaker: fakeStorageMaker{
					err: errors.New("fake unable to open storage"),
				},
				driveRepo: fakeDriveRepo{
					listResult: []domain.File{
						{ID: "1", Name: "file"},
					},
				},
			},
			args: args{
				cmd: validCMD(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "fail when listing the habits fail",
			fields: fields{
				fileRepo: fakeFileRepo{},
				storageMaker: fakeStorageMaker{
					storage: fakeStorage{
						stats: nil,
						err:   errors.New("fake listing habits failure"),
					},
				},
				driveRepo: fakeDriveRepo{
					listResult: []domain.File{
						{ID: "1", Name: "file"},
					},
				},
			},
			args: args{
				cmd: validCMD(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "get the habits",
			fields: fields{
				fileRepo: fakeFileRepo{},
				storageMaker: fakeStorageMaker{
					storage: fakeStorage{
						stats: make([]domain.Habit, 0),
					},
				},
				driveRepo: fakeDriveRepo{
					listResult: []domain.File{
						{ID: "1", Name: "file"},
					},
				},
			},
			args: args{
				cmd: validCMD(),
			},
			want:    make([]domain.Habit, 0),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := domain.NewHabits(
				tt.fields.fileRepo,
				tt.fields.storageMaker,
				tt.fields.driveRepo,
			)
			got, err := h.GetAll(tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func validCMD() domain.GetAllCMD {
	return domain.GetAllCMD{
		From:   time.Now(),
		To:     time.Now(),
		Prefix: "prefix",
	}
}
