package application

import (
	"github.com/ZMS-DevOps/search-service/domain"
	//"github.com/ZMS-DevOps/search-service/infrastructure/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

func (service *SearchService) Search(location string, guestNumber int, startTime time.Time, endTime time.Time, minPrice float32, maxPrice float32) ([]*domain.SearchResponse, error) {
	accommodation, err := service.store.Search(location, guestNumber, startTime, endTime, minPrice, maxPrice)

	return accommodation, err
}

//func CalculateTotalPrice(accommodation domain.Accommodation, start, end time.Time, numPeople int) float32 {
//	totalPrice := float32(0)
//	numDays := int(end.Sub(start).Hours() / 24)
//	for i := 0; i < numDays; i++ {
//		currentDate := start.AddDate(0, 0, i)
//		dailyPrice := accommodation.DefaultPrice.Price
//
//		for _, sp := range accommodation.SpecialPrice {
//			if currentDate.After(sp.DateRange.Start) && currentDate.Before(sp.DateRange.End.AddDate(0, 0, 1)) {
//				dailyPrice = sp.Price
//				break
//			}
//		}
//
//		if accommodation.DefaultPrice.Type == PerGuest {
//			dailyPrice *= float32(numPeople)
//		}
//
//		totalPrice += dailyPrice
//	}
//
//	return totalPrice
//}
