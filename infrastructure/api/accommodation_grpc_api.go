package api

import (
	"context"
	"fmt"
	"github.com/ZMS-DevOps/search-service/application"
	"github.com/ZMS-DevOps/search-service/domain"
	"github.com/ZMS-DevOps/search-service/infrastructure/dto"
	pb "github.com/ZMS-DevOps/search-service/proto"
	"github.com/ZMS-DevOps/search-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"go.mongodb.org/mongo-driver/bson/primitive"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type AccommodationGrpcHandler struct {
	pb.UnimplementedSearchServiceServer
	service       *application.AccommodationService
	traceProvider *sdktrace.TracerProvider
	loki          promtail.Client
}

func NewAccommodationGrpcHandler(service *application.AccommodationService, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *AccommodationGrpcHandler {
	return &AccommodationGrpcHandler{
		service:       service,
		traceProvider: traceProvider,
		loki:          loki,
	}
}

func (handler *AccommodationGrpcHandler) AddAccommodation(ctx context.Context, request *pb.AddAccommodationRequest) (*pb.AddAccommodationResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "add-accommodation-grpc")
	defer func() { span.End() }()

	accommodation := request.Accommodation

	if accommodation == nil {
		util.HttpTraceError(fmt.Errorf("accommodation is nil"), "accommodation is nil", span, handler.loki, "AddAccommodation", "")
		return nil, fmt.Errorf("accommodation is nil")
	}

	accommodationId, err := primitive.ObjectIDFromHex(accommodation.AccommodationId)
	if err != nil {
		util.HttpTraceError(err, "invalid accommodation id", span, handler.loki, "AddAccommodation", "")
		return nil, fmt.Errorf("invalid accommodation ID: %v", err)
	}
	mappedAccommodation := dto.MapAccommodation(accommodationId, accommodation)

	if err := handler.service.AddAccommodation(*mappedAccommodation, span, handler.loki); err != nil {
		util.HttpTraceError(err, "error adding accommodation", span, handler.loki, "AddAccommodation", "")
		return nil, fmt.Errorf("error adding accommodation: %v", err)
	}

	return &pb.AddAccommodationResponse{}, nil
}

func (handler *AccommodationGrpcHandler) EditAccommodation(ctx context.Context, request *pb.EditAccommodationRequest) (*pb.EditAccommodationResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "edit-accommodation-grpc")
	defer func() { span.End() }()
	accommodation := request.Accommodation

	if accommodation == nil {
		util.HttpTraceError(fmt.Errorf("accommodation is nil"), "accommodation is nil", span, handler.loki, "EditAccommodation", "")
		return nil, fmt.Errorf("accommodation is nil")
	}

	accommodationId, err := primitive.ObjectIDFromHex(accommodation.AccommodationId)
	if err != nil {
		util.HttpTraceError(fmt.Errorf("invalid accommodation id"), "invalid accommodation id", span, handler.loki, "EditAccommodation", accommodation.AccommodationId)
		return nil, fmt.Errorf("invalid accommodation ID: %v", err)
	}
	mappedAccommodation := dto.MapAccommodation(accommodationId, accommodation)
	if err := handler.service.EditAccommodation(*mappedAccommodation, span, handler.loki); err != nil {
		return nil, fmt.Errorf("error editting accommodation: %v", err)
	}

	return &pb.EditAccommodationResponse{}, nil
}

func (handler *AccommodationGrpcHandler) DeleteAccommodation(ctx context.Context, request *pb.DeleteAccommodationRequest) (*pb.DeleteAccommodationResponse, error) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(ctx, "delete-accommodation-grpc")
	defer func() { span.End() }()
	accommodationId, err := primitive.ObjectIDFromHex(request.AccommodationId)
	if err != nil {
		util.HttpTraceError(fmt.Errorf("invalid accommodation id"), "invalid accommodation id", span, handler.loki, "DeleteAccommodation", request.AccommodationId)
		return nil, fmt.Errorf("invalid accommodation id: %v", err)
	}

	if err := handler.service.DeleteAccommodation(accommodationId, span, handler.loki); err != nil {
		return nil, fmt.Errorf("error deleting accommodation: %v", err)
	}

	return &pb.DeleteAccommodationResponse{}, nil
}
