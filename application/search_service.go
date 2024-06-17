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
	log.Printf("Before search")
	accommodation, err := service.store.Search(location, guestNumber, startTime, endTime, minPrice, maxPrice)
	searchResponses := []*domain.SearchResponse{}
	for _, accommodation := range accommodation {
		log.Printf("Search accommodation: %v", accommodation)
		log.Printf("Search accommodation2: %s", accommodation.DefaultPrice.Type)
		searchResponse := &domain.SearchResponse{
			Id:         accommodation.Id,
			HostId:     accommodation.HostId,
			Location:   accommodation.Location,
			Rating:     accommodation.Rating,
			Name:       accommodation.Name,
			TotalPrice: service.CalculateTotalPrice(accommodation, startTime, endTime, guestNumber),
			UnitPrice:  accommodation.DefaultPrice.Price,
			PriceType:  accommodation.DefaultPrice.Type.String(),
		}
		searchResponses = append(searchResponses, searchResponse)
		log.Printf("Search response: %v", searchResponse)
	}
	if err != nil {
		return nil, err
	}

	availableIds, err := external.FilterAvailableAccommodation(service.bookingClient, getIds(searchResponses), startTime, endTime)
	if err != nil {
		return nil, err
	}

	filteredAccommodation := filterAccommodationsByIds(searchResponses, availableIds.AccommodationIds)
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
			HostId:     acc.HostId,
			Name:       acc.Name,
			Location:   acc.Location,
			MainPhoto:  acc.MainPhoto,
			Rating:     acc.Rating,
			TotalPrice: acc.DefaultPrice.Price,
			UnitPrice:  acc.DefaultPrice.Price,
			PriceType:  priceType.String(),
		}
		searchResponses = append(searchResponses, searchResponse)
	}
	return searchResponses
}

func (service *SearchService) CalculateTotalPrice(accommodation *domain.Accommodation, startDate time.Time, endDate time.Time, numberOfGuests int) float32 {
	var totalPrice = float32(0)
	days := int(endDate.Sub(startDate).Hours() / 24)
	coefficient := float32(1)
	if accommodation.DefaultPrice.Type == domain.PerGuest {
		coefficient = float32(numberOfGuests)
	}
	for i := 0; i < days; i++ {
		date := startDate.Add(time.Hour * 24 * time.Duration(i))
		totalPrice += service.GetPriceForDate(date, accommodation.DefaultPrice.Price, accommodation.SpecialPrice) * coefficient
		log.Printf("Total price: %f", totalPrice)
	}
	return totalPrice
}

func (service *SearchService) GetPriceForDate(date time.Time, defaultPrice float32, specialPrices []domain.SpecialPrice) float32 {
	log.Printf("Search: %v", date)
	for _, specialPrice := range specialPrices {
		log.Printf("Search accommodation: %v", specialPrice)
		if !date.Before(specialPrice.DateRange.Start) && date.Before(specialPrice.DateRange.End) {
			return specialPrice.Price
		}
	}
	return defaultPrice
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
