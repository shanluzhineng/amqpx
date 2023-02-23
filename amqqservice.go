package amqpx

import "fmt"

type IAMQPService interface {
	IAMQPPublisher

	ExchangeDeclare(declare ExchangeDeclare) error
	QueueDeclare(declare QueueDeclare) error
	QueueBind(bind QueueBind) error
}

type IAMQPPublisher interface {
	Publish(v interface{}, opts ...PublishOption) error
}

type amqpService struct {
	client *AMQPClient
}

// create IAMQPService instance
func NewAMQPService(client *AMQPClient) IAMQPService {
	return &amqpService{
		client: client,
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
	)
}

// #endregion
