package drive

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"path"
	"time"
)

const allHabitsQuery = `select Habits.Id, Habits.name, count(Repetitions.id)
	from Habits
	left join Repetitions on Habits.Id = Repetitions.habit
	where not Habits.archived and Repetitions.value = 2
	and Repetitions.timestamp >= ? and Repetitions.timestamp <= ?
	group by Habits.Id
	order by Habits.name, Habits.Id`

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

type stat struct {
	ID    int
	name  string
	count int
}

func (d *storage) AllHabits(from, to time.Time) ([]stat, error) {
	result, err := d.db.Query(allHabitsQuery, from.Unix()*1000, to.Unix()*1000)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	stats := make([]stat, 0)
	for result.Next() {
		s := stat{}
		if err = result.Scan(&s.ID, &s.name, &s.count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
