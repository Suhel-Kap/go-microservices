package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string    `json:"name" bson:"name"`
	Data      string    `json:"data" bson:"data"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")

	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry")
		return err
	}

	return nil
}

func (l *LogEntry) FindAll() ([]*LogEntry, error) {
	// this is done in case this func takes > 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	opts := options.Find()
	opts.SetSort(bson.D{{"created_at", -1}})

	cursor, err := collection.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		log.Println("Error finding log entries")
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []*LogEntry

	for cursor.Next(ctx) {
		var item LogEntry
		err := cursor.Decode(&item)
		if err != nil {
			log.Println("Error decoding log entry")
			return nil, err
		} else {
			entries = append(entries, &item)
		}
	}

	return entries, nil
}

func (l *LogEntry) FindOne(id string) (*LogEntry, error) {
	// this is done in case this func takes > 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting id to object id")
		return nil, err
	}

	var entry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&entry)
	if err != nil {
		log.Println("Error finding log entry")
		return nil, err
	}

	return &entry, nil
}

func (l *LogEntry) DropCollection() error {
	// this is done in case this func takes > 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	if err := collection.Drop(ctx); err != nil {
		log.Println("Error dropping collection")
		return err
	}

	return nil
}

func (l *LogEntry) Update() (*mongo.UpdateResult, error) {
	// this is done in case this func takes > 15 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	collection := client.Database("logs").Collection("logs")

	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("Error converting id to object id")
		return nil, err
	}

	result, err := collection.UpdateOne(ctx, bson.M{"_id": docId}, bson.D{
		{"$set", bson.D{
			{"name", l.Name},
			{"data", l.Data},
			{"updated_at", time.Now()},
		}},
	})

	if err != nil {
		log.Println("Error updating log entry")
		return nil, err
	}

	return result, nil
}
