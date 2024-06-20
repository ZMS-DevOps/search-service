package application

import (
	"fmt"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/search-service/application/external"
	"github.com/ZMS-DevOps/search-service/domain"
	"github.com/ZMS-DevOps/search-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"go.mongodb.org/mongo-driver/bson/primitive"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

type SearchService struct {
	store         domain.AccommodationStore
	bookingClient booking.BookingServiceClient
	traceProvider *sdktrace.TracerProvider
	loki          promtail.Client
}

func NewSearchService(store domain.AccommodationStore, bookingClient booking.BookingServiceClient, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *SearchService {
	return &SearchService{
		store:         store,
		bookingClient: bookingClient,
		traceProvider: traceProvider,
		loki:          loki,
	}
}

func (service *SearchService) Get(id primitive.ObjectID, span trace.Span, loki promtail.Client) (*domain.Accommodation, error) {
	util.HttpTraceInfo("Fetching accommodation by id...", span, loki, "Get", "")
	accommodation, err := service.store.Get(id)
	if err != nil {
		util.HttpTraceError(err, "Error getting accommodation", span, service.loki, "Get", id.Hex())
		return nil, err
	}
	util.HttpTraceInfo("Accommodation retrieved successfully", span, service.loki, "Get", id.Hex())

	return accommodation, nil
}

func (service *SearchService) GetAll(span trace.Span, loki promtail.Client) ([]*domain.Accommodation, error) {
	util.HttpTraceInfo("Fetching accommodations...", span, loki, "GetAll", "")
	accommodations, err := service.store.GetAll()
	if err != nil {
		util.HttpTraceError(err, "Error getting all accommodations", span, service.loki, "GetAll", "")
		return nil, err
	}
	util.HttpTraceInfo("All accommodations retrieved successfully", span, service.loki, "GetAll", "")

	return accommodations, nil
}

func (service *SearchService) Search(location string, guestNumber int, startTime time.Time, endTime time.Time, minPrice float32, maxPrice float32, span trace.Span, loki promtail.Client) ([]*domain.SearchResponse, error) {
	util.HttpTraceInfo("Searching accommodations...", span, loki, "Search", "")
	accommodation, err := service.store.Search(location, guestNumber, startTime, endTime, minPrice, maxPrice)
	if err != nil {
		util.HttpTraceError(err, "Error searching accommodations", span, service.loki, "Search", location)
		return nil, err
	}

	var searchResponses []*domain.SearchResponse
	for _, acc := range accommodation {
		searchResponse := &domain.SearchResponse{
			Id:         acc.Id,
			HostId:     acc.HostId,
			Location:   acc.Location,
			Rating:     acc.Rating,
			Name:       acc.Name,
			TotalPrice: service.CalculateTotalPrice(acc, startTime, endTime, guestNumber, span, loki),
			UnitPrice:  acc.DefaultPrice.Price,
			PriceType:  acc.DefaultPrice.Type.String(),
		}
		searchResponses = append(searchResponses, searchResponse)
	}

	availableIds, err := external.FilterAvailableAccommodation(service.bookingClient, getIds(searchResponses), startTime, endTime, span, service.loki)
	if err != nil {
		util.HttpTraceError(err, "Error filtering available accommodations", span, service.loki, "Search", location)
		return nil, err
	}

	filteredAccommodation := service.filterAccommodationsByIds(searchResponses, availableIds.AccommodationIds, span)
	util.HttpTraceInfo("Search completed successfully", span, service.loki, "Search", location)

	return filteredAccommodation, nil
}

func (service *SearchService) GetByHostId(hostId string, span trace.Span, loki promtail.Client) ([]*domain.Accommodation, error) {
	util.HttpTraceInfo("Fetching accommodations by host id...", span, loki, "GetByHostId", "")
	accommodations, err := service.store.GetByHostId(hostId)
	if err != nil {
		util.HttpTraceError(err, "Error getting accommodations by host ID", span, service.loki, "GetByHostId", hostId)
		return nil, err
	}
	util.HttpTraceInfo("Accommodations by host ID retrieved successfully", span, service.loki, "GetByHostId", hostId)

	return accommodations, nil
}

func (service *SearchService) MapToGetByHostIdResponse(accommodations []*domain.Accommodation, span trace.Span) []domain.SearchResponse {
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
	util.HttpTraceInfo("Mapped accommodations to search responses successfully", span, service.loki, "MapToGetByHostIdResponse", "")

	return searchResponses
}

func (service *SearchService) CalculateTotalPrice(accommodation *domain.Accommodation, startDate time.Time, endDate time.Time, numberOfGuests int, span trace.Span, loki promtail.Client) float32 {
	var totalPrice = float32(0)
	days := int(endDate.Sub(startDate).Hours() / 24)
	coefficient := float32(1)
	if accommodation.DefaultPrice.Type == domain.PerGuest {
		coefficient = float32(numberOfGuests)
	}
	util.HttpTraceInfo("Calculating total price...", span, loki, "CalculateTotalPrice", "")
	for i := 0; i < days; i++ {
		date := startDate.Add(time.Hour * 24 * time.Duration(i))
		totalPrice += service.GetPriceForDate(date, accommodation.DefaultPrice.Price, accommodation.SpecialPrice, span) * coefficient
	}
	util.HttpTraceInfo("Calculated total price successfully", span, service.loki, "CalculateTotalPrice", accommodation.Id.Hex())

	return totalPrice
}

func (service *SearchService) GetPriceForDate(date time.Time, defaultPrice float32, specialPrices []domain.SpecialPrice, span trace.Span) float32 {
	for _, specialPrice := range specialPrices {
		if !date.Before(specialPrice.DateRange.Start) && date.Before(specialPrice.DateRange.End) {
			util.HttpTraceInfo("Retrieved special price for date", span, service.loki, "GetPriceForDate", fmt.Sprintf("%f", specialPrice.Price))
			return specialPrice.Price
		}
	}
	util.HttpTraceInfo("Retrieved default price for date", span, service.loki, "GetPriceForDate", fmt.Sprintf("%f", defaultPrice))

	return defaultPrice
}

func getIds(response []*domain.SearchResponse) []primitive.ObjectID {
	accommodationIDs := make([]primitive.ObjectID, len(response))
	for i, searchResponse := range response {
		accommodationIDs[i] = searchResponse.Id
	}
	return accommodationIDs
}

func (service *SearchService) filterAccommodationsByIds(accommodation []*domain.SearchResponse, availableIds []string, span trace.Span) []*domain.SearchResponse {
	availableMap := make(map[primitive.ObjectID]bool)
	for _, id := range availableIds {
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			util.HttpTraceError(err, "Failed to convert string to ObjectID", span, service.loki, "filterAccommodationsByIds", id)
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
	util.HttpTraceInfo("Filtered accommodations by IDs successfully", span, service.loki, "filterAccommodationsByIds", "")

	return filteredAccommodations
}
