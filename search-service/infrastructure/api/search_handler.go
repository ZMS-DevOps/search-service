package api

import (
	"encoding/json"
	"fmt"
	"github.com/ZMS-DevOps/search-service/application"
	"github.com/ZMS-DevOps/search-service/infrastructure/dto"
	"github.com/gorilla/mux"
	"net/http"
)

type SearchHandler struct {
	service *application.SearchService
}

//type GetAllResponse struct {
//	//Search *domain.Search `json:"search"`
//	Accommodations []*domain.Accommodation `json:"accommodations"`
//}
//
//type SearchResponse struct {
//	//Search *domain.Search `json:"search"`
//	Accommodations []*domain.SearchResponse `json:"accommodations"`
//}

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
	router.HandleFunc(`/search/search/all`, handler.GetAll).Methods("GET")
	router.HandleFunc("/search/all", handler.Search).Methods("POST")
	router.HandleFunc("/search/health", handler.GetHealthCheck).Methods("GET")
}

func (handler *SearchHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	accommodations, err := handler.service.GetAll()

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(accommodations)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	var searchDto dto.SearchDto
	if err := json.NewDecoder(r.Body).Decode(&searchDto); err != nil {
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dto.ValidateSearch(searchDto); err != nil {
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	accommodations, err := handler.service.Search(
		searchDto.Location, searchDto.GuestNumber, searchDto.Start, searchDto.End,
		searchDto.MinPrice, searchDto.MaxPrice)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(accommodations)
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
