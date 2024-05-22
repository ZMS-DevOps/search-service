package config

import "os"

type Config struct {
	Port            string
	GrpcPort        string
	HotelDBHost     string
	HotelDBPort     string
	HotelDBUsername string
	HotelDBPassword string
	BookingHost     string
	BookingPort     string
}

func NewConfig() *Config {
	return &Config{
		Port: os.Getenv("SERVICE_PORT"),
		//GrpcPort:        os.Getenv("GRPC_PORT"),
		HotelDBHost:     os.Getenv("DB_HOST"),
		HotelDBPort:     os.Getenv("DB_PORT"),
		HotelDBUsername: os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		HotelDBPassword: os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
		BookingHost:     os.Getenv("BOOKING_HOST"),
		BookingPort:     os.Getenv("BOOKING_PORT"),
	}
}
