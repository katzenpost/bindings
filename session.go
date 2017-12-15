// session.go - mixnet session client
// Copyright (C) 2017  Yawning Angel.
// Copyright (C) 2017  Ruben Pollan.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package minclient

import (
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/katzenpost/core/crypto/rand"
	"github.com/katzenpost/core/sphinx"
	"github.com/katzenpost/core/sphinx/constants"
	"github.com/katzenpost/core/utils"
	"github.com/katzenpost/minclient"
	"github.com/katzenpost/minclient/block"
	"github.com/op/go-logging"
)

// Session holds the client session
type Session struct {
	client *minclient.Client
	queue  chan string
	log    *logging.Logger

	// TODO: we'll need to add persistency to the surb keys at some point
	surbKeys map[[constants.SURBIDLength]byte][]byte
}

// NewSession stablishes a session with provider using key
func (client Client) NewSession(user string, provider string, key Key) (Session, error) {
	var err error
	var session Session

	clientCfg := &minclient.ClientConfig{
		User:        user,
		Provider:    provider,
		LinkKey:     key.priv,
		LogBackend:  client.log,
		PKIClient:   client.pki,
		OnConnFn:    session.onConn,
		OnEmptyFn:   session.onEmpty,
		OnMessageFn: session.onMessage,
		OnACKFn:     session.onACK,
	}

	session.queue = make(chan string, 100)
	session.surbKeys = make(map[[constants.SURBIDLength]byte][]byte)
	session.client, err = minclient.New(clientCfg)
	session.log = client.log.GetLogger(fmt.Sprintf("callbacks:%s@%s", user, provider))
	return session, err
}

// Shutdown the session
func (s Session) Shutdown() {
	s.client.Shutdown()
}

// SendMessage into the mix network
func (s Session) SendMessage(recipient, provider, msg string) error {
	surbID := [constants.SURBIDLength]byte{}
	_, err := rand.Reader.Read(surbID[:])
	if err != nil {
		return err
	}

	chunk := [block.BlockCiphertextLength]byte{}
	copy(chunk[:], []byte(msg))
	surbKey, _, err := s.client.SendCiphertext(recipient, provider, &surbID, chunk[:])
	if err != nil {
		return err
	}

	s.surbKeys[surbID] = surbKey
	return nil
}

// GetMessage blocks until there is a message in the inbox
func (s *Session) GetMessage(timeout int) (string, error) {
	if timeout == 0 {
		return <-s.queue, nil
	}

	select {
	case msg := <-s.queue:
		return msg, nil
	case <-time.After(time.Second * timeout):
		return "", errors.New("Timeout")
	}
}

func (s *Session) onMessage(b []byte) error {
	s.log.Noticef("Received Message: %v", len(b))
	s.log.Noticef("====> %v", string(b))

	s.queue <- string(b)
	return nil
}

func (s *Session) onACK(id *[constants.SURBIDLength]byte, b []byte) error {
	s.log.Noticef("Received SURB-ACK: %v", len(b))
	s.log.Noticef("SURB-ID: %v", hex.EncodeToString(id[:]))

	// surbKeys should have a lock in production code, but lazy.
	k, ok := s.surbKeys[*id]
	if !ok {
		s.log.Errorf("Failed to find SURB SPRP key")
		return nil
	}

	payload, err := sphinx.DecryptSURBPayload(b, k)
	if err != nil {
		s.log.Errorf("Failed to decrypt SURB: %v", err)
		return nil
	}
	if utils.CtIsZero(payload) {
		s.log.Noticef("SURB Payload: %v bytes of 0x00", len(payload))
	} else {
		s.log.Noticef("SURB Payload: %v", hex.Dump(payload))
	}

	return nil
}

func (s *Session) onConn(isConnected bool) {
	s.log.Noticef("Peer connection status changed: %v", isConnected)
}

func (s *Session) onEmpty() error {
	s.log.Noticef("Providers queue is empty")
	return nil
}
