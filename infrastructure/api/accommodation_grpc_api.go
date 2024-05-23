package api

import (
	"context"
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
	accommodationId, err := primitive.ObjectIDFromHex(accommodation.AccommodationId)
	if err != nil {
		return nil, err
	}

	if err := handler.service.AddAccommodation(*dto.MapAccommodation(accommodationId, accommodation)); err != nil {
		return nil, err
	}
	return &pb.AddAccommodationResponse{}, nil
}

func (handler *AccommodationGrpcHandler) GetHealth(ctx context.Context, request *pb.HealtRequest) (*pb.HealtResponse, error) {
	return &pb.HealtResponse{}, nil
}
