package application

import (
	"errors"
	"fmt"
	"time"
)

type TimeRepository interface {
	Now() time.Time
}

const dateLayout = "2006-01-02"

type quarters struct {
	from time.Time
	to   time.Time
}

type DatesService struct {
	timeRepository TimeRepository
}

func NewDatesService(t TimeRepository) *DatesService {
	return &DatesService{
		timeRepository: t,
	}
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
	dates.To = endOfDay(dates.To)

	now := s.timeRepository.Now()
	qs := makeQuarters(now.Year())

	if quarter != 0 {
		q := quarter - 1 // 0 index access
		if q < 0 || q >= len(qs) {
			err = fmt.Errorf("invalid quarter %v. valid quarters go from 1 to 4", quarter)
			return
		}
		dates.From = qs[q].from
		dates.To = qs[q].to
		return
	}

	// No dates defined, set the current quarter
	if dates.From.IsZero() || dates.To.IsZero() {
		for _, qt := range qs {
			if now.Sub(qt.from) >= 0 && now.Sub(qt.to) <= 0 {
				dates.From, dates.To = qt.from, qt.to
				return
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

func makeQuarters(year int) []quarters {
	mkDate := func(month time.Month, day int) time.Time {
		return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	}

	qs := []quarters{
		{from: mkDate(1, 1), to: mkDate(3, 31)},
		{from: mkDate(4, 1), to: mkDate(6, 30)},
		{from: mkDate(7, 1), to: mkDate(9, 30)},
		{from: mkDate(10, 1), to: mkDate(12, 31)},
	}

	for i := range qs {
		qs[i].to = endOfDay(qs[i].to)
	}

	return qs
}

func (s *DatesService) parseDate(d string) (time.Time, error) {
	if d == "" {
		return time.Time{}, nil
	}

	return time.Parse(dateLayout, d)
}

func endOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}
