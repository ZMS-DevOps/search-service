package persistence

import (
	"context"
	"github.com/ZMS-DevOps/search-service/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const (
	DATABASE   = "searchaccommodationdb"
	COLLECTION = "searchaccommodation"
)

type AccommodationMongoDBStore struct {
	accommodations *mongo.Collection
}

func NewAccommodationMongoDBStore(client *mongo.Client) domain.AccommodationStore {
	accommodations := client.Database(DATABASE).Collection(COLLECTION)
	return &AccommodationMongoDBStore{
		accommodations: accommodations,
	}
}

func (store *AccommodationMongoDBStore) Search(location string, guestNumber int, startDate time.Time, endDate time.Time, minPrice float32, maxPrice float32) ([]*domain.SearchResponse, error) {
	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"location":         location,
			"guest_number.min": bson.M{"$lte": guestNumber},
			"guest_number.max": bson.M{"$gte": guestNumber},
		}},
		bson.M{"$addFields": bson.M{
			"total_price": bson.M{
				"$sum": bson.A{
					calculateDailyPrices(startDate, endDate, guestNumber),
				},
			},
			"number_of_quests": bson.M{
				"$toInt": guestNumber,
			},
		}},
		bson.M{"$match": bson.M{
			"total_price": bson.M{"$gte": minPrice, "$lte": maxPrice},
		}},
	}

	cursor, err := store.accommodations.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var accommodations []*domain.SearchResponse
	if err = cursor.All(context.TODO(), &accommodations); err != nil {
		return nil, err
	}

	return accommodations, nil
}

func calculateDailyPrices(startDate time.Time, endDate time.Time, guestNumber int) bson.A {
	days := int(endDate.Sub(startDate).Hours() / 24)
	dailyPrices := bson.A{}

	for i := 0; i < days; i++ {
		date := startDate.Add(time.Hour * 24 * time.Duration(i))
		dailyPrices = append(dailyPrices, bson.M{
			"$cond": bson.M{
				"if": bson.M{
					"$and": bson.A{
						bson.M{"$gte": bson.A{date, "$special_price.date_range.start"}},
						bson.M{"$lte": bson.A{date, "$special_price.date_range.end"}},
					},
				},
				"then": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.A{"$default_price.type", domain.PerGuest}},
						"then": bson.M{"$multiply": bson.A{"$special_price.price", guestNumber}},
						"else": "$special_price.price",
					},
				},
				"else": bson.M{
					"$cond": bson.M{
						"if":   bson.M{"$eq": bson.A{"$default_price.type", domain.PerGuest}},
						"then": bson.M{"$multiply": bson.A{"$default_price.price", guestNumber}},
						"else": "$default_price.price",
					},
				},
			},
		})
	}

	return dailyPrices
}

func (store *AccommodationMongoDBStore) Get(id primitive.ObjectID) (*domain.Accommodation, error) {
	filter := bson.M{"_id": id}
	return store.filterOne(filter)
}

func (store *AccommodationMongoDBStore) GetAll() ([]*domain.Accommodation, error) {
	filter := bson.D{{}}
	return store.filter(filter)
}

func (store *AccommodationMongoDBStore) Insert(accommodation *domain.Accommodation) error {
	accommodation.Id = primitive.NewObjectID()
	result, err := store.accommodations.InsertOne(context.TODO(), accommodation)
	if err != nil {
		return err
	}
	accommodation.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *AccommodationMongoDBStore) InsertWithId(accommodation *domain.Accommodation) error {
	result, err := store.accommodations.InsertOne(context.TODO(), accommodation)
	if err != nil {
		return err
	}
	accommodation.Id = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (store *AccommodationMongoDBStore) DeleteAll() {
	store.accommodations.DeleteMany(context.TODO(), bson.D{{}})
}

func (store *AccommodationMongoDBStore) Delete(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := store.accommodations.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}

func (store *AccommodationMongoDBStore) Update(id primitive.ObjectID, accommodation *domain.Accommodation) error {
	filter := bson.M{"_id": id}

	updateFields := bson.D{
		{"name", accommodation.Name},
		{"location", accommodation.Location},
		{"main_photo", accommodation.MainPhoto},
		{"guest_number", accommodation.GuestNumber},
		{"default_price", accommodation.DefaultPrice},
	}
	update := bson.D{{"$set", updateFields}}

	_, err := store.accommodations.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (store *AccommodationMongoDBStore) filter(filter interface{}) ([]*domain.Accommodation, error) {
	cursor, err := store.accommodations.Find(context.TODO(), filter)
	defer cursor.Close(context.TODO())

	if err != nil {
		return nil, err
	}
	return decode(cursor)
}

func (store *AccommodationMongoDBStore) filterOne(filter interface{}) (accommodation *domain.Accommodation, err error) {
	result := store.accommodations.FindOne(context.TODO(), filter)
	err = result.Decode(&accommodation)
	return
}

func decode(cursor *mongo.Cursor) (accommodations []*domain.Accommodation, err error) {
	for cursor.Next(context.TODO()) {
		var accommodation domain.Accommodation
		err = cursor.Decode(&accommodation)
		if err != nil {
			return
		}
		accommodations = append(accommodations, &accommodation)
	}
	err = cursor.Err()
	return
}

func (store *AccommodationMongoDBStore) UpdateDefaultPrice(id primitive.ObjectID, price *float64) error {
	if price == nil {
		return nil
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"default_price.price": *price}}

	_, err := store.accommodations.UpdateOne(context.TODO(), filter, update)
	return err
}

func (store *AccommodationMongoDBStore) UpdateSpecialPrice(id primitive.ObjectID, updatedSpecialPrices []domain.SpecialPrice) error {
	filter := bson.M{"_id": id}
	update := bson.M{}

	_, err := store.GetSpecialPrices(id)
	if err != nil {
		return err
	}

	update = bson.M{"$set": bson.M{"special_price": updatedSpecialPrices}}

	_, err = store.accommodations.UpdateOne(context.TODO(), filter, update)
	return err
}

func (store *AccommodationMongoDBStore) GetSpecialPrices(id primitive.ObjectID) ([]domain.SpecialPrice, error) {
	var accommodation domain.Accommodation
	filter := bson.M{"_id": id}
	err := store.accommodations.FindOne(context.TODO(), filter).Decode(&accommodation)
	if err != nil {
		return nil, err
	}
	return accommodation.SpecialPrice, nil
}

//func (store *AccommodationMongoDBStore) UpdateTypeOfPayment(id primitive.ObjectID, typeOfPayment *string) error {
//	if typeOfPayment == nil {
//		return fmt.Errorf("payment type is nil but should not be")
//	}
//
//	var pricingType = dto.MapPricingType(typeOfPayment)
//	if pricingType == nil {
//		return fmt.Errorf("payment type is nil but should not be")
//	}
//	filter := bson.M{"_id": id}
//	update := bson.M{"$set": bson.M{"default_price.type": pricingType}}
//
//	_, err := store.accommodations.UpdateOne(context.TODO(), filter, update)
//	return err
//}
