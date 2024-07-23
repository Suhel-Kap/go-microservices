package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	ampq "github.com/rabbitmq/amqp091-go"
	"github.com/suhel-kap/listener-service/event"
)

func main() {
	// try to connect to RabbitMQ
	rabbitConn, err := connect()
	if err != nil {
		log.Println("Failed to connect to RabbitMQ")
		os.Exit(1)
	}
	defer rabbitConn.Close()

	// start listening for messages
	log.Println("Listening for messages")

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	// watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*ampq.Connection, error) {
	var counts int64
	var backOff = 1 * time.Second
	var connection *ampq.Connection

	// dont continue until rabit is ready
	for {
		c, err := ampq.Dial("amqp://guest:guest@rabbitmq:5672")
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ")
			time.Sleep(backOff)
			counts++
		} else {
			connection = c
			log.Println("Connected to RabbitMQ")
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
