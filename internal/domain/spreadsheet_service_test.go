package domain_test

import (
	"errors"
	"habitsSync/internal/domain"
	"testing"
)

func TestSpreadsheet_Update(t *testing.T) {
	type fields struct {
		driveRepo  domain.DriveRepository
		sheetsRepo domain.SheetsRepository
	}
	type args struct {
		cmd domain.UpdateCMD
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update spreadsheet with passed habits",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: make([]domain.File, 1),
				},
				sheetsRepo: &fakeSheetRepo{
					createErr: nil,
					updateErr: nil,
				},
			},
			args: args{
				cmd: validUpdateCMD(),
			},
			wantErr: false,
		},
		{
			name:   "fail on invalid command",
			fields: fields{},
			args: args{
				cmd: invalidCMD(),
			},
			wantErr: true,
		},
		{
			name:   "fail on invalid command",
			fields: fields{},
			args: args{
				cmd: invalidCMD(),
			},
			wantErr: true,
		},
		{
			name:   "do nothing if no habits to be updated",
			fields: fields{},
			args: args{
				cmd: validCMDnoHabits(),
			},
			wantErr: false,
		},
		{
			name: "fail on finding spreadsheet error",
			fields: fields{
				driveRepo: fakeDriveRepo{
					err: errors.New("fake list error"),
				},
			},
			args: args{
				cmd: validUpdateCMD(),
			},
			wantErr: true,
		},
		{
			name: "fail when finding spreadsheets returns no results",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: nil,
				},
			},
			args: args{
				cmd: validUpdateCMD(),
			},
			wantErr: true,
		},
		{
			name: "fail when finding spreadsheets returns more than one result",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: make([]domain.File, 2),
				},
			},
			args: args{
				cmd: validUpdateCMD(),
			},
			wantErr: true,
		},
		{
			name: "fail on creation of Sheet Name error",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: make([]domain.File, 1),
				},
				sheetsRepo: &fakeSheetRepo{
					createErr: errors.New("fake create error"),
					updateErr: nil,
				},
			},
			args: args{
				cmd: validUpdateCMD(),
			},
			wantErr: true,
		},
		{
			name: "fail when updating Sheet with Habits and error",
			fields: fields{
				driveRepo: fakeDriveRepo{
					listResult: make([]domain.File, 1),
				},
				sheetsRepo: &fakeSheetRepo{
					createErr: nil,
					updateErr: errors.New("fake update error"),
				},
			},
			args: args{
				cmd: validUpdateCMD(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := domain.NewSpreadsheet(
				tt.fields.driveRepo,
				tt.fields.sheetsRepo,
			)
			if err := s.Update(tt.args.cmd); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCMD_Validate(t *testing.T) {
	type fields struct {
		Spreadsheet string
		SheetName   string
		Habits      []domain.Habit
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "fail on empty spreadsheet",
			fields: fields{
				Spreadsheet: "",
				SheetName:   "sheet name",
			},
			wantErr: true,
		},
		{
			name: "fail on empty sheet name",
			fields: fields{
				Spreadsheet: "spreadsheet",
				SheetName:   "",
			},
			wantErr: true,
		},
		{
			name: "valid command",
			fields: fields{
				Spreadsheet: "spreadsheet",
				SheetName:   "sheet name",
				Habits:      nil,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &domain.UpdateCMD{
				Spreadsheet: tt.fields.Spreadsheet,
				SheetName:   tt.fields.SheetName,
				Habits:      tt.fields.Habits,
			}
			if err := c.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func validUpdateCMD() domain.UpdateCMD {
	return domain.UpdateCMD{
		Spreadsheet: "spreadsheet",
		SheetName:   "sheetname",
		Habits: []domain.Habit{{
			ID:    1,
			Name:  "habit1",
			Count: 10,
		}},
	}
}

func validCMDnoHabits() domain.UpdateCMD {
	habit := validUpdateCMD()
	habit.Habits = nil
	return habit
}

func invalidCMD() domain.UpdateCMD {
	return domain.UpdateCMD{}
}
