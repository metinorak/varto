package varto

import "sync"

type Middleware interface {
	// OnAddConnection is called when a new connection is added.
	OnAddConnection(conn Connection) error

	// OnRemoveConnection is called when a connection is removed.
	OnRemoveConnection(conn Connection) error

	// OnSubscribe is called when a connection subscribes to a topic.
	OnSubscribe(conn Connection, topic string) error

	// OnUnsubscribe is called when a connection unsubscribes from a topic.
	OnUnsubscribe(conn Connection, topic string) error

	// OnPublish is called when a connection publishes a message to a topic.
	OnPublish(topic string, data []byte) error

	// OnBrodcastToAll is called when a connection publishes a message to everyone.
	OnBroadcastToAll(data []byte) error
}

type middlewareContext struct {
	sync.RWMutex
	items []Middleware
}

func newMiddlewareContext() *middlewareContext {
	return &middlewareContext{}
}

func (c *middlewareContext) Add(middleware Middleware) {
	if middleware == nil {
		return
	}

	c.Lock()
	defer c.Unlock()

	c.items = append(c.items, middleware)
}

func (c *middlewareContext) GetAll() []Middleware {
	c.RLock()
	defer c.RUnlock()

	return c.items
}
