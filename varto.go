package varto

import "sync"

// Varto is the main struct of the package.
type Varto struct {
	sync.RWMutex
	store *inMemoryStore
}

// New returns a new Varto instance.
func New() *Varto {
	return &Varto{
		store: newInMemoryStore(),
	}
}

func (v *Varto) AddConnection(conn Connection) error {
	v.Lock()
	defer v.Unlock()

	if conn == nil {
		return ErrNilConnection
	}

	return v.store.AddConnection(conn)
}

func (v *Varto) RemoveConnection(conn Connection) error {
	v.Lock()
	defer v.Unlock()

	for _, topic := range v.store.topics {
		topic.Unsubscribe(conn)
	}

	return v.store.RemoveConnection(conn)
}

// Subscribe subscribes a connection to a topic.
func (v *Varto) Subscribe(topicName string, conn Connection) error {
	v.Lock()
	defer v.Unlock()

	if topicName == "" {
		return ErrInvalidTopicName
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

func (v *Varto) Unsubscribe(topic string, conn Connection) error {
	v.Lock()
	defer v.Unlock()

	t, err := v.store.GetTopic(topic)
	if err != nil {
		return err
	}

	t.Unsubscribe(conn)

	return nil
}

// Publish publishes data to a topic.
func (v *Varto) Publish(topic string, data []byte) error {
	t, err := v.store.GetTopic(topic)
	if err != nil {
		return err
	}

	t.Publish(data)
	return nil
}

// BroadcastToAll broadcasts data to all connections.
func (v *Varto) BroadcastToAll(data []byte) error {
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
