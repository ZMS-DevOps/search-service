package startup

import (
	"fmt"
	booking "github.com/ZMS-DevOps/booking-service/proto"
	"github.com/ZMS-DevOps/search-service/application"
	"github.com/ZMS-DevOps/search-service/application/external"
	"github.com/ZMS-DevOps/search-service/domain"
	"github.com/ZMS-DevOps/search-service/infrastructure/api"
	"github.com/ZMS-DevOps/search-service/infrastructure/persistence"
	accommodationSearch "github.com/ZMS-DevOps/search-service/proto"
	"github.com/ZMS-DevOps/search-service/startup/config"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

type Server struct {
	config                   *config.Config
	router                   *mux.Router
	SearchHandler            *api.SearchHandler
	AccommodationGrpcHandler *api.AccommodationGrpcHandler
	AccommodationHandler     *api.AccommodationHandler
}

func NewServer(config *config.Config) *Server {
	server := &Server{
		config: config,
		router: mux.NewRouter(),
	}
	server.SearchHandler, server.AccommodationGrpcHandler, server.AccommodationHandler = server.setupHandlers()
	return server
}

func (server *Server) setupHandlers() (*api.SearchHandler, *api.AccommodationGrpcHandler, *api.AccommodationHandler) {
	mongoClient := server.initMongoClient()
	bookingClient := external.NewBookingClient(server.getBookingAddress())
	accommodationStore := server.initAccommodationStore(mongoClient)
	searchService := server.initSearchService(accommodationStore, bookingClient)
	searchHandler := server.initSearchHandler(searchService)
	searchHandler.Init(server.router)
	accommodationService := server.initAccommodationService(accommodationStore)
	accommodationGrpcHandler := server.initAccommodationGrpcHandler(accommodationService)
	accommodationHandler := server.initAccommodationHandler(accommodationService)
	return searchHandler, accommodationGrpcHandler, accommodationHandler
}

func (server *Server) Start() {
	go server.startGrpcServer(server.AccommodationGrpcHandler)
	fmt.Println(fmt.Sprintf(":%s", server.config.Port), server.router)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", server.config.Port), server.router))
}

func (server *Server) getBookingAddress() string {
	return fmt.Sprintf("%s:%s", server.config.BookingHost, server.config.BookingPort)
}

func (server *Server) initMongoClient() *mongo.Client {
	client, err := persistence.GetClient(server.config.HotelDBUsername, server.config.HotelDBPassword, server.config.HotelDBHost, server.config.HotelDBPort)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func (server *Server) initAccommodationStore(client *mongo.Client) domain.AccommodationStore {
	store := persistence.NewAccommodationMongoDBStore(client)
	store.DeleteAll()
	for _, accommodation := range accommodations {
		err := store.InsertWithId(accommodation)
		if err != nil {
			log.Fatal(err)
		}
	}
	return store
}

func (server *Server) startGrpcServer(bookingHandler *api.AccommodationGrpcHandler) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	accommodationSearch.RegisterSearchServiceServer(grpcServer, bookingHandler)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}

func (server *Server) initSearchService(store domain.AccommodationStore, bookingClient booking.BookingServiceClient) *application.SearchService {
	return application.NewSearchService(store, bookingClient)
}

func (server *Server) initAccommodationService(store domain.AccommodationStore) *application.AccommodationService {
	return application.NewAccommodationService(store)
}

func (server *Server) initSearchHandler(service *application.SearchService) *api.SearchHandler {
	return api.NewSearchHandler(service)
}

func (server *Server) initAccommodationGrpcHandler(service *application.AccommodationService) *api.AccommodationGrpcHandler {
	return api.NewAccommodationGrpcHandler(service)
}

func (server *Server) initAccommodationHandler(service *application.AccommodationService) *api.AccommodationHandler {
	return api.NewAccommodationHandler(service)
}
