package application

import (
	"errors"
	"fmt"
	"time"
)

const dateLayout = "2006-01-02"

type quarters struct {
	start time.Time
	to    time.Time
}

type DatesService struct {
}

func NewDatesService() *DatesService {
	return &DatesService{}
}

type DatesCMD struct {
	From    string
	To      string
	Quarter int
	now     time.Time
}

type DatesOut struct {
	From time.Time
	To   time.Time
}

func (s *DatesService) Handle(cmd DatesCMD) (DatesOut, error) {
	cmd.now = time.Now()
	return s.parseDates(cmd)
}

func (s *DatesService) parseDates(cmd DatesCMD) (dates DatesOut, err error) {
	if (cmd.From != "" && cmd.To == "") || (cmd.From == "" && cmd.To != "") {
		err = errors.New("you must define both from and to")
		return
	}

	dates.From, err = s.parseDate(cmd.From)
	if err != nil {
		return
	}
	dates.To, err = s.parseDate(cmd.To)
	if err != nil {
		return
	}

	qs := makeQuarters(cmd.now)

	if cmd.Quarter != 0 {
		q := cmd.Quarter - 1 // 0 index access
		if q < 0 && q > len(qs) {
			err = fmt.Errorf("invalid quarter %v. valid quarters go from 1 to 4", cmd.Quarter)
			return
		}
		dates.From = qs[q].start
		dates.To = qs[q].to
	}

	// No dates defined, set the current quarter
	if dates.From.IsZero() || dates.To.IsZero() {
		for _, qt := range qs {
			if cmd.now.Sub(qt.start) >= 0 && cmd.now.Sub(qt.to) <= 0 {
				dates.From, dates.To = qt.start, qt.to
				break
			}
		}
	}
	return
}

func makeQuarters(now time.Time) []quarters {
	makeDate := dateMaker(now)
	return []quarters{
		{start: makeDate(1, 1, false), to: makeDate(3, 31, true)},
		{start: makeDate(4, 1, false), to: makeDate(6, 30, true)},
		{start: makeDate(7, 1, false), to: makeDate(9, 30, true)},
		{start: makeDate(10, 1, false), to: makeDate(12, 31, true)},
	}
}

func dateMaker(now time.Time) func(month time.Month, day int, end bool) time.Time {
	return func(month time.Month, day int, end bool) time.Time {
		hour, min, sec, nano := 0, 0, 0, 0
		if end {
			hour, min, sec, nano = 23, 59, 59, 999999999
		}
		return time.Date(now.Year(), month, day, hour, min, sec, nano, time.UTC)
	}
}

func (s *DatesService) parseDate(d string) (time.Time, error) {
	if d == "" {
		return time.Time{}, nil
	}

	return time.Parse(dateLayout, d)
}
