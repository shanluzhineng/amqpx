package amqpx

import (
	"context"
	"time"
)

const (
	_defaultTimeout = 5 * time.Second
)

type (
	UnmarshalFunc func([]byte, interface{}) error
	MarshalFunc   func(interface{}) ([]byte, error)
)

type PublishContext struct {
	ctx        context.Context
	cancelFunc context.CancelFunc

	exchange  string
	key       string
	mandatory bool
	immediate bool

	Unmarshal UnmarshalFunc
	Marshal   MarshalFunc
}

type PublishOption func(c *PublishContext)

func NewDefaultPublishContext() *PublishContext {
	ctx, cancelFunc := context.WithTimeout(context.Background(), _defaultTimeout)
	return &PublishContext{
		ctx:        ctx,
		cancelFunc: cancelFunc,
		Marshal:    _marshal,
		Unmarshal:  _unmarshal,
	}
}

func WithContext(ctx context.Context, cancelFunc context.CancelFunc) PublishOption {
	return func(c *PublishContext) {
		c.ctx = ctx
		c.cancelFunc = cancelFunc
	}
}

func WithExchange(exchange string) PublishOption {
	return func(c *PublishContext) {
		c.exchange = exchange
	}
}

func WithKey(key string) PublishOption {
	return func(c *PublishContext) {
		c.key = key
	}
}

func WithMandatory(mandatory bool) PublishOption {
	return func(c *PublishContext) {
		c.mandatory = mandatory
	}
}

func WithImmediate(immediate bool) PublishOption {
	return func(c *PublishContext) {
		c.immediate = immediate
	}
}
