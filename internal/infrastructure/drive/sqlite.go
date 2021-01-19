package drive

import (
	"database/sql"
	"habitsSync/internal/domain"
	"path"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

func (s *storageFactory) Make(name string) (domain.Storage, error) {
	return NewStorage(path.Join(s.path, name))
}

type Storage struct {
	db *sql.DB
}

func NewStorage(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	return &Storage{
		db: db,
	}, nil
}

func (d *Storage) AllHabits(from, to time.Time) ([]domain.Stat, error) {
	result, err := d.db.Query(allHabitsQuery, from.Unix()*1000, to.Unix()*1000)
	if err != nil {
		return nil, err
	}
	defer func() { _ = result.Close() }()

	stats := make([]domain.Stat, 0)
	for result.Next() {
		s := domain.Stat{}
		if err = result.Scan(&s.ID, &s.Name, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}
