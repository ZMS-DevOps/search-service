package application

import (
	"github.com/ZMS-DevOps/search-service/domain"
	"github.com/ZMS-DevOps/search-service/infrastructure/dto"
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

func (service *AccommodationService) getById(accommodationId primitive.ObjectID) (*domain.Accommodation, error) {
	return service.store.Get(accommodationId)
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

func (service *AccommodationService) OnCreateRatingChangeNotification(ratingChangedRequest dto.RatingChangedRequest) {
	accommodationId, err := primitive.ObjectIDFromHex(ratingChangedRequest.AccommodationId)
	if err != nil {
		return
	}
	err = service.store.UpdateRating(accommodationId, ratingChangedRequest.Rating)
	if err != nil {
		return
	}
}
