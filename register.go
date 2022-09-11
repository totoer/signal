package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"signal/sip"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

var HASH_FUNCTION = "MD5"

type RegistrationType string

const (
	AuthRegistration        RegistrationType = "AUTH_REGISTRATION"
	NonAuthRegistration     RegistrationType = "NON_AUTH_REGISTRATION"
	TransparentRegistration RegistrationType = "TRANSPARENT_REGISTRATION"
)

type Account struct {
	RegistrationType RegistrationType `json:"registration_type"`
	Login            string           `json:"login"`
	Password         string           `json:"password"`
	Incoming         *ScenarioConfig  `json:"incoming"`
	Outgoing         *ScenarioConfig  `json:"outgoing"`
}

type Registration struct {
	ID              uuid.UUID                       `json:"id"`
	Destination     sip.Destination                 `json:"destination"`
	Host            string                          `json:"host"`
	Login           string                          `json:"login"`
	Authorized      bool                            `json:"authorized"`
	Contacts        []sip.Contact                   `json:"contacts"`
	SourceAddres    net.Addr                        `json:"source_addres"`
	WWWAuthenticate map[string]*sip.WWWAuthenticate `json:"www-authenticate"`
	Expires         int
	Account         *Account `json:"account"`
}

func NewRegistration(acc *Account, contacts []sip.Contact, addr net.Addr, destination sip.Destination, host, login string, authorized bool) *Registration {
	return &Registration{
		ID:              uuid.New(),
		Destination:     destination,
		Host:            host,
		Login:           login,
		Authorized:      authorized,
		Contacts:        contacts,
		SourceAddres:    addr,
		WWWAuthenticate: make(map[string]*sip.WWWAuthenticate),
		Expires:         60 * 60 * 24,
		Account:         acc,
	}
}

type Register struct {
	server  *Server
	pool    map[string]*Registration
	callMap map[string]string
}

var ErrRegistrationNotExists = errors.New("registration not exists")

// var ErrAccRegistrationNotTransparant = errors.New("account registration not transparant")

func (r *Register) loadRegistration(ctx context.Context, host, login string) (*Registration, error) {
	key := fmt.Sprintf("/register/%s/%s", host, login)

	if registration, ok := r.pool[key]; ok {
		return registration, nil
	} else if err := r.server.db.Get(ctx, key, registration); err == nil {
		r.pool[key] = registration
		return registration, nil
	}

	return nil, ErrRegistrationNotExists
}

func (r *Register) storeRegistration(ctx context.Context, host, login string, registration *Registration) error {
	key := fmt.Sprintf("/register/%s/%s", host, login)
	if err := r.server.db.Put(ctx, key, registration); err != nil {
		return err
	} else {
		r.pool[key] = registration
		return nil
	}
}

func (r *Register) loadRegistrationByDestination(ctx context.Context, dest sip.Destination) (*Registration, error) {
	return nil, nil
}

func (r *Register) registration(ctx context.Context, cid string, acc *Account, registration *Registration, req *sip.Request) error {
	if from, err := req.GetHeaders().GetFrom(); err != nil {
		log.Error().Err(err).Str("Call-ID", cid).
			Str("where", "Register.registration").
			Msg("While registration")
		return err
	} else if host, login, err := req.GetHeaders().GetHostLoginByFrom(); err != nil {
		log.Error().Err(err).Str("Call-ID", cid).
			Str("where", "Register.registration").
			Msg("While registration")
		return err
	} else if resp, err := req.MakeResponse(sip.Unauthorized); err != nil {
		log.Error().Err(err).Str("Call-ID", cid).
			Msg("While registration")
		return err
	} else {
		if registration == nil {
			contacts, _ := req.GetHeaders().GetContacts()
			registration = NewRegistration(acc, contacts, req.GetSourceAddres(), from, host, login, false)
			log.Info().Str("Call-ID", cid).
				Str("where", "Register.registration").
				Str("login", login).
				Str("host", host).
				Str("registration_id", registration.ID.String()).
				Msg("Create new registration")
		}

		registration.WWWAuthenticate[cid] = &sip.WWWAuthenticate{
			Realm:     host,
			Nonce:     uuid.New().String(),
			Algorithm: HASH_FUNCTION,
		}

		if err := r.storeRegistration(ctx, host, login, registration); err != nil {
			log.Error().Err(err).Str("Call-ID", cid).
				Str("where", "Register.registration").
				Str("login", login).
				Str("host", host).
				Str("registration_id", registration.ID.String()).
				Msg("While store registration")
		}

		resp.Headers.To.Tag = uuid.New().String()
		resp.Headers.WWWAuthenticate = registration.WWWAuthenticate[cid]

		if err := r.server.transport.SendSIP(resp); err != nil {
			log.Error().Err(err).Str("Call-ID", cid).
				Str("where", "Register.registration").
				Str("login", login).
				Str("host", host).
				Msg("While send response")
		}

		return nil
	}
}

var ErrUnsupportedRegistration = errors.New("unsuported registration")
var ErrRegistrationNotPrepared = errors.New("registration not prepared")

func (r *Register) auth(ctx context.Context, cid string, req *sip.Request, handle func(context.Context, *Registration) error) error {
	if from, err := req.GetHeaders().GetFrom(); err != nil {
		log.Error().Err(err).Str("Call-ID", cid).
			Str("where", "Register.auth").
			Msg("While auth")
		return err
	} else if host, login, err := req.GetHeaders().GetHostLoginByFrom(); err != nil {
		log.Error().Err(err).Str("Call-ID", cid).
			Str("where", "Register.auth").
			Msg("While auth")
		return err
	} else {
		if registration, err := r.loadRegistration(ctx, host, login); err == nil && registration.Authorized {
			log.Info().Str("Call-ID", cid).
				Str("where", "Register.auth").
				Str("host", host).
				Str("login", login).
				Msg("Registration exists and authorized, continues message handle")
			return handle(ctx, registration)
		} else {
			log.Info().Str("Call-ID", cid).
				Str("where", "Register.auth").
				Str("login", login).
				Str("host", host).
				Msg("Auth check")
			accountKey := fmt.Sprintf("/account/%s/%s", host, login)
			acc := &Account{}

			if err := r.server.db.Get(ctx, accountKey, acc); err != nil {
				log.Error().Err(err).Str("Call-ID", cid).
					Str("where", "Register.auth").
					Str("login", login).
					Str("host", host).
					Msg("Account not exists")
				return err
			} else if acc.RegistrationType == NonAuthRegistration {
				log.Info().Str("Call-ID", cid).
					Str("where", "Register.auth").
					Str("login", login).
					Str("host", host).
					Msg("Non auth registration")
				if registration == nil {
					contacts, _ := req.GetHeaders().GetContacts()
					registration = NewRegistration(acc, contacts, req.GetSourceAddres(), from, host, login, true)
					log.Info().Str("Call-ID", cid).
						Str("where", "Register.auth").
						Str("login", login).
						Str("host", host).
						Str("registration_id", registration.ID.String()).
						Msg("Create new registration")
					if err := r.storeRegistration(ctx, host, login, registration); err != nil {
						log.Error().Err(err).Str("Call-ID", cid).
							Str("where", "Register.auth").
							Str("login", login).
							Str("host", host).
							Str("registration_id", registration.ID.String()).
							Msg("While store registration")
						return err
					}
				}
				log.Info().Str("Call-ID", cid).
					Str("where", "Register.auth").
					Str("host", host).
					Str("login", login).
					Msg("Registration exists and authorized, continues message handle")
				return handle(ctx, registration)
			} else if registration == nil {
				return r.registration(ctx, cid, acc, nil, req)
			} else if authorization, err := req.GetHeaders().GetAuthorization(); err != nil {
				log.Info().Str("Call-ID", cid).
					Str("where", "Register.auth").
					Str("login", login).
					Str("host", host).
					Msg("Message does not contain Authorization")
				return r.registration(ctx, cid, acc, registration, req)
			} else if wwwAuthenticate, ok := registration.WWWAuthenticate[cid]; !ok {
				log.Info().Str("Call-ID", cid).
					Str("where", "Register.auth").
					Str("login", login).
					Str("host", host).
					Msg("Not found WWWAuthenticate for Call-ID")
				return r.registration(ctx, cid, acc, registration, req)
			} else {
				ha1s := fmt.Sprintf("%s:%s:%s", login, wwwAuthenticate.Realm, acc.Password)
				ha2s := fmt.Sprintf("%s:%s", req.Method, req.URI.String())
				ha1 := md5.Sum([]byte(ha1s))
				ha2 := md5.Sum([]byte(ha2s))
				rs := fmt.Sprintf("%s:%s:%s", hex.EncodeToString(ha1[:]), wwwAuthenticate.Nonce, hex.EncodeToString(ha2[:]))
				response := md5.Sum([]byte(rs))

				if authorization.Response != hex.EncodeToString(response[:]) {
					log.Info().Str("Call-ID", cid).
						Str("where", "Register.auth").
						Str("login", login).
						Str("host", host).
						Msg("Authorization response not expected")
					return r.registration(ctx, cid, acc, registration, req)
				} else {
					registration.Authorized = true
					if err := r.storeRegistration(ctx, host, login, registration); err != nil {
						log.Error().Err(err).Str("Call-ID", cid).
							Str("where", "Register.auth").
							Str("login", login).
							Str("host", host).
							Str("registration_id", registration.ID.String()).
							Msg("While store registration")
					}
					log.Info().Str("Call-ID", cid).
						Str("where", "Register.auth").
						Str("host", host).
						Str("login", login).
						Msg("Registration exists and authorized, continues message handle")
					return handle(ctx, registration)
				}
			}
		}
	}
}

func (r *Register) bind(cid string, host, login string) {
	key := fmt.Sprintf("/register/%s/%s", host, login)
	r.callMap[cid] = key
}

func (r *Register) loadRegistrationByCallID(ctx context.Context, cid string) (*Registration, error) {
	if key, ok := r.callMap[cid]; !ok {
		return nil, ErrRegistrationNotExists
	} else {
		if registration, ok := r.pool[key]; ok {
			return registration, nil
		} else if err := r.server.db.Get(ctx, key, registration); err == nil {
			r.pool[key] = registration
			return registration, nil
		}

		return nil, ErrRegistrationNotExists
	}
}

func NewRegister(s *Server) *Register {
	return &Register{
		server: s,
		pool:   make(map[string]*Registration),
	}
}
