package main

import (
	"context"
	"fmt"
	"signal/media"
	"signal/sip"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type UAC struct {
	callID       string
	tag          string
	toTag        string
	from         sip.Destination
	server       *Server
	meeting      *Meeting
	registration *Registration
	history      *History
	mediaChanal  *media.MediaChanal
}

func (uac *UAC) handleRequest(ctx context.Context, cid string, req *sip.Request) error {
	uac.history.writeRequest(req)
	switch req.Method {
	case sip.BYE:
		uac.sendResponse(sip.Ok, nil)
		uac.meeting.scenario.uacEmit(UAC_END, ctx, uac)
	}
	return nil
}

func (uac *UAC) handleResponse(ctx context.Context, cid string, resp *sip.Response) error {
	uac.history.writeResponse(resp)

	if to, err := resp.GetHeaders().GetTo(); err == nil && to.Tag != "" {
		uac.toTag = to.Tag
	}

	switch resp.Code {
	case sip.Trying:
		log.Info().Str("Call-ID", uac.callID).
			Str("where", "UAC.onTrying").
			Str("meeting_id", uac.meeting.id.String()).
			Msg("Trying received")
	case sip.Ringing:
		log.Info().Str("Call-ID", uac.callID).
			Str("where", "UAC.onRinging").
			Str("meeting_id", uac.meeting.id.String()).
			Msg("Ringing received")
		uac.meeting.scenario.uacEmit(UAC_RINGING, ctx, uac)
	case sip.Ok:
		log.Info().Str("Call-ID", uac.callID).
			Str("where", "UAC.onRinging").
			Str("meeting_id", uac.meeting.id.String()).
			Msg("Ringing received")
		uac.meeting.scenario.uacEmit(UAC_READY, ctx, uac)
	}
	return nil
}

func (uac *UAC) getBaseHeaders() sip.Headers {
	headers := sip.NewHeaders(nil)
	headers.CallID = &sip.PlainHeader{
		Value: uac.callID,
	}
	headers.From = &uac.from
	headers.To = &uac.registration.Destination
	headers.To.Tag = uac.toTag
	headers.Vias = make([]sip.Via, 0)
	host := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))
	headers.PushVia(sip.Via{
		Host:     host,
		Branch:   uuid.NewString(),
		Received: "",
		Rport:    false,
	})
	headers.CSeq = &sip.CSeq{
		Value:  0,
		Method: sip.INVITE,
	}
	headers.ContentLength = &sip.IntegerHeader{
		Value: 0,
	}

	return headers
}

func (uac *UAC) sendRequest(m sip.MethodType, f func(sip.Request) sip.Request) error {
	return nil
}

func (uac *UAC) sendResponse(c sip.ResponseCode, f func(sip.Response) sip.Response) error {
	return nil
}

func (uac *UAC) call(from sip.Destination) error {
	log.Info().Str("Call-ID", uac.callID).
		Str("where", "UAC.call").
		Msg("Start call")
	h := sip.NewHeaders(nil)
	h.CallID = &sip.PlainHeader{
		Value: uac.callID,
	}
	h.Contacts = uac.registration.Contacts
	h.From = &from
	h.To = &uac.registration.Destination
	h.To.Tag = ""
	h.Vias = make([]sip.Via, 0)
	host := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))
	h.PushVia(sip.Via{
		Host:     host,
		Branch:   uuid.NewString(),
		Received: "",
		Rport:    false,
	})
	h.ContentLength = &sip.IntegerHeader{
		Value: 0,
	}
	h.CSeq = &sip.CSeq{
		Value:  0,
		Method: sip.INVITE,
	}
	h.Contacts = make([]sip.Contact, 0)
	h.Contacts = append(h.Contacts, sip.Contact{
		Address: sip.Address{
			URI: sip.URI{
				Host:  host,
				Login: from.Address.URI.Login,
			},
		},
	})

	req := sip.NewRequest(sip.INVITE, "", h.To.Address.URI, h)
	req.SourceAddres = uac.registration.SourceAddres
	if err := uac.server.transport.SendSIP(req); err != nil {
		log.Error().Err(err).Str("Call-ID", uac.callID).
			Str("where", "UAC.call").
			Msg("While start call")
		return err
	} else {
		return nil
	}
}

func (uac *UAC) accept() error {
	return uac.sendRequest(sip.ACK, nil)
}

func NewUAC(cid string, s *Server, r *Registration) (*UAC, error) {
	uac := &UAC{
		callID:       cid,
		server:       s,
		registration: r,
		history:      NewHistory(),
	}

	return uac, nil
}
