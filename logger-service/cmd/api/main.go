package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	WEB_PORT  = "80"
	RPC_PORT  = "5001"
	MONGO_URL = "mongodb://mongo:27017/?authSource=admin"
	GRPC_PORT = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create a context in order to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	// start the web server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", WEB_PORT),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connectToMongo() (*mongo.Client, error) {
	// create the connection options
	clientOptions := options.Client().ApplyURI(MONGO_URL)
	clientOptions.SetAuth(options.Credential{
		Username: "root",
		Password: "password",
	})

	// connect to the mongo server
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting to mongo")
		return nil, err
	}

	log.Println("Connected to mongo")

	return c, nil
}
