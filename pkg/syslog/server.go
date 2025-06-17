package syslog

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/cardil/kleio/pkg/server"
	"gopkg.in/mcuadros/go-syslog.v2"
)

var ErrSyslogInit = errors.New("Syslog init failed")

type Handler func(syslog.LogPartsChannel)

func Serve(handle Handler) server.Server {
	channel := make(syslog.LogPartsChannel)
	handler := syslog.NewChannelHandler(channel)

	srv := syslog.NewServer()
	srv.SetFormat(syslog.RFC5424)
	srv.SetHandler(handler)
	port := 8514
	if eport, set := os.LookupEnv("PORT"); set {
		iport, perr := strconv.Atoi(eport)
		if perr == nil {
			port = iport
		}
	}
	bind := fmt.Sprint("0.0.0.0:", port)

	return &syslogServer{
		server:  srv,
		bind:    bind,
		channel: channel,
		handler: handle,
	}
}

type syslogServer struct {
	server  *syslog.Server
	bind    string
	handler Handler
	channel syslog.LogPartsChannel
}

func (s *syslogServer) Run() (err error) {
	slog.Info("Starting Syslog server", "bind", s.bind)

	defer func() {
		if err != nil {
			err = fmt.Errorf("%w: %w", ErrSyslogInit, err)
		}
	}()
	if err = s.server.ListenUDP(s.bind); err != nil {
		return
	}
	if err = s.server.ListenTCP(s.bind); err != nil {
		return
	}

	if err = s.server.Boot(); err != nil {
		return
	}

	go s.handler(s.channel)

	s.server.Wait()
	return nil
}

func (s *syslogServer) Kill() error {
	return s.server.Kill()
}
