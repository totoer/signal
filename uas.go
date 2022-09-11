package main

import (
	"context"
	"fmt"
	"signal/media"
	"signal/sip"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

type UAS struct {
	callID       string
	tag          string
	server       *Server
	meeting      *Meeting
	registration *Registration
	history      *History
	mediaChanal  *media.MediaChanal
}

func (uas *UAS) handleRequest(ctx context.Context, cid string, req *sip.Request) error {
	uas.history.writeRequest(req)

	switch req.Method {
	case sip.ACK:
		uas.meeting.scenario.uasEmit(UAS_READY, ctx, uas)
	case sip.CANCEL, sip.BYE:
		uas.sendResponse(sip.Ok, nil)
		uas.meeting.scenario.uasEmit(UAS_END, ctx, uas)
	}
	return nil
}

func (uas *UAS) handleResponse(ctx context.Context, cid string, resp *sip.Response) error {
	uas.history.writeResponse(resp)
	switch resp.Code {
	case sip.Ok:
		uas.meeting.scenario.uasEmit(UAS_END, ctx, uas)
	}
	return nil
}

func (uas *UAS) getBaseHeaders() sip.Headers {
	headers := sip.NewHeaders(nil)
	invite := uas.history.getInvite()
	copy(headers.Vias, invite.Headers.Vias)
	host := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))
	headers.PushVia(sip.Via{
		Host:     host,
		Branch:   uuid.NewString(),
		Received: "",
		Rport:    false,
	})
	to, _ := invite.GetHeaders().GetTo()
	to.Tag = uas.tag
	headers.To = &to
	from, _ := invite.GetHeaders().GetFrom()
	headers.From = &from
	headers.CallID = &sip.PlainHeader{
		Value: uas.callID,
	}
	topRequest := uas.history.topRequest()
	cseq, _ := topRequest.GetHeaders().GetCSeq()
	headers.CSeq = &cseq
	headers.ContentLength = &sip.IntegerHeader{
		Value: 0,
	}

	return headers
}

func (uas *UAS) sendRequest(m sip.MethodType, f func(sip.Request) sip.Request) error {
	headers := uas.getBaseHeaders()

	invite := uas.history.getInvite()
	from, _ := invite.GetHeaders().GetFrom()

	req := sip.NewRequest(m, "", from.Address.URI, headers)

	if f != nil {
		req = f(req)
	}

	if err := uas.server.transport.SendSIP(req); err != nil {
		return err
	}
	return nil
}

func (uas *UAS) sendResponse(c sip.ResponseCode, f func(sip.Response) sip.Response) error {
	headers := uas.getBaseHeaders()

	resp := sip.NewResponse(c, sip.ResponseCodes[int(c)], headers)

	if f != nil {
		resp = f(resp)
	}

	if err := uas.server.transport.SendSIP(resp); err != nil {
		return err
	}
	return nil
}

func (uas *UAS) trying() error {
	return uas.sendResponse(sip.Trying, nil)
}

func (uas *UAS) ringing() error {
	return uas.sendResponse(sip.Ringing, nil)
}

func (uas *UAS) accept() error {
	return uas.sendResponse(sip.Ok, nil)
}

func (uas *UAS) bye() error {
	return uas.sendRequest(sip.BYE, nil)
}

func (uas *UAS) cencel() error {
	return uas.sendRequest(sip.CANCEL, nil)
}

func NewUAS(cid string, s *Server) (*UAS, error) {
	uas := &UAS{
		callID:  cid,
		tag:     uuid.NewString(),
		server:  s,
		history: NewHistory(),
	}

	return uas, nil
}
