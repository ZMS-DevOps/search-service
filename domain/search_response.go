package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchResponse struct {
	Id         primitive.ObjectID `json:"id" bson:"_id"`
	HostId     string             `json:"host_id"`
	Name       string             `json:"name"`
	Location   string             `json:"location"`
	MainPhoto  string             `json:"main_photo"`
	Rating     float32            `json:"rating"`
	TotalPrice float32            `json:"total_price"`
	UnitPrice  float32            `json:"unit_price"`
	PriceType  string             `json:"price_type"`
}
