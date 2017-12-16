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

package fenced

import (
	"errors"

	"github.com/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/minclient/block"
	"github.com/op/go-logging"
)

// StorageStub implements the Storage interface
// as defined in the client library.
// XXX This should be replaced by something useful.
type StorageStub struct {
}

// GetBlocks returns a slice of blocks
func (s StorageStub) GetBlocks(*[block.MessageIDLength]byte) ([][]byte, error) {
	return nil, errors.New("failure: StorageStub GetBlocks not yet implemented")
}

// PutBlock puts a block into storage
func (s StorageStub) PutBlock(*[block.MessageIDLength]byte, []byte) error {
	return errors.New("failure: StorageStub PutBlock not yet implemented")
}

type MessageConsumer struct {
	log             *logging.Logger
	ingressMsgQueue chan string
}

// ReceivedMessage is used to receive a message.
// This is a method on the MessageConsumer interface
// which is defined in the client library.
// XXX fix me
func (c MessageConsumer) ReceivedMessage(senderPubKey *ecdh.PublicKey, message []byte) {
	c.log.Debug("ReceivedMessage")
	c.ingressMsgQueue <- string(message)
}

// GetMessage blocks until there is a message in the inbox
func (c MessageConsumer) GetMessage() string {
	c.log.Debug("GetMessage")
	return <-c.ingressMsgQueue
}

// ReceivedACK is used to receive a signal that a message was received by
// the recipient Provider. This is a method on the MessageConsumer interface
// which is defined in the client library.
// XXX fix me
func (c MessageConsumer) ReceivedACK(messageID *[block.MessageIDLength]byte, message []byte) {
	c.log.Debug("ReceivedACK")
}

func NewMessageConsumer(log *logging.Logger) MessageConsumer {
	c := MessageConsumer{
		log:             log,
		ingressMsgQueue: make(chan string, 100),
	}
	return c
}

type UserKeyDiscoveryStub struct{}

// Get returns the identity public key for a given identity.
// This is part of the UserKeyDiscovery interface defined
// in the client library.
// XXX fix me
func (u UserKeyDiscoveryStub) Get(identity string) (*ecdh.PublicKey, error) {
	//u.log.Debugf("Get identity %s", identity)
	return nil, nil
}
