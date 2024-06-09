package api

import (
	"context"
	"fmt"
	"github.com/ZMS-DevOps/search-service/application"
	"github.com/ZMS-DevOps/search-service/infrastructure/dto"
	pb "github.com/ZMS-DevOps/search-service/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccommodationGrpcHandler struct {
	pb.UnimplementedSearchServiceServer
	service *application.AccommodationService
}

func NewAccommodationGrpcHandler(service *application.AccommodationService) *AccommodationGrpcHandler {
	return &AccommodationGrpcHandler{
		service: service,
	}
}

func (handler *AccommodationGrpcHandler) AddAccommodation(ctx context.Context, request *pb.AddAccommodationRequest) (*pb.AddAccommodationResponse, error) {
	accommodation := request.Accommodation

	if accommodation == nil {
		return nil, fmt.Errorf("accommodation is nil")
	}

	accommodationId, err := primitive.ObjectIDFromHex(accommodation.AccommodationId)
	if err != nil {
		return nil, fmt.Errorf("invalid accommodation ID: %v", err)
	}
	mappedAccommodation := dto.MapAccommodation(accommodationId, accommodation)

	if err := handler.service.AddAccommodation(*mappedAccommodation); err != nil {
		return nil, fmt.Errorf("error adding accommodation: %v", err)
	}

	return &pb.AddAccommodationResponse{}, nil
}

func (handler *AccommodationGrpcHandler) EditAccommodation(ctx context.Context, request *pb.EditAccommodationRequest) (*pb.EditAccommodationResponse, error) {
	accommodation := request.Accommodation

	if accommodation == nil {
		return nil, fmt.Errorf("accommodation is nil")
	}

	accommodationId, err := primitive.ObjectIDFromHex(accommodation.AccommodationId)
	if err != nil {
		return nil, fmt.Errorf("invalid accommodation ID: %v", err)
	}
	mappedAccommodation := dto.MapAccommodation(accommodationId, accommodation)
	if err := handler.service.EditAccommodation(*mappedAccommodation); err != nil {
		return nil, fmt.Errorf("error editting accommodation: %v", err)
	}

	return &pb.EditAccommodationResponse{}, nil
}

func (handler *AccommodationGrpcHandler) DeleteAccommodation(ctx context.Context, request *pb.DeleteAccommodationRequest) (*pb.DeleteAccommodationResponse, error) {
	accommodationId, err := primitive.ObjectIDFromHex(request.AccommodationId)
	if err != nil {
		return nil, fmt.Errorf("invalid accommodation ID: %v", err)
	}

	if err := handler.service.DeleteAccommodation(accommodationId); err != nil {
		return nil, fmt.Errorf("error deleting accommodation: %v", err)
	}

	return &pb.DeleteAccommodationResponse{}, nil
}
