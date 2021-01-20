package application_test

import (
	"habitsSync/internal/application"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testTimeRepository struct{}

func (r *testTimeRepository) Now() time.Time {
	return time.Date(2021, 01, 05, 0, 0, 0, 0, time.UTC)
}

func TestDatesService_Handle(t *testing.T) {
	type args struct {
		cmd application.DatesCMD
	}
	tests := []struct {
		name    string
		args    args
		want    application.DatesOut
		wantErr bool
	}{
		{
			name: "Parse from and to",
			args: args{
				cmd: application.DatesCMD{
					From: "2021-01-01",
					To:   "2021-03-31",
				},
			},
			want: application.DatesOut{
				From: time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2021, 03, 31, 23, 59, 59, 999999999, time.UTC),
			},
		},
		{
			name: "Parse quarter",
			args: args{
				cmd: application.DatesCMD{
					Quarter: 1,
				},
			},
			want: application.DatesOut{
				From: time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2021, 03, 31, 23, 59, 59, 999999999, time.UTC),
			},
		},
		{
			name: "Default quarter",
			args: args{
				cmd: application.DatesCMD{
					// No dates nor quarter
				},
			},
			want: application.DatesOut{
				From: time.Date(2021, 01, 01, 0, 0, 0, 0, time.UTC),
				To:   time.Date(2021, 03, 31, 23, 59, 59, 999999999, time.UTC),
			},
		},
		{
			name: "Fail on invalid date From",
			args: args{
				cmd: application.DatesCMD{
					From: "2021-01-01111",
					To:   "2021-01-02",
				},
			},
			wantErr: true,
		},
		{
			name: "Fail on invalid date To",
			args: args{
				cmd: application.DatesCMD{
					From: "2021-01-01",
					To:   "2021-01-022",
				},
			},
			wantErr: true,
		},
		{
			name: "Fail on just From passed",
			args: args{
				cmd: application.DatesCMD{
					From: "2021-01-01",
				},
			},
			wantErr: true,
		},
		{
			name: "Fail on just To passed",
			args: args{
				cmd: application.DatesCMD{
					To: "2021-01-02",
				},
			},
			wantErr: true,
		},
		{
			name: "Fail on negative quarter",
			args: args{
				cmd: application.DatesCMD{
					Quarter: -1,
				},
			},
			wantErr: true,
		},
		{
			name: "Fail on invalid positive quarter",
			args: args{
				cmd: application.DatesCMD{
					Quarter: 5,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := application.NewDatesService(&testTimeRepository{})
			got, err := s.Handle(tt.args.cmd)

			if tt.wantErr {
				assert.NotNil(t, t, err)
				return
			}

			if !got.From.Equal(tt.want.From) || !got.To.Equal(tt.want.To) {
				t.Errorf("Handle() got = %v, want %v", got, tt.want)
			}
		})
	}
}
