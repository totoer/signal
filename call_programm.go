package main

import (
	"context"
	"signal/sip"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type CallProgrammConfig struct {
	Target   *sip.Destination `json:"target"`
	Greeting string           `json:"greeting"`
}

type CallProgramm struct {
	uas    *UAS
	uac    *UAC
	config *CallProgrammConfig
}

func (cp *CallProgramm) isGreeting() bool {
	return cp.config.Greeting != ""
}

func (cp *CallProgramm) init(ctx context.Context, m *Meeting, uas *UAS) error {
	m.scenario.onUASEvent(UAS_READY, cp.onUASReady)
	m.scenario.onUASEvent(UAS_END, cp.onUASEnd)

	m.scenario.onUACEvent(UAC_RINGING, cp.onUACRinging)
	m.scenario.onUACEvent(UAC_READY, cp.onUACReady)
	m.scenario.onUACEvent(UAC_END, cp.onUACEnd)

	cp.uas = uas
	log.Info().Str("Call-ID", uas.callID).
		Str("where", "CallProgramm.init").
		Msg("init")

	invite := uas.history.getInvite()

	var target sip.Destination
	if cp.config.Target != nil {
		target = *cp.config.Target
	} else {
		target, _ = invite.GetHeaders().GetTo()
	}

	if registration, err := uas.server.register.loadRegistrationByDestination(ctx, target); err != nil {
		log.Error().Err(err).Str("Call-ID", uas.callID).
			Str("where", "CallProgramm.init").
			Msg("While init")
		return err
	} else if uac, err := NewUAC(uuid.NewString(), uas.server, registration); err != nil {
		log.Error().Err(err).Str("Call-ID", uas.callID).
			Str("where", "CallProgramm.init").
			Msg("While init")
		return err
	} else {
		cp.uac = uac
		uas.meeting.appendUAC(uac)
		uas.server.userAgentPool[uac.callID] = uac
		if cp.isGreeting() {
			return cp.uas.accept()
		} else {
			from, _ := invite.GetHeaders().GetFrom()
			return uac.call(from)
		}
	}
}

func (cp *CallProgramm) onUASReady(ctx context.Context, uas *UAS) {
	if cp.isGreeting() {
		cp.uas.mediaChanal.Play()
		cp.uac.mediaChanal.Beeps()
	} else {
		cp.uac.accept()
		cp.uas.meeting.mediaMixer.Join(cp.uas.mediaChanal)
		cp.uac.meeting.mediaMixer.Join(cp.uac.mediaChanal)
	}
}

func (cp *CallProgramm) onUASEnd(ctx context.Context, uas *UAS) {}

func (cp *CallProgramm) onUACRinging(ctx context.Context, uac *UAC) {}

func (cp *CallProgramm) onUACReady(ctx context.Context, uac *UAC) {}

func (cp *CallProgramm) onUACEnd(ctx context.Context, uac *UAC) {}
