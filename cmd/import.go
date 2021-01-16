package main

import (
	"flag"
	"habitsSync/internal/auth"
	"habitsSync/internal/drive"
	"log"
	"os"
)

type args struct {
	credentialsPath string
	tokenPath       string
	tmpPath         string
	prefix          string
	authorize       bool
}

func parseArgs() (a args) {
	flag.StringVar(&a.credentialsPath, "credentials", "", "credentials path")
	flag.StringVar(&a.tokenPath, "token", "", "token path")
	flag.StringVar(&a.prefix, "prefix", "Loop Habits Backup", "prefix of the backup name")
	flag.StringVar(&a.tmpPath, "tmp", "/tmp", "Temporary directory where to store the DB")
	flag.BoolVar(&a.authorize, "auth", false, "authorize")
	flag.Parse()
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

	srv := drive.NewService(r,
		drive.NewDBFile(arg.tmpPath),
		drive.NewStorageFactory(arg.tmpPath),
		os.Stdout)

	err = srv.Handle(drive.CMD{
		Prefix: arg.prefix,
	})
	failOnErr(err)
}

func main() {
	arg := parseArgs()

	if arg.authorize {
		authorize(arg)
		return
	}

	importData(arg)
}

func failOnErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
