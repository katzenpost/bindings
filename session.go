package minclient

import (
	"encoding/hex"

	"github.com/katzenpost/core/crypto/rand"
	"github.com/katzenpost/core/log"
	"github.com/katzenpost/core/sphinx"
	"github.com/katzenpost/core/sphinx/constants"
	"github.com/katzenpost/core/utils"
	"github.com/katzenpost/minclient"
	"github.com/katzenpost/minclient/block"
)

// Session holds the client session
type Session struct {
	client *minclient.Client
	queue  chan string
	log    *log.Backend

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
		OnMessageFn: session.onMessage,
		OnACKFn:     session.onACK,
	}

	session.queue = make(chan string, 100)
	session.surbKeys = make(map[[constants.SURBIDLength]byte][]byte)
	session.client, err = minclient.New(clientCfg)
	session.log = client.log
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
func (s *Session) GetMessage() string {
	return <-s.queue
}

func (s *Session) onMessage(b []byte) error {
	lm := s.log.GetLogger("callbacks:onMessage")
	lm.Noticef("Received Message: %v", len(b))
	lm.Noticef("====> %v", string(b))

	s.queue <- string(b)
	return nil
}

func (s *Session) onACK(id *[constants.SURBIDLength]byte, b []byte) error {
	lm := s.log.GetLogger("callbacks:onACK")
	lm.Noticef("Received SURB-ACK: %v", len(b))
	lm.Noticef("SURB-ID: %v", hex.EncodeToString(id[:]))

	// surbKeys should have a lock in production code, but lazy.
	k, ok := s.surbKeys[*id]
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
}

func (s *Session) onConn(isConnected bool) {
	lm := s.log.GetLogger("callbacks:onConn")
	lm.Noticef("Peer connection status changed: %v", isConnected)
}
