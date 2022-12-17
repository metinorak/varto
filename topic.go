package varto

import (
	"fmt"
	"sync"
)

type Topic interface {
	Subscribe(Connection)
	Unsubscribe(Connection)
	IsEmpty() bool
	Publish([]byte)
}

type topic struct {
	sync.RWMutex
	name        string
	connections map[Connection]bool
	Channel     chan []byte
}

func NewTopic(name string) Topic {
	t := &topic{
		name:        name,
		connections: make(map[Connection]bool),
		Channel:     make(chan []byte, 100),
	}

	go t.listen()

	return t
}

func (t *topic) Subscribe(conn Connection) {
	t.Lock()
	defer t.Unlock()

	t.connections[conn] = true
}

func (t *topic) Unsubscribe(conn Connection) {
	t.Lock()
	defer t.Unlock()

	delete(t.connections, conn)
}

func (t *topic) IsEmpty() bool {
	t.RLock()
	defer t.RUnlock()

	return len(t.connections) == 0
}

func (t *topic) Publish(data []byte) {
	t.Channel <- data
}

func (t *topic) listen() {
	for data := range t.Channel {
		if err := t.publish(data); err != nil {
			fmt.Println(err)
		}
	}
}

func (t *topic) publish(data []byte) error {
	t.RLock()
	connections := t.connections
	t.RUnlock()

	wg := sync.WaitGroup{}
	chErr := make(chan error)

	for conn := range connections {
		wg.Add(1)

		go func(c Connection) {
			defer wg.Done()

			if err := c.Write(data); err != nil {
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
