package external

import (
	"context"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/search-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func NewBookingClient(address string) booking.BookingServiceClient {
	conn, err := getConnection(address)
	if err != nil {
		log.Fatalf("Failed to start gRPC connection to Catalogue service: %v", err)
	}
	return booking.NewBookingServiceClient(conn)
}

func getConnection(address string) (*grpc.ClientConn, error) {
	return grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
}

func FilterAvailableAccommodation(bookingClient booking.BookingServiceClient, ids []primitive.ObjectID, startDate time.Time, endDate time.Time, span trace.Span, loki promtail.Client) (*booking.FilterAvailableAccommodationResponse, error) {
	util.HttpTraceInfo("Filtering available accommodations...", span, loki, "FilterAvailableAccommodation", "")
	accommodationIDs := make([]string, len(ids))
	for i, id := range ids {
		accommodationIDs[i] = id.Hex()
	}

	response, err := bookingClient.FilterAvailableAccommodation(context.TODO(), &booking.FilterAvailableAccommodationRequest{
		AccommodationIds: accommodationIDs,
		StartDate:        startDate.Format(time.RFC3339),
		EndDate:          endDate.Format(time.RFC3339),
	})

	if err != nil {
		util.HttpTraceError(err, "Error filtering available accommodations", span, loki, "FilterAvailableAccommodation", "")
		return nil, err
	}
	util.HttpTraceInfo("Available accommodations filtered successfully", span, loki, "FilterAvailableAccommodation", "")
	return response, nil
}
