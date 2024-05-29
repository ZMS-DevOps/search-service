package external

import (
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
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

//func FilterAvailableAccommodation(bookingClient booking.BookingServiceClient, ids []primitive.ObjectID, startDate time.Time, endDate time.Time) (*booking.FilterAvailableAccommodationResponse, error) {
//	accommodationIDs := make([]string, len(ids))
//	for i, id := range ids {
//		accommodationIDs[i] = id.Hex()
//	}
//	return bookingClient.FilterAvailableAccommodation(context.TODO(), &booking.FilterAvailableAccommodationRequest{AccommodationIds: accommodationIDs, StartDate: startDate.Format(time.RFC3339), EndDate: endDate.Format(time.RFC3339)})
//}
