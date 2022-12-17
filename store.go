package varto

import "sync"

type inMemoryStore struct {
	sync.RWMutex
	connections map[string]Connection
	topics      map[string]Topic
}

func newInMemoryStore() *inMemoryStore {
	return &inMemoryStore{
		connections: make(map[string]Connection),
		topics:      make(map[string]Topic),
	}
}

func (s *inMemoryStore) GetTopic(topicName string) (Topic, error) {
	s.Lock()
	defer s.Unlock()

	if _, ok := s.topics[topicName]; !ok {
		return nil, ErrTopicNotFound
	}

	return s.topics[topicName], nil
}

func (s *inMemoryStore) AddTopic(topicName string) (Topic, error) {
	s.Lock()
	defer s.Unlock()

	newTopic := NewTopic(topicName)
	s.topics[topicName] = newTopic
	return newTopic, nil
}

func (s *inMemoryStore) RemoveTopic(topicName string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.topics, topicName)
	return nil
}

func (s *inMemoryStore) AddConnection(conn Connection) error {
	s.Lock()
	defer s.Unlock()

	s.connections[conn.GetId()] = conn

	return nil
}

func (s *inMemoryStore) RemoveConnection(conn Connection) error {
	s.Lock()
	defer s.Unlock()

	delete(s.connections, conn.GetId())

	for _, topic := range s.topics {
		topic.Unsubscribe(conn)
	}

	return nil
}

func (s *inMemoryStore) GetAllConnections() ([]Connection, error) {
	s.RLock()
	defer s.RUnlock()

	var connections []Connection

	for conn := range s.connections {
		connections = append(connections, s.connections[conn])
	}

	return connections, nil
}
