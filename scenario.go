package main

import (
	"context"
	"errors"
	"fmt"
	"signal/sip"
)

type ProgrammType string

const (
	CallProgrammType   ProgrammType = "CALL_PROGRAMM"
	RejectProgrammType ProgrammType = "REJECT_PROGRAMM"
)

type RequestEvent struct {
	method sip.MethodType
}

type ResponseEvent struct {
	code sip.ResponseCode
}
type Programm interface {
	init(context.Context, *Meeting, *UAS) error
}

type ScenarioConfig struct {
	ID        string                  `json:"id"`
	RootID    string                  `json:"root_id"`
	Programms map[string]ProgrammType `json:"programms"`
}

type UASEvent int

const (
	UAS_READY UASEvent = iota
	UAS_END
)

type UACEvent int

const (
	UAC_RINGING UACEvent = iota
	UAC_READY
	UAC_END
)

type UASEventHandler func(context.Context, *UAS)

type UACEventHandler func(context.Context, *UAC)

type Scenario struct {
	id               string
	meeting          *Meeting
	rootId           string
	programm         Programm
	programms        map[string]Programm
	uasEventHandlers map[UASEvent]UASEventHandler
	uacEventHandlers map[UACEvent]UACEventHandler
}

var ErrUnknownProgramm = errors.New("unknown programm")

func (s *Scenario) onUASEvent(event UASEvent, f UASEventHandler) {
	s.uasEventHandlers[event] = f
}

func (s *Scenario) onUACEvent(event UACEvent, f UACEventHandler) {
	s.uacEventHandlers[event] = f
}

func (s *Scenario) uasEmit(event UASEvent, ctx context.Context, uas *UAS) {
	if f, ok := s.uasEventHandlers[event]; ok {
		f(ctx, uas)
	}
}

func (s *Scenario) uacEmit(event UACEvent, ctx context.Context, uac *UAC) {
	if f, ok := s.uacEventHandlers[event]; ok {
		f(ctx, uac)
	}
}

func (s *Scenario) next(ctx context.Context, uas *UAS, programmID string) error {
	if programm, ok := s.programms[programmID]; !ok {
		return ErrUnknownProgramm
	} else {
		s.programm = programm
		s.programm.init(ctx, s.meeting, uas)
		return nil
	}
}

func (s *Scenario) run(ctx context.Context, uas *UAS) error {
	return s.next(ctx, uas, s.rootId)
}

var ErrEmptyScenario = errors.New("unknown scenario")

func NewScenario(ctx context.Context, s *Server, config *ScenarioConfig) (*Scenario, error) {
	if config == nil {
		return nil, ErrEmptyScenario
	} else {
		scenario := &Scenario{
			programms: make(map[string]Programm),
		}

		for programmID, programmType := range config.Programms {
			key := fmt.Sprintf("/scenario/%s/programm/%s", config.ID, programmID)
			switch programmType {
			case CallProgrammType:
				programmConfig := &CallProgrammConfig{}
				if err := s.db.Get(ctx, key, programmConfig); err != nil {
					return nil, err
				} else {
					scenario.programms[programmID] = &CallProgramm{
						config: programmConfig,
					}
				}
				// case RejectProgrammType:
				// 	config := &RejectProgrammConfig{}
				// 	if err := s.db.Get(ctx, key, config); err != nil {
				// 		return nil, err
				// 	} else {
				// 		scenario.programms[programmID] = &RejectProgramm{
				// 			config: config,
				// 		}
				// 	}
			}
		}

		scenario.programm = scenario.programms[config.RootID]
		return scenario, nil
	}
}
