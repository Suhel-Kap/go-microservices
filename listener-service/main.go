package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
)

func main() {
	// try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
		os.Exit(1)
	}
	defer rabbitConn.Close()

	log.Println("Connected to RabbitMQ")
	// start listening for messages

	// create consumer

	// watch the queue and consume events

}

func connect() (*ampq.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *ampq.Connection

	// dont continue until rabit is ready
	for {
		c, err := ampq.Dial("amqp://guest:guest@localhost:5672")
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ")
			time.Sleep(backOff)
			counts++
		} else {
			connection = c
			break
		}

		if counts > 5 {
			fmt.Println("Failed to connect to RabbitMQ after 5 attempts", err)
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		log.Printf("Retrying connection in %v", backOff)
		time.Sleep(backOff)
		continue
	}

	return connection, nil
}
