package main

import (
	"os"

	"github.com/cardil/kleio/internal/collector"
)

func main() {
	collector.ServeOrDie(os.Exit)
}
