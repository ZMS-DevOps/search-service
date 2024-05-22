package application

import (
	"github.com/ZMS-DevOps/search-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchService struct {
	store domain.AccommodationStore
}

func NewSearchService(store domain.AccommodationStore) *SearchService {
	return &SearchService{
		store: store,
	}
}

func (service *SearchService) Get(id primitive.ObjectID) (*domain.Accommodation, error) {
	return service.store.Get(id)
}

func (service *SearchService) GetAll() ([]*domain.Accommodation, error) {
	return service.store.GetAll()
}
