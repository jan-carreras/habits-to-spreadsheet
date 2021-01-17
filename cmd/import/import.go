package main

import (
	"errors"
	"flag"
	"fmt"
	"habitsSync/internal/application"
	"habitsSync/internal/infrastructure/auth"
	"habitsSync/internal/infrastructure/drive"
	"log"
	"os"
	"time"
)

const dateLayout = "2006-01-02"

var now = time.Now()
var qs = []struct {
	start time.Time
	to    time.Time
}{
	{start: makeDate(1, 1, false), to: makeDate(3, 31, true)},
	{start: makeDate(4, 1, false), to: makeDate(6, 30, true)},
	{start: makeDate(7, 1, false), to: makeDate(9, 30, true)},
	{start: makeDate(10, 1, false), to: makeDate(12, 31, true)},
}

type args struct {
	credentialsPath string
	tokenPath       string
	tmpPath         string
	prefix          string
	authorize       bool
	fromStr         string
	toStr           string
	quarter         int
	from            time.Time
	to              time.Time
}

func parseDates(a *args) {
	if (a.fromStr != "" && a.toStr == "") || (a.fromStr == "" && a.toStr != "") {
		failOnErr(errors.New("you must define both from and to"))
	}

	a.from = parseDate(a.fromStr)
	a.to = parseDate(a.toStr)

	if a.quarter != 0 {
		q := a.quarter - 1 // 0 index access
		if q < 0 && q > len(qs) {
			failOnErr(fmt.Errorf("invalid quarter %v. valid quarters go from 1 to 4", a.quarter))
		}
		a.from = qs[q].start
		a.to = qs[q].to
	}

	// No dates defined, set the current quarter
	if a.from.IsZero() || a.to.IsZero() {
		for _, qt := range qs {
			if now.Sub(qt.start) >= 0 && now.Sub(qt.to) <= 0 {
				a.from, a.to = qt.start, qt.to
				break
			}
		}
	}
}

func parseArgs() (a args) {
	flag.StringVar(&a.credentialsPath, "credentials", "credentials.json", "credentials file")
	flag.StringVar(&a.tokenPath, "token", "auth.json", "token file")
	flag.StringVar(&a.prefix, "prefix", "Loop Habits Backup", "prefix of the backup name")
	flag.StringVar(&a.tmpPath, "tmp", "/tmp", "temporary directory where to store the DB")
	flag.StringVar(&a.fromStr, "from", "", "yyyy-mm-dd date from where start importing Habits records")
	flag.StringVar(&a.toStr, "to", "", "yyy-mm-dd date from where stop importing Habits records")
	flag.BoolVar(&a.authorize, "auth", false, "authorize")
	flag.IntVar(&a.quarter, "quarter", 0, "date range for the quarter of the current year")
	flag.Parse()

	parseDates(&a)

	return a
}

func parseDate(d string) time.Time {
	if d == "" {
		return time.Time{}
	}

	t, err := time.Parse(dateLayout, d)
	failOnErr(err)
	return t
}

func makeDate(month time.Month, day int, end bool) time.Time {
	hour, min, sec, nano := 0, 0, 0, 0
	if end {
		hour, min, sec, nano = 23, 59, 59, 999999999
	}
	return time.Date(now.Year(), month, day, hour, min, sec, nano, time.UTC)
}

func authorize(arg args) {
	service := auth.NewService(
		auth.NewReadWriter(),
		auth.NewAuthRepository(arg.credentialsPath, arg.tokenPath),
	)
	err := service.Handle()
	failOnErr(err)
}

func importData(arg args) {
	r, err := drive.NewRepository(arg.credentialsPath, arg.tokenPath)
	failOnErr(err)

	srv := application.NewSyncService(r,
		drive.NewDBFile(arg.tmpPath),
		drive.NewStorageFactory(arg.tmpPath),
		os.Stdout)

	err = srv.Handle(application.SyncCMD{
		Prefix: arg.prefix,
		From:   arg.from,
		To:     arg.to,
	})
	failOnErr(err)
}

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
