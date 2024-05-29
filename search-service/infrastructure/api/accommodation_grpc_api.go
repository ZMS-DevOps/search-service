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
	fmt.Println("Received accommodation:", accommodation)

	if accommodation == nil {
		return nil, fmt.Errorf("accommodation is nil")
	}

	fmt.Println("Accommodation ID:", accommodation.AccommodationId)
	accommodationId, err := primitive.ObjectIDFromHex(accommodation.AccommodationId)
	if err != nil {
		return nil, fmt.Errorf("invalid accommodation ID: %v", err)
	}
	fmt.Println("Stop 1")
	mappedAccommodation := dto.MapAccommodation(accommodationId, accommodation)

	if err := handler.service.AddAccommodation(*mappedAccommodation); err != nil {
		return nil, fmt.Errorf("error adding accommodation: %v", err)
	}
	fmt.Println("Stop 3")

	return &pb.AddAccommodationResponse{}, nil
}

func (handler *AccommodationGrpcHandler) EditAccommodation(ctx context.Context, request *pb.EditAccommodationRequest) (*pb.EditAccommodationResponse, error) {
	accommodation := request.Accommodation
	fmt.Println("Received accommodation:", accommodation)

	if accommodation == nil {
		return nil, fmt.Errorf("accommodation is nil")
	}

	fmt.Println("Accommodation ID:", accommodation.AccommodationId)
	accommodationId, err := primitive.ObjectIDFromHex(accommodation.AccommodationId)
	if err != nil {
		return nil, fmt.Errorf("invalid accommodation ID: %v", err)
	}
	fmt.Println("Stop 1")
	mappedAccommodation := dto.MapAccommodation(accommodationId, accommodation)
	fmt.Println("Stop 2")
	fmt.Println("mappedAccommodation ", mappedAccommodation)
	if err := handler.service.EditAccommodation(*mappedAccommodation); err != nil {
		return nil, fmt.Errorf("error editting accommodation: %v", err)
	}
	fmt.Println("Stop 3")

	return &pb.EditAccommodationResponse{}, nil
}

func (handler *AccommodationGrpcHandler) DeleteAccommodation(ctx context.Context, request *pb.DeleteAccommodationRequest) (*pb.DeleteAccommodationResponse, error) {
	fmt.Println("Accommodation ID:", request.AccommodationId)
	accommodationId, err := primitive.ObjectIDFromHex(request.AccommodationId)
	if err != nil {
		return nil, fmt.Errorf("invalid accommodation ID: %v", err)
	}

	if err := handler.service.DeleteAccommodation(accommodationId); err != nil {
		return nil, fmt.Errorf("error deleting accommodation: %v", err)
	}

	return &pb.DeleteAccommodationResponse{}, nil
}
