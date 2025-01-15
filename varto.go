package varto

import "sync"

// If this is not provided, default options will be used.
type Options struct {
	// AllowedTopics is a list of topics that are allowed to be subscribed.
	// If this list is empty, all topics are allowed.
	AllowedTopics []string
}

func getDefaultOptions() *Options {
	return &Options{}
}

// Varto is the main struct of the package.
type Varto struct {
	store             Store
	opts              *Options
	allowedTopics     map[string]bool
	middlewareContext *middlewareContext
}

// New returns a new Varto instance.
// If opts is nil, default options will be used.
func New(opts *Options) *Varto {
	return NewWithStore(opts, newInMemoryStore())
}

// NewWithStore creates a new Varto instance with a custom store implementation.
func NewWithStore(opts *Options, store Store) *Varto {
	v := &Varto{
		store:             store,
		middlewareContext: newMiddlewareContext(),
	}

	if opts == nil {
		v.opts = getDefaultOptions()
	} else {
		v.opts = opts
	}

	if len(v.opts.AllowedTopics) > 0 {
		v.allowedTopics = make(map[string]bool)
		for _, topic := range v.opts.AllowedTopics {
			v.allowedTopics[topic] = true
		}
	}

	return v
}

// Use adds a middleware to the middleware chain.
func (v *Varto) Use(middleware Middleware) {
	if middleware == nil {
		return
	}

	v.middlewareContext.Add(middleware)
}

func (v *Varto) AddConnection(conn Connection) error {
	if conn == nil {
		return ErrNilConnection
	}

	for _, m := range v.middlewareContext.GetAll() {
		if err := m.OnAddConnection(conn); err != nil {
			return err
		}
	}

	return v.store.AddConnection(conn)
}

func (v *Varto) RemoveConnection(conn Connection) error {
	if conn == nil {
		return ErrNilConnection
	}

	for _, m := range v.middlewareContext.GetAll() {
		if err := m.OnRemoveConnection(conn); err != nil {
			return err
		}
	}

	return v.store.RemoveConnection(conn)
}

// Subscribe subscribes a connection to a topic.
func (v *Varto) Subscribe(conn Connection, topicName string) error {
	if topicName == "" {
		return ErrInvalidTopicName
	}

	if conn == nil {
		return ErrNilConnection
	}

	for _, m := range v.middlewareContext.GetAll() {
		if err := m.OnSubscribe(conn, topicName); err != nil {
			return err
		}
	}

	if v.allowedTopics != nil {
		if _, ok := v.allowedTopics[topicName]; !ok {
			return ErrTopicIsNotAllowed
		}
	}

	topic, err := v.store.GetTopic(topicName)
	if err == ErrTopicNotFound {
		if t, err := v.store.AddTopic(topicName); err != nil {
			return err
		} else {
			topic = t
		}
	} else if err != nil {
		return err
	}

	topic.Subscribe(conn)
	return nil
}

func (v *Varto) Unsubscribe(conn Connection, topic string) error {
	if topic == "" {
		return ErrInvalidTopicName
	}

	if conn == nil {
		return ErrNilConnection
	}

	for _, m := range v.middlewareContext.GetAll() {
		if err := m.OnUnsubscribe(conn, topic); err != nil {
			return err
		}
	}

	t, err := v.store.GetTopic(topic)
	if err != nil {
		return err
	}

	t.Unsubscribe(conn)

	if t.IsEmpty() {
		if err := v.store.RemoveTopic(topic); err != nil {
			return err
		}
	}

	return nil
}

// Publish publishes data to a topic.
func (v *Varto) Publish(topic string, data []byte) error {
	if topic == "" {
		return ErrInvalidTopicName
	}

	for _, m := range v.middlewareContext.GetAll() {
		if err := m.OnPublish(topic, data); err != nil {
			return err
		}
	}

	t, err := v.store.GetTopic(topic)
	if err != nil {
		return err
	}

	t.Publish(data)
	return nil
}

// BroadcastToAll broadcasts data to all connections.
func (v *Varto) BroadcastToAll(data []byte) error {
	for _, m := range v.middlewareContext.GetAll() {
		if err := m.OnBroadcastToAll(data); err != nil {
			return err
		}
	}

	connections, err := v.store.GetAllConnections()
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}

	chErr := make(chan error, len(connections))
	defer close(chErr)

	for _, conn := range connections {
		wg.Add(1)

		go func(conn Connection) {
			defer wg.Done()
			if err := conn.Write(data); err != nil {
				chErr <- err
			}

		}(conn)
	}

	wg.Wait()

	select {
	case err := <-chErr:
		return err
	default:
	}

	return nil
}
