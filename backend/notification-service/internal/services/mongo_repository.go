package services

import (
	"context"
	"time"

	"notification-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	preferences *mongo.Collection
	records     *mongo.Collection
}

func NewMongoRepository(db *mongo.Database) (*MongoRepository, error) {
	repo := &MongoRepository{preferences: db.Collection("notification_preferences"), records: db.Collection("notification_records")}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := repo.records.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "eventId", Value: 1}}, Options: options.Index().SetUnique(true)},
		{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}},
		{Keys: bson.D{{Key: "expireAt", Value: 1}}, Options: options.Index().SetExpireAfterSeconds(0)},
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}

func (r *MongoRepository) GetPreference(ctx context.Context, userID string) (*models.Preference, error) {
	var preference models.Preference
	err := r.preferences.FindOne(ctx, bson.M{"_id": userID}).Decode(&preference)
	if err == mongo.ErrNoDocuments {
		return &models.Preference{UserID: userID, EmailEnabled: true}, nil
	}
	return &preference, err
}

func (r *MongoRepository) SavePreference(ctx context.Context, preference models.Preference) error {
	_, err := r.preferences.UpdateOne(ctx, bson.M{"_id": preference.UserID}, bson.M{"$set": bson.M{
		"emailEnabled": preference.EmailEnabled,
		"updatedAt":    preference.UpdatedAt,
	}}, options.Update().SetUpsert(true))
	return err
}

func (r *MongoRepository) CreateRecord(ctx context.Context, record models.Record) (bool, error) {
	_, err := r.records.InsertOne(ctx, record)
	if mongo.IsDuplicateKeyError(err) {
		return false, nil
	}
	return err == nil, err
}

func (r *MongoRepository) UpdateRecord(ctx context.Context, eventID, status, summary string) error {
	_, err := r.records.UpdateOne(ctx, bson.M{"eventId": eventID}, bson.M{"$set": bson.M{
		"status": status, "errorSummary": summary, "updatedAt": time.Now().UTC(),
	}})
	return err
}

func (r *MongoRepository) GetRecord(ctx context.Context, eventID string) (*models.Record, error) {
	var record models.Record
	err := r.records.FindOne(ctx, bson.M{"eventId": eventID}).Decode(&record)
	return &record, err
}

func (r *MongoRepository) ListRecords(ctx context.Context, userID string, page, size int) ([]models.Record, int64, error) {
	filter := bson.M{"userId": userID}
	total, err := r.records.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	cursor, err := r.records.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}).SetSkip(int64((page-1)*size)).SetLimit(int64(size)))
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)
	items := make([]models.Record, 0)
	if err := cursor.All(ctx, &items); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
