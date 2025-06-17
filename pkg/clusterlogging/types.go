package clusterlogging

import (
	"time"

	"github.com/cardil/kleio/pkg/kubernetes"
)

type Message struct {
	Timestamp                time.Time `json:"timestamp"`
	Message                  string    `json:"message"`
	kubernetes.ContainerInfo `json:"kubernetes"`
}
