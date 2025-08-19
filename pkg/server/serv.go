package server

import (
	"errors"
)

var ErrAlreadyStopped = errors.New("server already stopped")

type Server interface {
	Run() error
	Close() error
}
