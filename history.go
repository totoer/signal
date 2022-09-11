package main

import "signal/sip"

type History struct {
	reqs  []*sip.Request
	resps []*sip.Response
}

func (h *History) writeRequest(req *sip.Request) {
	h.reqs = append(h.reqs, req)
}

func (h *History) writeResponse(resp *sip.Response) {
	h.resps = append(h.resps, resp)
}

func (h *History) topRequest() *sip.Request {
	return h.reqs[len(h.reqs)-1]
}

func (h *History) getInvite() *sip.Request {
	for _, req := range h.reqs {
		if req.Method == sip.INVITE {
			return req
		}
	}
	return nil
}

func NewHistory() *History {
	return &History{
		reqs:  make([]*sip.Request, 0),
		resps: make([]*sip.Response, 0),
	}
}
