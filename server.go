package main

import (
	"context"
	"errors"
	"fmt"
	"signal/db"
	"signal/transport"
	"sync"
	"time"

	"signal/sip"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type UserAgent interface {
	handleRequest(context.Context, string, *sip.Request) error
	handleResponse(context.Context, string, *sip.Response) error
}

type Server struct {
	timeout       int
	db            db.DB
	transport     transport.Transport
	messages      chan sip.Message
	register      *Register
	userAgentPool map[string]UserAgent
}

var ErrWrongRequest = errors.New("wrong request")

func (s *Server) handleRequest(ctx context.Context, cid string, req sip.Request) error {
	log.Info().Str("Call-ID", cid).
		Str("Method", string(req.Method)).
		Str("RURI", req.URI.String()).
		Msg("Handle request")
	switch req.Method {
	case sip.REGISTER:
		s.onRegister(ctx, cid, &req)
	case sip.INVITE:
		s.onInvite(ctx, cid, &req)
	case sip.OPTIONS:
	case sip.INFO:
	default:
		if ua, ok := s.userAgentPool[cid]; ok {
			return ua.handleRequest(ctx, cid, &req)
		} else {
			log.Error().Str("Call-ID", cid).
				Msg("Not found UA for handle request")
			return ErrUnknownUserAgent
		}
	}
	return nil
	// if ua, ok := s.userAgentPool[cid]; ok {
	// 	return ua.handleRequest(ctx, cid, &req)
	// } else if uas, err := NewUAS(cid, s); err != nil {
	// 	return err
	// } else {
	// 	s.userAgentPool[cid] = uas
	// 	return uas.handleRequest(ctx, cid, &req)
	// }
}

var ErrUnknownUserAgent = errors.New("unknown user agent")

func (s *Server) handleResponse(ctx context.Context, cid string, resp sip.Response) error {
	log.Info().Str("Call-ID", cid).
		Str("Code", fmt.Sprint(resp.Code)).
		Str("Text", sip.ResponseCodes[int(resp.Code)]).
		Msg("Handle response")

	if ua, ok := s.userAgentPool[cid]; ok {
		return ua.handleResponse(ctx, cid, &resp)
	} else {
		log.Error().Str("Call-ID", cid).
			Msg("Not found UA for handle response")
		return ErrUnknownUserAgent
	}
}

func (s *Server) serve() {
	for {
		select {
		case m := <-s.messages:
			// if cid, err := m.GetHeaders().GetCallID(); err != nil {
			// 	return
			// } else if from, err := m.GetHeaders().GetFrom(); err != nil {
			// 	return
			// } else if to, err := m.GetHeaders().GetTo(); err != nil {
			// 	return
			// } else {
			// 	transactionKey := fmt.Sprintf("%s:%s:%s", cid, from.Tag, to.Tag)
			// 	if _, ok := s.transactionPool[transactionKey]; ok {
			// 		return
			// 	} else {
			// 		s.transactionPool[transactionKey] = Transaction{}
			// 	}
			// }

			if cid, err := m.GetHeaders().GetCallID(); err != nil {
				log.Error().Err(err).
					Str("body", m.GetRawBody()).
					Msg("Wrong message")
			} else {
				ctx, cencel := context.WithDeadline(context.Background(), time.Now().Add(time.Duration(s.timeout)*time.Second))
				switch m.(type) {
				case sip.Request:
					if err := s.handleRequest(ctx, cid, m.(sip.Request)); err != nil {
						log.Error().Err(err).
							Str("body", m.GetRawBody()).
							Msg("Error while request")
					}

				case sip.Response:
					if err := s.handleResponse(ctx, cid, m.(sip.Response)); err != nil {
						log.Error().Err(err).
							Str("body", m.GetRawBody()).
							Msg("Error while response")
					}
				}
				cencel()
			}
		}
	}
}

func (s *Server) Run() {
	log.Info().Msg("Server run.")
	var wg sync.WaitGroup
	wg.Add(1)
	go s.transport.Run(s.messages)
	go s.serve()
	wg.Wait()
}

func NewServer() (*Server, error) {
	if db, err := db.NewSQLiteDB("./fixtures/store.db"); err != nil {
		return nil, err
	} else {
		host := viper.GetString("server.host")
		port := viper.GetInt("server.port")

		messages := make(chan sip.Message)
		transportType := viper.GetString("server.transport")
		var t transport.Transport
		switch transportType {
		case "UDP":
			t = transport.NewUDPTransport(host, port)
		}

		s := &Server{
			timeout:       viper.GetInt("server.timeout"),
			messages:      messages,
			transport:     t,
			db:            db,
			userAgentPool: make(map[string]UserAgent),
		}

		s.register = NewRegister(s)

		return s, nil
	}
}
