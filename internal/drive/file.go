package drive

import (
	"io/ioutil"
	"os"
	"path"
)

type dbFile struct {
	path string
}

func NewDBFile(path string) *dbFile {
	return &dbFile{
		path: path,
	}
}

func (d *dbFile) Exists(name string) bool {
	_, err := os.Stat(path.Join(d.path, name))
	if err != nil {
		return false
	}
	return true
}

func (d *dbFile) Store(name string, db []byte) error {
	return ioutil.WriteFile(path.Join(d.path, name), db, 0600)
}
