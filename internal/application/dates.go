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
}

type DatesOut struct {
	From time.Time
	To   time.Time
}

func (s *DatesService) Handle(cmd DatesCMD) (DatesOut, error) {
	return s.parseDates(cmd.From, cmd.To, cmd.Quarter)
}

/** TODO Things to change about this code
- Methods are too long — should be defined in smaller functions
- Time management should be injected for testability purposes
- Some functions are verbose and repeat themselves (date creation). Could be more simple
- Adding unit tests, of course
*/

func (s *DatesService) parseDates(from, to string, quarter int) (dates DatesOut, err error) {
	if err = noDatesOrBothDatesRequired(from, to); err != nil {
		return
	}

	dates.From, err = s.parseDate(from)
	if err != nil {
		return
	}
	dates.To, err = s.parseDate(to)
	if err != nil {
		return
	}

	now := time.Now() // TODO: That's untestable. Should come from a repository
	qs := makeQuarters(now)

	if quarter != 0 {
		q := quarter - 1 // 0 index access
		if q < 0 && q > len(qs) {
			err = fmt.Errorf("invalid quarter %v. valid quarters go from 1 to 4", quarter)
			return
		}
		dates.From = qs[q].start
		dates.To = qs[q].to
	}

	// No dates defined, set the current quarter
	if dates.From.IsZero() || dates.To.IsZero() {
		for _, qt := range qs {
			if now.Sub(qt.start) >= 0 && now.Sub(qt.to) <= 0 {
				dates.From, dates.To = qt.start, qt.to
				break
			}
		}
	}
	return
}

func noDatesOrBothDatesRequired(from, to string) error {
	if (from != "" && to == "") || (from == "" && to != "") {
		return errors.New("you must define both from and to")
	}
	return nil
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
