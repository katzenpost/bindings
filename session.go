package minclient

import (
	"encoding/hex"
	"errors"

	"github.com/katzenpost/core/crypto/rand"
	"github.com/katzenpost/core/sphinx"
	"github.com/katzenpost/core/sphinx/constants"
	"github.com/katzenpost/core/utils"
	"github.com/katzenpost/minclient"
	"github.com/katzenpost/minclient/block"
)

// TODO: we might need that being long lived
var surbKeys = make(map[[constants.SURBIDLength]byte][]byte)

// Session holds the client session
type Session struct {
	client *minclient.Client
}

// NewSession stablishes a session with provider using key
func NewSession(user string, provider string, key Key) (Session, error) {
	if pkiClient == nil {
		return Session{}, errors.New("PKI is not configured")
	}
	lm := clientLog.GetLogger("callbacks:main")

	clientCfg := &minclient.ClientConfig{
		User:        user,
		Provider:    provider,
		LinkKey:     key.priv,
		LogBackend:  clientLog,
		PKIClient:   pkiClient,
		OnConnFn:    func(isConnected bool) {
			lm.Noticef("Peer connection status changed: %v", isConnected)
		},
		OnMessageFn: func(b []byte) error {
			// TODO: we need to handle incomming messages
			lm.Noticef("Received Message: %v", len(b))
			lm.Noticef("====> %v", string(b))

			return nil
		},
		OnACKFn: func(id *[constants.SURBIDLength]byte, b []byte) error {
			lm.Noticef("Received SURB-ACK: %v", len(b))
			lm.Noticef("SURB-ID: %v", hex.EncodeToString(id[:]))

			// surbKeys should have a lock in production code, but lazy.
			k, ok := surbKeys[*id]
			if !ok {
				lm.Errorf("Failed to find SURB SPRP key")
				return nil
			}

			payload, err := sphinx.DecryptSURBPayload(b, k)
			if err != nil {
				lm.Errorf("Failed to decrypt SURB: %v", err)
				return nil
			}
			if utils.CtIsZero(payload) {
				lm.Noticef("SURB Payload: %v bytes of 0x00", len(payload))
			} else {
				lm.Noticef("SURB Payload: %v", hex.Dump(payload))
			}

			return nil
		},
	}

	var session Session
	var err error
	session.client, err = minclient.New(clientCfg)
	return session, err
}

// Shutdown the session
func (s Session) Shutdown() {
    s.client.Shutdown()
}

// SendMessage into the mix network
func (s Session) SendMessage(recipient, provider, msg string) error {
	// TODO: deal with ACKs
	surbID := [constants.SURBIDLength]byte{}
	_, err := rand.Reader.Read(surbID[:])
	if err != nil {
		return err
	}

	chunk := [block.BlockCiphertextLength]byte{}
	copy(chunk[:], []byte(msg))
	_, _, err = s.client.SendCiphertext(recipient, provider, &surbID, chunk[:])
	return err
}
