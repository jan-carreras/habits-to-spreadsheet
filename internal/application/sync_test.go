package application_test

import (
	"errors"
	"habitsSync/internal/application"
	"habitsSync/internal/domain"
	"io"
	"io/ioutil"
	"testing"
)

func TestSyncService_Handle(t *testing.T) {
	type fields struct {
		habitsGetter       application.HabitsGetter
		spreadsheetUpdater application.SpreadsheetUpdater
		output             io.Writer
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
			name: "fail when getting all the habits fail",
			fields: fields{
				habitsGetter: &fakeHabitsGetter{
					err: errors.New("fake get all error"),
				},
				output: ioutil.Discard,
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "fail when updating spreadsheets with habits returns error",
			fields: fields{
				habitsGetter: &fakeHabitsGetter{
					habits: []domain.Habit{
						{
							ID:    1,
							Name:  "habit 1",
							Count: 10,
						},
					},
				},
				spreadsheetUpdater: &fakeSpreadsheetUpdater{
					err: errors.New("fake update error"),
				},
				output: ioutil.Discard,
			},
			args:    args{},
			wantErr: true,
		},
		{
			name: "get all habits and update spreadsheet",
			fields: fields{
				habitsGetter: &fakeHabitsGetter{
					habits: []domain.Habit{
						{
							ID:    1,
							Name:  "habit 1",
							Count: 10,
						},
					},
				},
				spreadsheetUpdater: &fakeSpreadsheetUpdater{},
				output:             ioutil.Discard,
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := application.NewSyncService(
				tt.fields.habitsGetter,
				tt.fields.spreadsheetUpdater,
				tt.fields.output,
			)
			if err := s.Handle(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Handle() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
