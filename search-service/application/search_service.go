package application

import (
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/search-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type SearchService struct {
	store         domain.AccommodationStore
	bookingClient booking.BookingServiceClient
}

func NewSearchService(store domain.AccommodationStore, bookingClient booking.BookingServiceClient) *SearchService {
	return &SearchService{
		store:         store,
		bookingClient: bookingClient,
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
	accommodation, err = service.bookingClient.FilterAvailableAccommodation(getIds(accommodation), startTime, endTime)
	return accommodation, err
}

func getIds(response []domain.SearchResponse) {
	accommodationIDs := make([]string, len(response))
	for i, searchResponse := range response {
		accommodationIDs[i] = searchResponse.Id.Hex()
	}
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
