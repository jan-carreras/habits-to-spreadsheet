package application_test

import "habitsSync/internal/domain"

type fakeHabitsGetter struct {
	habits []domain.Habit
	err    error
}

func (f *fakeHabitsGetter) GetAll(cmd domain.GetAllCMD) ([]domain.Habit, error) {
	return f.habits, f.err
}

type fakeSpreadsheetUpdater struct {
	err error
}

func (f *fakeSpreadsheetUpdater) Update(cmd domain.UpdateCMD) error {
	return f.err
}
