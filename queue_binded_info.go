package amqpx

// queueBindKey
type QueueBind struct {
	// Queue name
	Queue string
	// Exchange name
	Exchange string
	// Routing key
	RoutingKey string
	NoWait     bool
	Arguments  map[string]interface{}
}

func NewQueueBind(queue string,
	routingKey string,
	exchange string,
	noWait bool) *QueueBind {

	return &QueueBind{
		Queue:      queue,
		Exchange:   exchange,
		RoutingKey: routingKey,
		NoWait:     noWait,
		Arguments:  make(map[string]interface{}),
	}
}
