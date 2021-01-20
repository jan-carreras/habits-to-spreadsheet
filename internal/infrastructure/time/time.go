package time

import "time"

type Repository struct {
}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) Now() time.Time {
	return time.Now()
}
