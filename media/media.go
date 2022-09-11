package media

import "github.com/google/uuid"

type MediaChanal struct{}

func (mc *MediaChanal) Start() {}

func (mc *MediaChanal) IsStarted() bool {
	return false
}

func (mc *MediaChanal) End() {}

func (mc *MediaChanal) Listen(am *MediaChanal) {}

func (mc *MediaChanal) Connect(am *MediaChanal) {
	mc.Listen(am)
	am.Listen(mc)
}

func (mc *MediaChanal) Play(f ...string) {}

func (ms *MediaChanal) Beeps() {}

func (mc *MediaChanal) Stop() {}

func NewMediaChanal(mid uuid.UUID, cid string) (*MediaChanal, error) {
	return &MediaChanal{}, nil
}

type MediaMixer struct{}

func (m MediaMixer) Join(mc *MediaChanal) {}

func (m *MediaMixer) Play(f string) {}

func (m *MediaMixer) Stop() {}

func NewMediaMixer() (*MediaMixer, error) {
	return &MediaMixer{}, nil
}
