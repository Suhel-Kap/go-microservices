package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/suhel-kap/toolbox"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

var tools = toolbox.Tools{}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(q.Name, s, "logs_topic", false, nil)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var paylaod Payload
			_ = json.Unmarshal(d.Body, &paylaod)

			go handlePayload(paylaod)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]", q.Name)

	<-forever

	return nil
}

func handlePayload(p Payload) {
	fmt.Println("Received message", p)

	switch p.Name {
	case "log", "event":
		err := logEvent(p)
		if err != nil {
			log.Println("Failed to log event", err)
		}
	case "auth":
	// handle auth

	default:
		log.Println("Unknown event", p.Name)
	}
}

func logEvent(p Payload) error {
	jsonData, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		return err
	}

	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return err
	}

	return nil
}
