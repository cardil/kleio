package syslog

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/mcuadros/go-syslog.v2"
)

type Handler func(syslog.LogPartsChannel)

type Waiter interface {
	Wait()
}

type WaiterFn func()

func (f WaiterFn) Wait() {
	f()
}

func Serve(handle Handler) (Waiter, error) {
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
		return nil, err
	}
	if err := server.ListenTCP(bind); err != nil {
		return nil, err
	}

	if err := server.Boot(); err != nil {
		return nil, err
	}

	go handle(channel)

	return WaiterFn(server.Wait), nil
}
