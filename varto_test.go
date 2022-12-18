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

	t.Run("TestAddConnection_WhenMiddlewareReturnsError_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		mockMiddleware := mock.NewMockMiddleware(gomock.NewController(t))
		mockMiddleware.EXPECT().OnAddConnection(mockConnection).Return(fmt.Errorf("error"))

		v.Use(mockMiddleware)

		err := v.AddConnection(mockConnection)
		assert.NotNil(t, err)
	})
}

func TestRemoveConnection(t *testing.T) {
	t.Run("TestRemoveConnection_WhenCall_ThenReturnNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().GetId().Return("id").AnyTimes()

		v.Subscribe(mockConnection, "topic")

		v.AddConnection(mockConnection)
		err := v.RemoveConnection(mockConnection)
		assert.Nil(t, err)
	})

	t.Run("TestRemoveConnection_WhenConnectionIsNil_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		err := v.RemoveConnection(nil)
		assert.Equal(t, varto.ErrNilConnection, err)
	})

	t.Run("TestRemoveConnection_WhenMiddlewareReturnsError_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)

		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().GetId().Return("id").AnyTimes()

		mockMiddleware := mock.NewMockMiddleware(gomock.NewController(t))
		mockMiddleware.EXPECT().OnAddConnection(mockConnection).Return(nil)
		mockMiddleware.EXPECT().OnRemoveConnection(mockConnection).Return(fmt.Errorf("error"))

		v.Use(mockMiddleware)

		v.AddConnection(mockConnection)
		err := v.RemoveConnection(mockConnection)
		assert.NotNil(t, err)
	})
}

func TestSubscribe(t *testing.T) {
	t.Run("TestSubscribe_WhenCall_ThenReturnNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().GetId().Return("id")

		err := v.Subscribe(mockConnection, "topic")
		assert.Nil(t, err)
	})

	t.Run("TestSubscribe_WhenConnectionIsNil_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		err := v.Subscribe(nil, "topic")
		assert.Equal(t, varto.ErrNilConnection, err)
	})

	t.Run("TestSubscribe_WhenTopicNameIsEmpty_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		err := v.Subscribe(mockConnection, "")
		assert.Equal(t, varto.ErrInvalidTopicName, err)
	})

	t.Run("TestSubscribe_WhenTopicIsNotAllowed_ThenReturnError", func(t *testing.T) {
		v := varto.New(&varto.Options{AllowedTopics: []string{"topic"}})
		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		err := v.Subscribe(mockConnection, "topic2")
		assert.Equal(t, varto.ErrTopicIsNotAllowed, err)
	})

	t.Run("TestSubscribe_WhenSuccessfullySubscribed_ThenShouldPublish", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().Write([]byte("data")).Return(nil)
		mockConnection.EXPECT().GetId().Return("id").AnyTimes()

		v.Subscribe(mockConnection, "topic")
		err := v.Publish("topic", []byte("data"))

		time.Sleep(10 * time.Millisecond)
		assert.Nil(t, err)
	})

	t.Run("TestSubscribe_WhenMiddlewareReturnsError_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)

		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		mockMiddleware := mock.NewMockMiddleware(gomock.NewController(t))
		mockMiddleware.EXPECT().OnSubscribe(mockConnection, "topic").Return(fmt.Errorf("error")).AnyTimes()

		v.Use(mockMiddleware)

		err := v.Subscribe(mockConnection, "topic")
		assert.NotNil(t, err)
	})
}

func TestUnsubscribe(t *testing.T) {
	t.Run("TestUnsubscribe_WhenTopicExists_ThenReturnNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().GetId().Return("id").AnyTimes()

		v.Subscribe(mockConnection, "topic")
		err := v.Unsubscribe(mockConnection, "topic")
		assert.Nil(t, err)
	})

	t.Run("TestUnsubscribe_WhenTopicDoesNotExist_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().GetId().Return("id").AnyTimes()

		err := v.Unsubscribe(mockConnection, "topic")
		assert.Equal(t, varto.ErrTopicNotFound, err)
	})

	t.Run("TestUnsubscribe_WhenConnectionIsNil_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		err := v.Unsubscribe(nil, "topic")
		assert.Equal(t, varto.ErrNilConnection, err)
	})

	t.Run("TestUnsubscribe_WhenTopicNameIsEmpty_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		err := v.Unsubscribe(mockConnection, "")
		assert.Equal(t, varto.ErrInvalidTopicName, err)
	})

	t.Run("TestUnsubscribe_WhenSuccessfullyUnsubscribed_ThenShouldNotPublish", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().GetId().Return("id").AnyTimes()
		mockConnection.EXPECT().Write(gomock.Any()).Times(0)

		v.Subscribe(mockConnection, "topic")
		v.Unsubscribe(mockConnection, "topic")
		err := v.Publish("topic", []byte("data"))

		assert.Equal(t, varto.ErrTopicNotFound, err)
	})

	t.Run("TestUnsubscribe_WhenMiddlewareReturnsError_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)

		mockConnection := mock.NewMockConnection(gomock.NewController(t))

		mockMiddleware := mock.NewMockMiddleware(gomock.NewController(t))
		mockMiddleware.EXPECT().OnUnsubscribe(mockConnection, "topic").Return(fmt.Errorf("error")).AnyTimes()

		v.Use(mockMiddleware)

		err := v.Unsubscribe(mockConnection, "topic")
		assert.NotNil(t, err)
	})
}

func TestPublish(t *testing.T) {
	t.Run("TestPublish_WhenTopicExists_ReturnsNil", func(t *testing.T) {
		v := varto.New(nil)
		mockConnection := mock.NewMockConnection(gomock.NewController(t))
		mockConnection.EXPECT().Write([]byte("data")).Return(nil)
		mockConnection.EXPECT().GetId().Return("id").AnyTimes()

		v.Subscribe(mockConnection, "topic")
		err := v.Publish("topic", []byte("data"))

		time.Sleep(10 * time.Millisecond)
		assert.Nil(t, err)
	})

	t.Run("TestPublish_WhenTopicDoesNotExist_ReturnsError", func(t *testing.T) {
		v := varto.New(nil)
		err := v.Publish("topic", []byte("data"))
		assert.Equal(t, varto.ErrTopicNotFound, err)
	})

	t.Run("TestPublish_WhenMiddlewareReturnsError_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)

		mockMiddleware := mock.NewMockMiddleware(gomock.NewController(t))
		mockMiddleware.EXPECT().OnPublish("topic", []byte("data")).Return(fmt.Errorf("error")).AnyTimes()

		v.Use(mockMiddleware)

		err := v.Publish("topic", []byte("data"))
		assert.NotNil(t, err)
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

	t.Run("TestBroadcastToAll_WhenMiddlewareReturnsError_ThenReturnError", func(t *testing.T) {
		v := varto.New(nil)

		mockMiddleware := mock.NewMockMiddleware(gomock.NewController(t))
		mockMiddleware.EXPECT().OnBroadcastToAll([]byte("data")).Return(fmt.Errorf("error")).AnyTimes()

		v.Use(mockMiddleware)

		err := v.BroadcastToAll([]byte("data"))
		assert.NotNil(t, err)
	})
}
