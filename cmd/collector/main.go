package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"gopkg.in/mcuadros/go-syslog.v2"
)

func main() {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	server := syslog.NewServer()
	server.SetFormat(syslog.RFC5424)
	server.SetHandler(handler)
	port := 514
	if eport, set := os.LookupEnv("PORT"); set {
		iport, perr := strconv.Atoi(eport)
		if perr == nil {
			port = iport
		}
	}
	bind := fmt.Sprint("0.0.0.0:", port)
	if err := server.ListenUDP(bind); err != nil {
		panic(err)
	}
	if err := server.ListenTCP(bind); err != nil {
		panic(err)
	}

	if err := server.Boot(); err != nil {
		panic(err)
	}

	go func(channel syslog.LogPartsChannel) {
		for logParts := range channel {
			data := ClusterLoggingMsg{}
			if err := json.Unmarshal([]byte(fmt.Sprint(logParts["message"])), &data); err != nil {
				panic(err)
			}
			bytes, _ := json.MarshalIndent(data, "", "  ")
			fmt.Printf("%s\n", string(bytes))
		}
	}(channel)

	server.Wait()
}

type Kubernetes struct {
	ContainerName  string `json:"container_name"`
	ContainerImage string `json:"container_image"`
	NamespaceName  string `json:"namespace_name"`
	PodName        string `json:"pod_name"`
}

type ClusterLoggingMsg struct {
	Timestamp  time.Time `json:"timestamp"`
	Message    string    `json:"message"`
	Kubernetes `json:"kubernetes"`
}
