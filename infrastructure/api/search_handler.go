package api

import (
	"encoding/json"
	"fmt"
	"github.com/ZMS-DevOps/search-service/application"
	"github.com/ZMS-DevOps/search-service/domain"
	"github.com/gorilla/mux"
	"net/http"
)

type SearchHandler struct {
	service *application.SearchService
}

type SearchResponse struct {
	//Search *domain.Search `json:"search"`
	Accommodations []*domain.Accommodation `json:"accommodations"`
}

type HealthCheckResponse struct {
	Size string `json:"size"`
}

func NewSearchHandler(service *application.SearchService) *SearchHandler {
	server := &SearchHandler{
		service: service,
	}
	return server
}

func (handler *SearchHandler) Init(router *mux.Router) {
	router.HandleFunc(`/search/search`, handler.GetAll).Methods("GET")
	//router.HandleFunc("/hotel/search/{id}", handler.GetById).Methods("GET")
	//router.HandleFunc("/hotel/search", handler.Add).Methods("POST")
	//router.HandleFunc("/hotel/search/{id}", handler.Update).Methods("PUT")
	//router.HandleFunc("/hotel/search/{id}", handler.Delete).Methods("DELETE")
	//router.HandleFunc("/hotel/search/price/{id}", handler.UpdatePrice).Methods("PUT")
	router.HandleFunc("/search/health", handler.GetHealthCheck).Methods("GET")
}

func (handler *SearchHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	accommodations, err := handler.service.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := SearchResponse{
		Accommodations: accommodations,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *SearchHandler) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthCheckResponse{
		Size: "Search SERVICE OK",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func handleError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, message)
}
