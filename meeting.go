package main

import (
	"context"
	"errors"
	"signal/media"

	"github.com/google/uuid"
)

var ErrUnknownMeeting = errors.New("unknown meeting")

type Meeting struct {
	id         uuid.UUID
	scenario   *Scenario
	uasPool    map[string]*UAS
	uacPool    map[string]*UAC
	mediaMixer *media.MediaMixer
}

func (m *Meeting) appendUAS(uas *UAS) {
	m.uasPool[uas.callID] = uas
	m.uasPool[uas.callID].meeting = m
}

func (m *Meeting) appendUAC(uac *UAC) {
	m.uacPool[uac.callID] = uac
	m.uacPool[uac.callID].meeting = m
}

var ErrMeetingHaseNotLag = errors.New("meeting hase not lag")

func NewMeeting(ctx context.Context, s *Scenario) (*Meeting, error) {
	if mm, err := media.NewMediaMixer(); err != nil {
		return nil, err
	} else {
		m := &Meeting{
			id:         uuid.New(),
			scenario:   s,
			uasPool:    make(map[string]*UAS),
			uacPool:    make(map[string]*UAC),
			mediaMixer: mm,
		}
		s.meeting = m
		return m, nil
	}
}
