package main

import (
	"flag"
	"habitsSync/internal/application"
	"habitsSync/internal/infrastructure/auth"
	"habitsSync/internal/infrastructure/drive"
	"log"
	"os"
	"time"
)

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

func parseDates(a *args) error {
	s := application.NewDatesService()
	cmd := application.DatesCMD{
		From:    a.fromStr,
		To:      a.toStr,
		Quarter: a.quarter,
	}
	d, err := s.Handle(cmd)
	if err != nil {
		return err
	}
	a.from = d.From
	a.to = d.To
	return nil
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

	failOnErr(parseDates(&a))

	return a
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

	s, err := sheets.NewRepository(arg.credentialsPath, arg.tokenPath)
	failOnErr(err)

	srv := application.NewSyncService(
		r,
		s,
		drive.NewDBFile(arg.tmpPath),
		drive.NewStorageFactory(arg.tmpPath),
		os.Stdout)

	err = srv.Handle(application.SyncCMD{
		Prefix:      arg.prefix,
		From:        arg.from,
		To:          arg.to,
		SheetName:   arg.sheetName,
		Spreadsheet: arg.spreadsheet,
	})
	failOnErr(err)
}

func failOnErr(err error) {
	if err != nil {
		var std = log.New(os.Stderr, "", log.LstdFlags)
		std.Fatal(err)
	}
}
