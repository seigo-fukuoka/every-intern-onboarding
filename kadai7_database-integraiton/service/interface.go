package service

import "kadai7_database-integration/repository"

type EventServiceInterface interface {
	GetAllEvents() ([]repository.Event, error)
	ScrapeAndSaveEvents(limit int) ([]repository.Event, error)
}
