package varto

import "errors"

var ErrTopicNotFound = errors.New("topic not found")
var ErrConnectionNotFound = errors.New("connection not found")
var ErrInvalidTopicName = errors.New("invalid topic name")
var ErrNilConnection = errors.New("connection is nil")
