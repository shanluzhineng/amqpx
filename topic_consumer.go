package amqpx

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// IConsumer define a amqp queue consumer
type ITopicConsumer interface {
	// observe a message
	Observe(fn func(msg *DeliveryMessage))

	// stop consume
	Stop() error
}

type defaultTopicConsumer struct {
	consumer string
	channel  *amqp.Channel
	ch       <-chan amqp.Delivery

	unmarshal       func([]byte, interface{}) error
	registedObserve []func(msg *DeliveryMessage)
}

var _ ITopicConsumer = (*defaultTopicConsumer)(nil)

func newDefaultConsumer(consumer string,
	channel *amqp.Channel,
	ch <-chan amqp.Delivery,
	observeFn func(msg *DeliveryMessage)) *defaultTopicConsumer {
	defaultConsumer := &defaultTopicConsumer{
		consumer:  consumer,
		channel:   channel,
		ch:        ch,
		unmarshal: _unmarshal,
	}

	if observeFn != nil {
		defaultConsumer.Observe(observeFn)
	}
	defaultConsumer.start()
	return defaultConsumer
}

// #region IConsumer Members

func (c *defaultTopicConsumer) Observe(fn func(msg *DeliveryMessage)) {
	c.registedObserve = append(c.registedObserve, fn)
}

func (c *defaultTopicConsumer) Stop() error {
	err := c.channel.Cancel(c.consumer, true)
	if err != nil {
		return err
	}
	return nil
}

// #endregion

func (c *defaultTopicConsumer) start() {
	go c.safeStart()
}

func (c *defaultTopicConsumer) safeStart() {
	defer func() {
		if p := recover(); p != nil {
			c.Stop()
		}
	}()
	for eachDelivery := range c.ch {
		c.notifyObserver(&eachDelivery)
	}
}

func (c *defaultTopicConsumer) notifyObserver(deliveryMessage *amqp.Delivery) {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("defaultConsumer.notifyObserver panic when notify registed observer, panic: %v", p)
		}
	}()

	clonedObserver := make([]func(msg *DeliveryMessage), len(c.registedObserve))
	copy(clonedObserver, c.registedObserve)

	newMessage := newDeliveryMessage(deliveryMessage)
	newMessage.unmarshal = c.unmarshal
	for _, eachObserver := range clonedObserver {
		eachObserver(newMessage)
	}
}
