package amqpx

// QueueDeclare declares a queue to hold messages and deliver to consumers.
// Declaring creates a queue if it doesn't already exist, or ensures that an
// existing queue matches the same parameters.
//
// Every queue declared gets a default binding to the empty exchange "" which has
// the type "direct" with the routing key matching the queue's name.  With this
// default binding, it is possible to publish messages that route directly to
// this queue by publishing to "" with the routing key of the queue name.
//
// Durable and Non-Auto-Deleted queues will survive server restarts and remain
// when there are no remaining consumers or bindings.  Persistent publishings will
// be restored in this queue on server restart.  These queues are only able to be
// bound to durable exchanges.
//
// Non-Durable and Auto-Deleted queues will not be redeclared on server restart
// and will be deleted by the server after a short time when the last consumer is
// canceled or the last consumer's channel is closed.  Queues with this lifetime
// can also be deleted normally with QueueDelete.  These durable queues can only
// be bound to non-durable exchanges.
//
// Non-Durable and Non-Auto-Deleted queues will remain declared as long as the
// server is running regardless of how many consumers.  This lifetime is useful
// for temporary topologies that may have long delays between consumer activity.
// These queues can only be bound to non-durable exchanges.

// Durable and Auto-Deleted queues will be restored on server restart, but without
// active consumers will not survive and be removed.  This Lifetime is unlikely
// to be useful.
type QueueDeclare struct {
	// The queue name may be empty, in which case the server will generate a unique name
	Name string

	Durable    bool
	AutoDelete bool //delete when unused

	// Exclusive queues are only accessible by the connection that declares them and
	// will be deleted when the connection closes.  Channels on other connections
	// will receive an error when attempting  to declare, bind, consume, purge or
	// delete a queue with the same name.
	Exclusive bool

	// When noWait is true, the queue will assume to be declared on the server.  A
	// channel exception will arrive if the conditions are met for existing queues
	// or attempting to modify an existing queue from a different connection.
	NoWait bool

	Args map[string]interface{}
}
