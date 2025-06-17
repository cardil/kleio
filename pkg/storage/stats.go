package storage

import (
	"time"

	"github.com/cardil/kleio/pkg/kubernetes"
)

type Stats []ContainerStat

type ContainerStat struct {
	kubernetes.ContainerInfo `json:"kubernetes"`
	MessageCount             int       `json:"message_count"`
	LastMessage              time.Time `json:"last_message"`
}
