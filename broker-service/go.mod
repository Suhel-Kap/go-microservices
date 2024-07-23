module github.com/suhel-kap/broker

go 1.22.4

require (
	github.com/go-chi/chi/v5 v5.1.0
	github.com/go-chi/cors v1.2.1
	github.com/suhel-kap/toolbox v0.0.0
)

require github.com/rabbitmq/amqp091-go v1.10.0 // indirect

replace github.com/suhel-kap/toolbox => ../toolbox
