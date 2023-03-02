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

// Ack delegates an acknowledgement through the Acknowledger interface that the
// client or server has finished work on a delivery.
//
// All deliveries in AMQP must be acknowledged.  If you called Channel.Consume
// with autoAck true then the server will be automatically ack each message and
// this method should not be called.  Otherwise, you must call Delivery.Ack after
// you have successfully processed this delivery.
//
// When multiple is true, this delivery and all prior unacknowledged deliveries
// on the same channel will be acknowledged.  This is useful for batch processing
// of deliveries.
//
// An error will indicate that the acknowledge could not be delivered to the
// channel it was sent from.
//
// Either Delivery.Ack, Delivery.Reject or Delivery.Nack must be called for every
// delivery that is not automatically acknowledged.
func (m *DeliveryMessage) Ack(multiple bool) error {
	return m._underlyingDelivery.Ack(multiple)
}

// Reject delegates a negatively acknowledgement through the Acknowledger interface.
//
// When requeue is true, queue this message to be delivered to a consumer on a
// different channel.  When requeue is false or the server is unable to queue this
// message, it will be dropped.
//
// If you are batch processing deliveries, and your server supports it, prefer
// Delivery.Nack.
//
// Either Delivery.Ack, Delivery.Reject or Delivery.Nack must be called for every
// delivery that is not automatically acknowledged.
func (m *DeliveryMessage) Reject(requeue bool) error {
	return m._underlyingDelivery.Reject(requeue)
}

// Nack negatively acknowledge the delivery of message(s) identified by the
// delivery tag from either the client or server.
//
// When multiple is true, nack messages up to and including delivered messages up
// until the delivery tag delivered on the same channel.
//
// When requeue is true, request the server to deliver this message to a different
// consumer.  If it is not possible or requeue is false, the message will be
// dropped or delivered to a server configured dead-letter queue.
//
// This method must not be used to select or requeue messages the client wishes
// not to handle, rather it is to inform the server that the client is incapable
// of handling this message at this time.
//
// Either Delivery.Ack, Delivery.Reject or Delivery.Nack must be called for every
// delivery that is not automatically acknowledged.
func (m *DeliveryMessage) Nack(multiple, requeue bool) error {
	return m._underlyingDelivery.Nack(multiple, requeue)
}
