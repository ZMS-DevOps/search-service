package dto

type RatingChangedRequest struct {
	AccommodationId string  `json:"id"`
	Rating          float32 `json:"rating"`
}
