package varto

type Connection interface {
	Read() ([]byte, error)
	Write([]byte) error
	GetId() string
}
