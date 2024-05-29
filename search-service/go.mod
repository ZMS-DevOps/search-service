module github.com/ZMS-DevOps/search-service
go 1.17

// replace github.com/tamararankovic/microservices_demo/common => ../common
replace github.com/ZMS-DevOps/booking-service => ./../../../booking-service/booking-service

require (
	github.com/ZMS-DevOps/booking-service v1.0.0
	github.com/go-playground/validator/v10 v10.20.0
	github.com/gorilla/mux v1.8.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.9.0
	go.mongodb.org/mongo-driver v1.8.4
	google.golang.org/grpc v1.45.0
	google.golang.org/protobuf v1.27.1
)

require (
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	google.golang.org/genproto v0.0.0-20220314164441-57ef72a4c106 // indirect
)

require (
	github.com/go-stack/stack v1.8.0 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	// github.com/tamararankovic/microservices_demo/common v1.0.0
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.0.2 // indirect
	github.com/xdg-go/stringprep v1.0.2 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e // indirect
	golang.org/x/text v0.14.0 // indirect
)
