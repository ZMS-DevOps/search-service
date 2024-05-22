package dto

import (
	"fmt"
	"github.com/ZMS-DevOps/search-service/domain"
	search "github.com/ZMS-DevOps/search-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func MapAccommodation(accommodationId primitive.ObjectID, accommodationDto *search.Accommodation) *domain.Accommodation {
	accommodation := &domain.Accommodation{
		Id:           accommodationId,
		Name:         accommodationDto.Name,
		Location:     accommodationDto.Location,
		MainPhoto:    accommodationDto.MainPhoto,
		GuestNumber:  *mapGuestNumber(int(accommodationDto.MinGuestNumber), int(accommodationDto.MaxGuestNumber)),
		DefaultPrice: *mapDefaultPrice(accommodationDto.DefaultPrice, accommodationDto.PriceType),
		SpecialPrice: *mapSpecialPrice(accommodationDto.SpecialPrice),
	}
	return accommodation
}

func mapSpecialPrice(prices []*search.SpecialPrice) *[]domain.SpecialPrice {
	if prices == nil {
		return nil
	}

	var specialPrices []domain.SpecialPrice
	layout := time.RFC3339
	for _, price := range prices {
		startDate, err := time.Parse(layout, price.StartDate)
		if err != nil {
			fmt.Println("Error parsing start time:", err)
			continue
		}
		endDate, err := time.Parse(layout, price.EndDate)
		if err != nil {
			fmt.Println("Error parsing end time:", err)
			continue
		}
		specialPrice := domain.SpecialPrice{
			Price: price.Price,
			DateRange: domain.DateRange{
				Start: startDate,
				End:   endDate,
			},
		}
		specialPrices = append(specialPrices, specialPrice)
	}

	return &specialPrices
}

func mapDefaultPrice(price float32, priceTypeName string) *domain.DefaultPrice {
	var priceType, _ = mapPriceType(&priceTypeName)

	return &domain.DefaultPrice{
		Price: price,
		Type:  *priceType,
	}
}

func mapPriceType(priceTypeName *string) (*domain.PricingType, error) {
	var priceType domain.PricingType
	switch *priceTypeName {
	case "PerApartmentUnit":
		priceType = domain.PerApartmentUnit
	case "PerGuest":
		priceType = domain.PerGuest
	default:
		return nil, fmt.Errorf("invalid pricing type: %s", *priceTypeName)
	}
	return &priceType, nil
}

func mapGuestNumber(min int, max int) *domain.GuestNumber {
	return &domain.GuestNumber{
		Min: min,
		Max: max,
	}
}
