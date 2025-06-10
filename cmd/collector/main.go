package main

import (
	"fmt"
	"os"
	"strconv"

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
			fmt.Printf("%#v\n", logParts)
		}
	}(channel)

	server.Wait()
}
