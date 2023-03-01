package amqpx

// consume queue request parameter structure
type QueueConsume struct {
	Queue string

	// The consumer is identified by a string that is unique and scoped for all
	// consumers on this channel.  If you wish to eventually cancel the consumer, use
	// the same non-empty identifier in Channel.Cancel.  An empty string will cause
	// the library to generate a unique identity.  The consumer identity will be
	// included in every Delivery in the ConsumerTag field
	Consumer string

	// When autoAck (also known as noAck) is true, the server will acknowledge
	// deliveries to this consumer prior to writing the delivery to the network.  When
	// autoAck is true, the consumer should not call Delivery.Ack. Automatically
	// acknowledging deliveries means that some deliveries may get lost if the
	// consumer is unable to process them after the server delivers them.
	// See http://www.rabbitmq.com/confirms.html for more details.
	AutoAck bool

	// When exclusive is true, the server will ensure that this is the sole consumer
	// from this queue. When exclusive is false, the server will fairly distribute
	// deliveries across multiple consumers.
	Exclusive bool

	// When noWait is true, do not wait for the server to confirm the request and
	// immediately begin deliveries.  If it is not possible to consume, a channel
	// exception will be raised and the channel will be closed.
	NoWait bool

	Args map[string]interface{}
}

func NewDefaultQueueConsume(queue string) *QueueConsume {
	c := &QueueConsume{
		Queue:     queue,
		AutoAck:   true,
		Exclusive: false,
		NoWait:    false,
	}
	return c
}
