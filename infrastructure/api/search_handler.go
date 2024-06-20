package api

import (
	"encoding/json"
	"fmt"
	"github.com/ZMS-DevOps/search-service/application"
	"github.com/ZMS-DevOps/search-service/domain"
	"github.com/ZMS-DevOps/search-service/infrastructure/dto"
	"github.com/ZMS-DevOps/search-service/util"
	"github.com/afiskon/promtail-client/promtail"
	"github.com/gorilla/mux"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"net/http"
)

type SearchHandler struct {
	service       *application.SearchService
	traceProvider *sdktrace.TracerProvider
	loki          promtail.Client
}

type HealthCheckResponse struct {
	Size string `json:"size"`
}

func NewSearchHandler(service *application.SearchService, traceProvider *sdktrace.TracerProvider, loki promtail.Client) *SearchHandler {
	server := &SearchHandler{
		service:       service,
		traceProvider: traceProvider,
		loki:          loki,
	}
	return server
}

func (handler *SearchHandler) Init(router *mux.Router) {
	router.HandleFunc("/search/all", handler.Search).Methods("POST")
	router.HandleFunc("/search/{id}", handler.GetByHostId).Methods("GET")
	router.HandleFunc("/search/health", handler.GetHealthCheck).Methods("GET")
}

func (handler *SearchHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-all-get")
	defer func() { span.End() }()
	accommodations, err := handler.service.GetAll(span, handler.loki)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResponse, err := json.Marshal(accommodations)
	if err != nil {
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "GetAll", "")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *SearchHandler) GetByHostId(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "get-by-host-id-get")
	defer func() { span.End() }()
	vars := mux.Vars(r)
	hostId := vars["id"]
	accommodations, err := handler.service.GetByHostId(hostId, span, handler.loki)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	responses := handler.service.MapToGetByHostIdResponse(accommodations, span)

	jsonResponse, err := json.Marshal(responses)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (handler *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	_, span := handler.traceProvider.Tracer(domain.ServiceName).Start(r.Context(), "search-post")
	defer func() { span.End() }()
	var searchDto dto.SearchDto
	if err := json.NewDecoder(r.Body).Decode(&searchDto); err != nil {
		util.HttpTraceError(err, "failed to parse request body", span, handler.loki, "Search", "")
		handleError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := dto.ValidateSearch(searchDto); err != nil {
		util.HttpTraceError(err, "failed to validate data", span, handler.loki, "Search", "")
		handleError(w, http.StatusBadRequest, err.Error())
		return
	}

	accommodations, err := handler.service.Search(
		searchDto.Location, searchDto.GuestNumber, searchDto.Start, searchDto.End,
		searchDto.MinPrice, searchDto.MaxPrice, span, handler.loki)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(accommodations)
	if err != nil {
		util.HttpTraceError(err, "failed to marshal data", span, handler.loki, "Search", "")
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
