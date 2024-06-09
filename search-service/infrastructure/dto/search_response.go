package dto

type SearchResponse struct {
	Accommodations []AccommodationResponse `json:"accommodations"`
}

type AccommodationResponse struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Location       string  `json:"location"`
	NumberOfGuests int     `json:"number_of_guests"`
	MainPhoto      string  `json:"main_photo"`
	Rating         float32 `json:"rating"`
	TotalPrice     float32 `json:"total_price"`
}
