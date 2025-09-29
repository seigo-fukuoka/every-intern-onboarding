package repository

type EventRepositoryInterface interface {
	GetAll() ([]Event, error)
	Create(event Event) (Event, error)
	ExistsByDateAndTitle(date, title string) (bool, error)
}
