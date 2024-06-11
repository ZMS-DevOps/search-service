package application

import (
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/search-service/application/external"
	"github.com/ZMS-DevOps/search-service/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
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
	if err != nil {
		return nil, err
	}

	availableIds, err := external.FilterAvailableAccommodation(service.bookingClient, getIds(accommodation), startTime, endTime)
	if err != nil {
		return nil, err
	}

	filteredAccommodation := filterAccommodationsByIds(accommodation, availableIds.AccommodationIds)
	return filteredAccommodation, nil
}

func (service *SearchService) GetByHostId(hostId string) ([]*domain.Accommodation, error) {
	return service.store.GetByHostId(hostId)
}

func (service *SearchService) MapToGetByHostIdResponse(accommodations []*domain.Accommodation) []domain.SearchResponse {
	var searchResponses []domain.SearchResponse
	for _, acc := range accommodations {
		priceType := acc.DefaultPrice.Type
		searchResponse := domain.SearchResponse{
			Id:         acc.Id,
			Name:       acc.Name,
			Location:   acc.Location,
			MainPhoto:  acc.MainPhoto,
			Rating:     acc.Rating,
			TotalPrice: acc.DefaultPrice.Price, // not displayed in get host by id
			UnitPrice:  acc.DefaultPrice.Price,
			PriceType:  priceType.String(),
		}
		searchResponses = append(searchResponses, searchResponse)
	}
	return searchResponses
}

func getIds(response []*domain.SearchResponse) []primitive.ObjectID {
	accommodationIDs := make([]primitive.ObjectID, len(response))
	for i, searchResponse := range response {
		accommodationIDs[i] = searchResponse.Id
	}
	return accommodationIDs
}

func filterAccommodationsByIds(accommodation []*domain.SearchResponse, availableIds []string) []*domain.SearchResponse {
	availableMap := make(map[primitive.ObjectID]bool)
	for _, id := range availableIds {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			log.Fatalf("Failed to convert string to ObjectID: %v", err)
		}
		availableMap[objectID] = true
	}

	var filteredAccommodations []*domain.SearchResponse
	for _, acc := range accommodation {
		if availableMap[acc.Id] {
			filteredAccommodations = append(filteredAccommodations, acc)
		}
	}
	return filteredAccommodations
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
