package varto_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/metinorak/varto"
	"github.com/metinorak/varto/mock"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	t.Run("TestNew_WhenNoOptionsProvided_ThenReturnInstance", func(t *testing.T) {
		v := varto.New(nil)
		assert.NotNil(t, v)
	})
}

func TestAddConnection(t *testing.T) {
	t.Run("TestAddConnection_WhenCall_ThenReturnNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().GetId().Return("id")

		err := v.AddConnection(mockConnection)
		assert.Nil(t, err)
	})

	t.Run("TestAddConnection_WhenConnectionIsNil_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		err := v.AddConnection(nil)
		assert.Equal(t, varto.ErrNilConnection, err)
	})
}

func TestRemoveConnection(t *testing.T) {
	t.Run("TestRemoveConnection_WhenCall_ThenReturnNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().GetId().Return("id").AnyTimes()

		v.Subscribe("topic", mockConnection)

		v.AddConnection(mockConnection)
		err := v.RemoveConnection(mockConnection)
		assert.Nil(t, err)
	})
}

func TestSubscribe(t *testing.T) {
	t.Run("TestSubscribe_WhenCall_ThenReturnNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		err := v.Subscribe("topic", mockConnection)
		assert.Nil(t, err)
	})

	t.Run("TestSubscribe_WhenTopicNameIsEmpty_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		err := v.Subscribe("", mockConnection)
		assert.Equal(t, varto.ErrInvalidTopicName, err)
	})

	t.Run("TestSubscribe_WhenTopicIsNotAllowed_ThenReturnError", func(t *testing.T) {
		v := varto.New(&varto.Options{AllowedTopics: []string{"topic"}})
		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		err := v.Subscribe("topic2", mockConnection)
		assert.Equal(t, varto.ErrTopicIsNotAllowed, err)
	})
}

func TestUnsubscribe(t *testing.T) {
	t.Run("TestUnsubscribe_WhenTopicExists_ThenReturnNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		v.Subscribe("topic", mockConnection)
		err := v.Unsubscribe("topic", mockConnection)
		assert.Nil(t, err)
	})

	t.Run("TestUnsubscribe_WhenTopicDoesNotExist_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		err := v.Unsubscribe("topic", mockConnection)
		assert.Equal(t, varto.ErrTopicNotFound, err)
	})
}

func TestPublish(t *testing.T) {
	t.Run("TestPublish_WhenTopicExists_ReturnsNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().Write([]byte("data")).Return(nil)

		v.Subscribe("topic", mockConnection)
		err := v.Publish("topic", []byte("data"))

		time.Sleep(10 * time.Millisecond)
		assert.Nil(t, err)
	})

	t.Run("TestPublish_WhenTopicDoesNotExist_ReturnsError", func(t *testing.T) {
		v := varto.New(nil)
		err := v.Publish("topic", []byte("data"))
		assert.Equal(t, varto.ErrTopicNotFound, err)
	})
}

func TestBroadcastToAll(t *testing.T) {
	t.Run("TestBroadcastToAll_WhenWriteSucceeds_ThenReturnNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().Write(gomock.Any()).Return(nil)
		mockConnection.EXPECT().GetId().Return("id")

		v.AddConnection(mockConnection)
		err := v.BroadcastToAll([]byte("data"))
		assert.Nil(t, err)
	})

	t.Run("TestBroadcastToAll_WhenWriteFails_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().Write(gomock.Any()).Return(fmt.Errorf("error"))
		mockConnection.EXPECT().GetId().Return("id")

		v.AddConnection(mockConnection)
		err := v.BroadcastToAll([]byte("data"))
		assert.NotNil(t, err)
	})
}
