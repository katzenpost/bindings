// session.go - mixnet session client
// Copyright (C) 2017  Yawning Angel, Ruben Pollan, David Stainton
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

package client

import (
	"encoding/hex"
	"fmt"

	"github.com/katzenpost/bindings/python/internal"
	"github.com/katzenpost/client"
	"github.com/op/go-logging"
)

// Session holds the client session
type Session struct {
	client     *client.Client
	log        *logging.Logger
	clientCfg  *client.Config
	sessionCfg *client.SessionConfig
	session    *client.Session
}

// NewSession stablishes a session with provider using key
func (c Client) NewSession(user string, provider string, key Key) (Session, error) {
	var err error
	var session Session
	clientCfg := &client.Config{
		User:       user,
		Provider:   provider,
		LinkKey:    key.priv,
		LogBackend: c.log,
		PKIClient:  c.pki,
	}
	gClient, err := client.New(clientCfg)
	if err != nil {
		return session, err
	}
	session.client = gClient
	session.log = c.log.GetLogger(fmt.Sprintf("session_%s@%s", user, provider))
	return session, err
}

// Connect connects the client to the Provider
func (s Session) Connect(identityKey Key) error {
	consumer := internal.NewMessageConsumer(s.log)
	userKeyDiscoveryStub := internal.UserKeyDiscoveryStub{}
	sessionCfg := client.SessionConfig{
		User:             s.clientCfg.User,
		Provider:         s.clientCfg.Provider,
		IdentityPrivKey:  identityKey.priv,
		LinkPrivKey:      s.clientCfg.LinkKey,
		MessageConsumer:  consumer,
		Storage:          new(internal.StorageStub),
		UserKeyDiscovery: userKeyDiscoveryStub,
	}
	s.sessionCfg = &sessionCfg
	var err error
	s.session, err = s.client.NewSession(&sessionCfg)
	return err
}

func (s Session) GetMessage() string {
	return s.sessionCfg.MessageConsumer.(internal.MessageConsumer).GetMessage()
}

// Shutdown the session
func (s Session) Shutdown() {
	s.Shutdown()
}

// Send into the mix network
func (s Session) Send(recipient, provider, msg string) error {
	raw, err := hex.DecodeString(msg)
	if err != nil {
		return err
	}
	messageID, err := s.session.Send(recipient, provider, raw)
	if err != nil {
		return err
	}
	s.log.Debugf("sent message with messageID %x", messageID)
	return nil
}

// SendUnreliable into the mix network
func (s Session) SendUnreliable(recipient, provider, msg string) error {
	raw, err := hex.DecodeString(msg)
	if err != nil {
		return err
	}
	return s.session.SendUnreliable(recipient, provider, raw)
}
