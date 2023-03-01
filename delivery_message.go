package amqpx

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type DeliveryMessage struct {
	_underlyingDelivery *amqp.Delivery
	unmarshal           func([]byte, interface{}) error
}

func newDeliveryMessage(underlyingDelivery *amqp.Delivery) *DeliveryMessage {
	return &DeliveryMessage{
		_underlyingDelivery: underlyingDelivery,
		unmarshal:           _unmarshal,
	}
}

// get Payload data
func (m *DeliveryMessage) Payload() []byte {
	return m._underlyingDelivery.Body
}

// convert Payload to value through json unmarsh
func (m *DeliveryMessage) ToValue(v interface{}) error {
	if len(m._underlyingDelivery.Body) <= 0 {
		return nil
	}
	return _unmarshal(m._underlyingDelivery.Body, v)
}
