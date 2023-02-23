package amqpx

type ExchangeKind string

const (
	Exchange_Direct ExchangeKind = "direct"
	Exchange_Fanout ExchangeKind = "fanout"
	Exchange_Topic  ExchangeKind = "topic"
)

// ExchangeDeclare declares an exchange on the server. If the exchange does not
// already exist, the server will create it.  If the exchange exists, the server
// verifies that it is of the provided type, durability and auto-delete flags.
//
// Durable and Non-Auto-Deleted exchanges will survive server restarts and remain
// declared when there are no remaining bindings.  This is the best lifetime for
// long-lived exchange configurations like stable routes and default exchanges.

// Non-Durable and Auto-Deleted exchanges will be deleted when there are no
// remaining bindings and not restored on server restart.  This lifetime is
// useful for temporary topologies that should not pollute the virtual host on
// failure or after the consumers have completed.

// Non-Durable and Non-Auto-deleted exchanges will remain as long as the server is
// running including when there are no remaining bindings.  This is useful for
// temporary topologies that may have long delays between bindings.

// Durable and Auto-Deleted exchanges will survive server restarts and will be
// removed before and after server restarts when there are no remaining bindings.
// These exchanges are useful for robust temporary topologies or when you require
// binding durable queues to auto-deleted exchanges.
//
// Note: RabbitMQ declares the default exchange types like 'amq.fanout' as
// durable, so queues that bind to these pre-declared exchanges must also be
// durable.
type ExchangeDeclare struct {
	//Exchange names starting with "amq." are reserved for pre-declared and
	// standardized exchanges. The client MAY declare an exchange starting with
	// "amq." if the passive option is set, or the exchange already exists.  Names can
	// consist of a non-empty sequence of letters, digits, hyphen, underscore,
	// period, or colon.
	Name string

	// Each exchange belongs to one of a set of exchange kinds/types implemented by
	// the server. The exchange types define the functionality of the exchange - i.e.
	// how messages are routed through it. Once an exchange is declared, its type
	// cannot be changed.  The common types are "direct", "fanout", "topic" and
	// "headers".
	Kind ExchangeKind

	Durable    bool
	AutoDelete bool

	// Exchanges declared as `internal` do not accept accept publishings. Internal
	// exchanges are useful when you wish to implement inter-exchange topologies
	// that should not be exposed to users of the broker.
	Internal bool

	// When noWait is true, declare without waiting for a confirmation from the server.
	// The channel may be closed as a result of an error.  Add a NotifyClose listener
	// to respond to any exceptions.
	NoWait bool

	// Optional amqp.Table of arguments that are specific to the server's implementation of
	// the exchange can be sent for exchange types that require extra parameters.
	Args map[string]interface{}
}
