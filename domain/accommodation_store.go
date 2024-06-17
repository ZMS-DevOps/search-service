package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AccommodationStore interface {
	Get(id primitive.ObjectID) (*Accommodation, error)
	GetAll() ([]*Accommodation, error)
	Insert(accommodation *Accommodation) error
	InsertWithId(accommodation *Accommodation) error
	DeleteAll()
	Delete(id primitive.ObjectID) error
	Update(id primitive.ObjectID, accommodation *Accommodation) error
	UpdateDefaultPrice(id primitive.ObjectID, price *float64) error
	UpdateSpecialPrice(id primitive.ObjectID, newSpecialPrices []SpecialPrice) error
	GetSpecialPrices(id primitive.ObjectID) ([]SpecialPrice, error)
	Search(location string, guestNumber int, startDate time.Time, endDate time.Time, minPrice float32, maxPrice float32) ([]*Accommodation, error)
	UpdateRating(accommodationId primitive.ObjectID, rating float32) error
	GetByHostId(hostId string) ([]*Accommodation, error)
}
