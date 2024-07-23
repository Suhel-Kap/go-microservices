module github.com/suhel-kap/listener-service

go 1.22.4

require (
	github.com/suhel-kap/toolbox v0.0.0
)

replace github.com/suhel-kap/toolbox => ../toolbox

require github.com/rabbitmq/amqp091-go v1.10.0 // indirect
