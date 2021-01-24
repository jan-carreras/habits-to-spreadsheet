# Habits to Spreadsheet

Simple Go application that imports habits recorded by
the [Loop Habit Tracker](https://play.google.com/store/apps/details?id=org.isoron.uhabits&hl=en&gl=US)
into a Google Spreadsheet.

## Show off code

I'm using this project to show "how do I code in Go", for whoever may be interested.

Interesting stuff:

* [Hexagonal architecture with SOLID principles](https://github.com/jan-carreras/habits-to-spreadsheet/tree/master/internal)
* Unit tests
  ([example](https://github.com/jan-carreras/habits-to-spreadsheet/blob/master/internal/domain/habit_service_test.go))
* Separation of concerns
* Simple [Makefile](https://github.com/jan-carreras/habits-to-spreadsheet/blob/master/Makefile)
* [Example of GIT usage](https://github.com/jan-carreras/habits-to-spreadsheet/commits/master) (commit size, signed
  commits, lah-dih-dah)
* [Minimal dependencies](https://github.com/jan-carreras/habits-to-spreadsheet/blob/master/go.mod)
* [TODO.md](https://github.com/jan-carreras/habits-to-spreadsheet/blob/master/TODO.md) driven development
* Proper [MIT license](https://github.com/jan-carreras/habits-to-spreadsheet/blob/master/LICENSE)

## Build

```bash
git clone https://github.com/jan-carreras/habits-to-spreadsheet
cd habits-to-spreadsheet
make
```

You will find a binary in `bin/hsync`

## Usage

The aim is to import a Habits backup into Google Spreadsheets. Let's start by creating a backup from the Habits app and
store it to the Drive. Leave the default name.

```bash
bin/hsync -spreadsheet "2021 - OKRs"
```

The spreadsheet must exist. A new Sheet called "Import" will be created with your habits, and it's count. The habits are
filtered by default by quarter. Use help to modify that or any other option:

```bash
bin/hsync -h
Usage of bin/hsync:
  -auth
        authorize
  -credentials string
        credentials file (default "credentials.json")
  -from string
        yyyy-mm-dd date from where start importing Habits records
  -prefix string
        prefix of the backup name (default "Loop Habits Backup")
  -quarter int
        date range for the quarter of the current year
  -sheet-name string
        the name of the Sheet where data is going to be imported (default "Import")
  -spreadsheet string
        name of the spreadsheet to import
  -tmp string
        temporary directory where to store the DB (default "/tmp")
  -to string
        yyy-mm-dd date from where stop importing Habits records
  -token string
        token file (default "auth.json")
```

## Contributing

Since the code is to show off, and I can hardly imagine anyone using it let alone contributing to the project, I don't
see much value in accepting PRs. _But_ feel free to create issues to highlight problems or requesting new features.

## License

[MIT](https://choosealicense.com/licenses/mit/)