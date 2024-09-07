package amqpx

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shanluzhineng/configurationx/options/rabbitmq"
)

var (
	_abqpClient *AMQPClient
)

type AMQPClient struct {
	options *rabbitmq.DialOptions

	isConnected bool

	conn    *amqp.Connection
	channel *amqp.Channel
}

type WithChannel struct {
	Channel *amqp.Channel
}

// new a AMQPClient
func NewAMQPClient(options *rabbitmq.DialOptions) (*AMQPClient, error) {
	if options.RawUrl == "" {
		return nil, fmt.Errorf("options.RawUrl value is empty")
	}
	client := &AMQPClient{
		options: options,
	}
	return client, nil
}

// set Global client
func GlobalABQPClient(abqpClient *AMQPClient) {
	_abqpClient = abqpClient
}

// Get global *AMQPClient Instance, you can use GlobalABQPClient methods to set this value
func Client() *AMQPClient {
	return _abqpClient
}

// connect to amqp server and create channel
func (c *AMQPClient) Connect() error {
	if c.isConnected {
		return nil
	}

	conn, err := amqp.Dial(c.options.RawUrl)
	if err != nil {
		return err
	}
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	c.conn = conn
	c.channel = ch

	c.isConnected = true

	return nil
}

// get a new channel from amqp
func (c *AMQPClient) GetNewChannel() (*amqp.Channel, error) {
	err := c.ensureConnect()
	if err != nil {
		return nil, err
	}
	return c.conn.Channel()
}

// Close client
func (c *AMQPClient) Close() error {
	if c.conn == nil || c.conn.IsClosed() {
		return nil
	}
	if c.channel != nil {
		err := c.channel.Close()
		if err != nil {
			return err
		}
	}
	err := c.conn.Close()
	if err != nil {
		return err
	}
	c.conn = nil
	c.channel = nil
	c.isConnected = false
	return nil
}

func (c *AMQPClient) IsConnected() bool {
	return c.isConnected
}

func (c *AMQPClient) ensureConnect() error {
	return c.Connect()
}

// declare exchange
func (c *AMQPClient) ExchangeDeclare(declare ExchangeDeclare) error {
	err := c.ensureConnect()
	if err != nil {
		return err
	}

	return c.channel.ExchangeDeclare(declare.Name,
		string(declare.Kind),
		declare.Durable,
		declare.AutoDelete,
		declare.Internal,
		declare.NoWait,
		declare.Args)
}

// declare queue
func (c *AMQPClient) QueueDeclare(declare QueueDeclare) (*amqp.Queue, error) {
	err := c.ensureConnect()
	if err != nil {
		return nil, err
	}
	q, err := c.channel.QueueDeclare(declare.Name,
		declare.Durable,
		declare.AutoDelete,
		declare.Exclusive,
		declare.NoWait,
		declare.Args)
	if err != nil {
		return nil, err
	}
	return &q, nil
}

// bind exchange to a queue
func (c *AMQPClient) QueueBind(bind QueueBind) error {
	err := c.ensureConnect()
	if err != nil {
		return err
	}
	err = c.channel.QueueBind(bind.Queue, bind.RoutingKey, bind.Exchange, bind.NoWait, bind.Arguments)
	if err != nil {
		return err
	}
	return nil
}

func (c *AMQPClient) Qos(prefetchCount, prefetchSize int, global bool, channel ...WithChannel) error {
	err := c.ensureConnect()
	if err != nil {
		return err
	}
	usedChannel := c.channel
	if len(channel) > 0 {
		usedChannel = channel[0].Channel
	}
	return usedChannel.Qos(prefetchCount, prefetchSize, global)
}

// publish data to exchange
func (c *AMQPClient) PublishWithContext(ctx context.Context,
	exchange string,
	key string,
	mandatory bool,
	immediate bool,
	data []byte,
	channel ...WithChannel) error {
	if len(data) == 0 {
		return fmt.Errorf("argument data is empty")
	}
	if key == "" {
		return fmt.Errorf("key is empty")
	}
	err := c.ensureConnect()
	if err != nil {
		return err
	}
	usedChannel := c.channel
	if len(channel) > 0 {
		usedChannel = channel[0].Channel
	}
	err = usedChannel.PublishWithContext(ctx,
		exchange,  // exchange
		key,       // routing key
		mandatory, // mandatory
		immediate, // immediate
		amqp.Publishing{
			ContentType: "text/plan",
			Body:        data,
		})
	return err
}

// consume queue
//
// channel parameter indicate used specified channel,if nil or empty,then used default channel
func (c *AMQPClient) Consume(consume *QueueConsume, channel ...WithChannel) (<-chan amqp.Delivery, error) {
	if consume.Queue == "" {
		return nil, fmt.Errorf("consum.Queue field value cannot be empty")
	}
	err := c.ensureConnect()
	if err != nil {
		return nil, err
	}
	usedChannel := c.channel
	if len(channel) > 0 {
		usedChannel = channel[0].Channel
	}
	return usedChannel.Consume(consume.Queue,
		consume.Consumer,
		consume.AutoAck,
		consume.Exclusive,
		false,
		consume.NoWait,
		consume.Args)
}
