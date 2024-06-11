package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Accommodation struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	HostId       string             `bson:"host_id"`
	Location     string             `bson:"location"`
	MainPhoto    string             `bson:"main_photo"`
	Rating       float32            `bson:"rating"`
	GuestNumber  GuestNumber        `bson:"guest_number"`
	DefaultPrice DefaultPrice       `bson:"default_price"`
	SpecialPrice []SpecialPrice     `bson:"special_price"`
}

type GuestNumber struct {
	Min int `bson:"min"`
	Max int `bson:"max"`
}

type PricingType int

const (
	PerApartmentUnit PricingType = iota
	PerGuest
)

func (p PricingType) String() string {
	switch p {
	case PerApartmentUnit:
		return "perApartment"
	case PerGuest:
		return "perPerson"
	default:
		return "Unknown"
	}
}

type DefaultPrice struct {
	Price float32     `bson:"price"`
	Type  PricingType `bson:"type"`
}

type SpecialPrice struct {
	Price     float32   `bson:"price"`
	DateRange DateRange `bson:"date_range"`
}

type DateRange struct {
	Start time.Time `bson:"start"`
	End   time.Time `bson:"end"`
}
