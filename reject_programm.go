package main

import (
	"context"
)

type RejectProgrammConfig struct {
	HoldMedia string `json:"media"`
}

type RejectProgramm struct {
	config *RejectProgrammConfig
}

func (r *RejectProgramm) onUASInit(ctx context.Context, uas *UAS) error {
	// sc := NewSignalContext(ctx)

	// if uac, err := NewUAC(); err != nil {
	// 	return err
	// } else if meeting, err := sc.GetMeeting(); err != nil {
	// 	return err
	// } else if meetingPool, err := sc.GetMeetingPool(); err != nil {
	// 	return err
	// } else {
	// 	ce.uas = uas
	// 	ce.uac = uac

	// 	meetingPool.bindOutgoing(meeting.id, uac)
	// 	uac.call()
	// 	uas.ringing()
	// 	meeting.mediaMixer.Join(uas.mediaChanal)

	// 	if ce.config.HoldMedia != "" {
	// 		uas.accept()
	// 	}

	// 	return nil
	// }
	return nil
}

func (r *RejectProgramm) onUASAccepted(ctx context.Context, uas *UAS) error {
	return nil
}

func (r *RejectProgramm) onUASEnded(ctx context.Context, uas *UAS) error {
	return nil
}

func (r *RejectProgramm) onUACTrying(ctx context.Context, uac *UAC) error {
	return nil
}

func (r *RejectProgramm) onUACRinging(ctx context.Context, uac *UAC) error {
	return nil
}

func (r *RejectProgramm) onUACAccepted(ctx context.Context, uac *UAC) error {
	return nil
}

func (r *RejectProgramm) onUACEnded(ctx context.Context, uac *UAC) error {
	return nil
}
