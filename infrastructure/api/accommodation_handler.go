package api

import (
	"context"
	"encoding/json"
	"github.com/ZMS-DevOps/search-service/application"
	"github.com/ZMS-DevOps/search-service/domain"
	"github.com/ZMS-DevOps/search-service/infrastructure/dto"
	"github.com/ZMS-DevOps/search-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"log"
)

type AccommodationHandler struct {
	service       *application.AccommodationService
	traceProvider *sdktrace.TracerProvider
	loki          promtail.Client
}

func NewAccommodationHandler(service *application.AccommodationService, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *AccommodationHandler {
	server := &AccommodationHandler{
		service:       service,
		traceProvider: traceProvider,
		loki:          loki,
	}
	return server
}

func (handler *AccommodationHandler) OnRatingChanged(message *kafka.Message) {
	ctx := context.Background()
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "on-rating-changed")
	defer func() { span.End() }()
	var ratingChangedRequest dto.RatingChangedRequest
	if err := json.Unmarshal(message.Value, &ratingChangedRequest); err != nil {
		util.HttpTraceError(err, "to unmarshal data", span, handler.loki, "OnRatingChanged", "")
		log.Printf("Error unmarshalling rating change request: %v", err)
	}

	handler.service.OnCreateRatingChangeNotification(ratingChangedRequest, span, handler.loki)
}
