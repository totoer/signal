package main

import (
	"context"
	"signal/sip"

	"github.com/rs/zerolog/log"
)

func (s *Server) onRegister(ctx context.Context, cid string, req *sip.Request) error {
	return s.register.auth(ctx, cid, req, func(ctx context.Context, registration *Registration) error {
		if resp, err := req.MakeResponse(sip.Ok); err != nil {
			log.Err(err).Str("Call-ID", cid).
				Str("where", "UAS.onRegister").
				Msg("While make response")
			return err
		} else if err := s.transport.SendSIP(resp); err != nil {
			log.Err(err).Str("Call-ID", cid).
				Str("where", "UAS.onRegister").
				Msg("While response")
			return err
		}
		return nil
	})
}

func (s *Server) onInvite(ctx context.Context, cid string, req *sip.Request) error {
	return s.register.auth(ctx, cid, req, func(ctx context.Context, registration *Registration) error {
		if uas, err := NewUAS(cid, s); err != nil {
			return err
		} else {
			log.Info().Str("Call-ID", cid).
				Str("where", "UAS.onInvite").
				Msg("Meeting not created")
			if scenario, err := NewScenario(ctx, uas.server, uas.registration.Account.Outgoing); err != nil {
				log.Err(err).Str("Call-ID", cid).
					Str("where", "UAS.onInvite").
					Msg("While create new meeting")
				return err
			} else if meeting, err := NewMeeting(ctx, scenario); err != nil {
				log.Err(err).Str("Call-ID", cid).
					Str("where", "UAS.onInvite").
					Msg("While create new meeting")
				return err
			} else {
				log.Info().Str("Call-ID", cid).
					Str("where", "UAS.onInvite").
					Str("meeting_id", meeting.id.String()).
					Str("scenario_id", scenario.id).
					Msg("Create new meeting")
				uas.meeting.appendUAS(uas)
			}

			if resp, err := req.MakeResponse(sip.Trying); err != nil {
				log.Err(err).Str("Call-ID", cid).
					Str("meeting_id", uas.meeting.id.String()).
					Str("scenario_id", uas.meeting.scenario.id).
					Str("response", sip.ResponseCodes[int(sip.Trying)]).
					Msg("While create response")
				return err
			} else if err := uas.server.transport.SendSIP(resp); err != nil {
				log.Err(err).Str("Call-ID", cid).
					Str("meeting_id", uas.meeting.id.String()).
					Str("scenario_id", uas.meeting.scenario.id).
					Str("response", sip.ResponseCodes[int(sip.Trying)]).
					Msg("While send response")
				return err
			} else {
				return uas.meeting.scenario.run(ctx, uas)
			}
		}
	})
}
