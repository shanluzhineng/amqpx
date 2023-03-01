package amqpx

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"sync/atomic"

	amqp "github.com/rabbitmq/amqp091-go"
)

var consumerSeq uint64

const consumerTagLengthMax = 0xFF // see writeShortstr

type IAMQPService interface {
	IAMQPPublisher
	IAMQPConsumer

	ExchangeDeclare(declare ExchangeDeclare) error
	QueueDeclare(declare QueueDeclare) error
	QueueBind(bind QueueBind) error
}

type IAMQPPublisher interface {
	Publish(v interface{}, opts ...PublishOption) error
}

type IAMQPConsumer interface {
	SimpleConsume(topic string, consumer string, observeFn func(msg *DeliveryMessage)) (ITopicConsumer, error)

	// general unique consumer name, this cannot guarante value is unique in distributed environment
	GenerateUniqueConsumerName() string
}

type amqpService struct {
	client *AMQPClient
	defaultTopicConsumer

	publishChannel *amqp.Channel
	// consuming channel, key is topic
	consumedChannel map[string]*amqp.Channel
	lock            sync.Mutex
}

// create IAMQPService instance
func NewAMQPService(client *AMQPClient) IAMQPService {
	return &amqpService{
		client:          client,
		consumedChannel: make(map[string]*amqp.Channel),
	}
}

// #region IAMQPService members

func (s *amqpService) ExchangeDeclare(declare ExchangeDeclare) error {
	err := s.client.ExchangeDeclare(declare)
	if err != nil {
		return err
	}
	return nil
}

func (s *amqpService) QueueDeclare(declare QueueDeclare) error {
	_, err := s.client.QueueDeclare(declare)
	if err != nil {
		return err
	}
	return nil
}

func (s *amqpService) QueueBind(bind QueueBind) error {
	err := s.client.QueueBind(bind)
	if err != nil {
		return err
	}
	return nil
}

// #endregion

// #region IAMQPPublisher Members

func (s *amqpService) Publish(v interface{}, opts ...PublishOption) error {
	err := s.ensurePublishChannelInit()
	if err != nil {
		return err
	}
	publishContext := NewDefaultPublishContext()

	for _, eachOpt := range opts {
		eachOpt(publishContext)
	}
	if publishContext.cancelFunc != nil {
		defer publishContext.cancelFunc()
	}

	data, err := publishContext.Marshal(v)
	if err != nil {
		return fmt.Errorf("cannot serialize object,v: %+V", v)
	}

	return s.client.PublishWithContext(publishContext.ctx,
		publishContext.exchange,
		publishContext.key,
		publishContext.mandatory,
		publishContext.immediate,
		data,
		WithChannel{Channel: s.publishChannel},
	)
}

// #endregion

// #region IAMQPConsumer Members

func (s *amqpService) SimpleConsume(topic string, consumer string, observeFn func(msg *DeliveryMessage)) (ITopicConsumer, error) {
	if topic == "" {
		return nil, fmt.Errorf("topic can not be empty")
	}
	if consumer == "" {
		consumer = s.GenerateUniqueConsumerName()
	}
	channel, err := s.getOrCreateChannel(topic)
	if err != nil {
		return nil, err
	}
	queueConsume := NewDefaultQueueConsume(topic)
	ch, err := s.client.Consume(queueConsume, WithChannel{channel})
	if err != nil {
		return nil, err
	}
	return newDefaultConsumer(consumer, channel, ch, observeFn), nil
}

func (s *amqpService) GenerateUniqueConsumerName() string {
	tagPrefix := "ctag-"
	tagInfix := os.Args[0]
	tagSuffix := "-" + strconv.FormatUint(atomic.AddUint64(&consumerSeq, 1), 10)
	if len(tagPrefix)+len(tagInfix)+len(tagSuffix) > consumerTagLengthMax {
		tagInfix = "abmpio/amqpx"
	}
	return tagPrefix + tagInfix + tagSuffix
}

// #endregion

func (s *amqpService) ensurePublishChannelInit() error {
	if s.publishChannel != nil {
		return nil
	}
	s.client.ensureConnect()
	var err error
	s.publishChannel, err = s.client.GetNewChannel()
	return err
}

func (s *amqpService) getOrCreateChannel(topic string) (*amqp.Channel, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	var err error
	c, ok := s.consumedChannel[topic]
	if !ok {
		c, err = s.client.GetNewChannel()
		if err != nil {
			return nil, err
		}
		s.consumedChannel[topic] = c
	}
	return c, nil
}
