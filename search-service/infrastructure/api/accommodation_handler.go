package api

import (
	"encoding/json"
	"github.com/ZMS-DevOps/search-service/application"
	"github.com/ZMS-DevOps/search-service/infrastructure/dto"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"log"
)

type AccommodationHandler struct {
	service *application.AccommodationService
}

func NewAccommodationHandler(service *application.AccommodationService) *AccommodationHandler {
	server := &AccommodationHandler{
		service: service,
	}
	return server
}

func (handler *AccommodationHandler) OnRatingChanged(message *kafka.Message) {
	var ratingChangedRequest dto.RatingChangedRequest
	if err := json.Unmarshal(message.Value, &ratingChangedRequest); err != nil {
		log.Printf("Error unmarshalling rating change request: %v", err)
	}

	handler.service.OnCreateRatingChangeNotification(ratingChangedRequest)
}
