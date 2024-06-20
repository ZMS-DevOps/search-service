package application

import (
	"github.com/ZMS-DevOps/search-service/domain"
	"github.com/ZMS-DevOps/search-service/infrastructure/dto"
	"github.com/ZMS-DevOps/search-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

type AccommodationService struct {
	store domain.AccommodationStore
	loki  promtail.Client
}

func NewAccommodationService(store domain.AccommodationStore, loki promtail.Client) *AccommodationService {
	return &AccommodationService{
		store: store,
		loki:  loki,
	}
}

func (service *AccommodationService) AddAccommodation(accommodation domain.Accommodation, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Adding accommodation...", span, loki, "AddAccommodation", "")
	err := service.store.InsertWithId(&accommodation)
	if err != nil {
		util.HttpTraceError(err, "Error adding accommodation", span, service.loki, "AddAccommodation", accommodation.Id.Hex())
		return err
	}
	util.HttpTraceInfo("Accommodation added successfully", span, service.loki, "AddAccommodation", accommodation.Id.Hex())

	return nil
}

func (service *AccommodationService) getById(accommodationId primitive.ObjectID, span trace.Span, loki promtail.Client) (*domain.Accommodation, error) {
	util.HttpTraceInfo("Fetching accommodations by id...", span, loki, "getById", "")
	accommodation, err := service.store.Get(accommodationId)
	if err != nil {
		util.HttpTraceError(err, "Error getting accommodation by ID", span, service.loki, "getById", accommodationId.Hex())
		return nil, err
	}
	util.HttpTraceInfo("Accommodation retrieved by ID successfully", span, service.loki, "getById", accommodationId.Hex())

	return accommodation, nil
}

func (service *AccommodationService) GetAll(span trace.Span, loki promtail.Client) ([]*domain.Accommodation, error) {
	util.HttpTraceInfo("Fetching all accommodations...", span, loki, "GetAll", "")
	accommodations, err := service.store.GetAll()
	if err != nil {
		util.HttpTraceError(err, "Error getting all accommodations", span, service.loki, "GetAll", "")
		return nil, err
	}
	util.HttpTraceInfo("All accommodations retrieved successfully", span, service.loki, "GetAll", "")

	return accommodations, nil
}

func (service *AccommodationService) EditAccommodation(accommodation domain.Accommodation, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Updating accommodation...", span, loki, "EditAccommodation", "")
	err := service.store.Update(accommodation.Id, &accommodation)
	if err != nil {
		util.HttpTraceError(err, "Error editing accommodation", span, service.loki, "EditAccommodation", accommodation.Id.Hex())
		return err
	}
	util.HttpTraceInfo("Accommodation edited successfully", span, service.loki, "EditAccommodation", accommodation.Id.Hex())

	return nil
}

func (service *AccommodationService) DeleteAccommodation(id primitive.ObjectID, span trace.Span, loki promtail.Client) error {
	util.HttpTraceInfo("Deleting accommodation...", span, loki, "DeleteAccommodation", "")
	err := service.store.Delete(id)
	if err != nil {
		util.HttpTraceError(err, "Error deleting accommodation", span, service.loki, "DeleteAccommodation", id.Hex())
		return err
	}
	util.HttpTraceInfo("Accommodation deleted successfully", span, service.loki, "DeleteAccommodation", id.Hex())

	return nil
}

func (service *AccommodationService) OnCreateRatingChangeNotification(ratingChangedRequest dto.RatingChangedRequest, span trace.Span, loki promtail.Client) {
	accommodationId, err := primitive.ObjectIDFromHex(ratingChangedRequest.AccommodationId)
	if err != nil {
		util.HttpTraceError(err, "Error converting accommodation ID from hex", span, service.loki, "OnCreateRatingChangeNotification", ratingChangedRequest.AccommodationId)
		return
	}
	util.HttpTraceInfo("Updating rating...", span, loki, "OnCreateRatingChangeNotification", "")
	err = service.store.UpdateRating(accommodationId, ratingChangedRequest.Rating)
	if err != nil {
		util.HttpTraceError(err, "Error updating rating", span, service.loki, "OnCreateRatingChangeNotification", ratingChangedRequest.AccommodationId)
		return
	}
	util.HttpTraceInfo("Rating change notification processed successfully", span, service.loki, "OnCreateRatingChangeNotification", ratingChangedRequest.AccommodationId)
}
