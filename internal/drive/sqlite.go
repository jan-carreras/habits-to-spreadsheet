package drive

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"path"
)

type storageFactory struct {
	path string
}

func NewStorageFactory(path string) *storageFactory {
	return &storageFactory{
		path: path,
	}
}

func (s *storageFactory) Make(name string) (*storage, error) {
	return NewStorage(path.Join(s.path, name))
}

type storage struct {
	db *sql.DB
}

func NewStorage(path string) (*storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &storage{
		db: db,
	}, nil
}

type stats struct {
}

func (d *storage) Task(name string) ([]stats, error) {
	stmt, err := d.db.Prepare("select Id, name from Habits where name = ? order by Id desc limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var taskName string
	var taskID int
	err = stmt.QueryRow(name).Scan(&taskID, &taskName)
	if err != nil {
		return nil, err
	}

	fmt.Println("Name", taskName)
	fmt.Println("ID", taskID)

	// TODO: Do a sub-query to select all the datapoints for this given task and accumulate the information

	return nil, nil
}
