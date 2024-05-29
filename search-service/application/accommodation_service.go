package application

import (
	"github.com/ZMS-DevOps/search-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	return service.store.InsertWithId(&accommodation)
}

func (service *AccommodationService) GetAll() ([]*domain.Accommodation, error) {
	return service.store.GetAll()
}

func (service *AccommodationService) EditAccommodation(accommodation domain.Accommodation) error {
	return service.store.Update(accommodation.Id, &accommodation)
}

func (service *AccommodationService) DeleteAccommodation(id primitive.ObjectID) interface{} {
	return service.store.Delete(id)
}
