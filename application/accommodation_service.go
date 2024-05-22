package application

import (
	"github.com/ZMS-DevOps/search-service/domain"
)

type AccommodationService struct {
	store domain.AccommodationStore
}

func NewAccommodationService(store domain.AccommodationStore) *AccommodationService {
	return &AccommodationService{
		store: store,
	}
}

func (service *AccommodationService) AddAccommodation(accommodation domain.Accommodation) error {
	return service.store.Insert(&accommodation)
}

func (service *AccommodationService) GetAll() ([]*domain.Accommodation, error) {
	return service.store.GetAll()
}
