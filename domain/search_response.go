package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type SearchResponse struct {
	Id             primitive.ObjectID `json:"id" bson:"_id"`
	Name           string             `json:"name"`
	Location       string             `json:"location"`
	NumberOfGuests int                `json:"number_of_quests"`
	MainPhoto      string             `bson:"main_photo"`
	Rating         float32            `bson:"rating"`
	TotalPrice     float32            `bson:"total_price"`
}
