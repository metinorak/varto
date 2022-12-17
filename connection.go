package varto

type Connection interface {
	Write([]byte) error
}
