// client.go - mixnet client
// Copyright (C) 2017  Yawning Angel.
// Copyright (C) 2018  Ruben Pollan.
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

package katzenpost

import (
	"encoding/hex"
	"errors"
	"time"

	"github.com/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/mailproxy"
	"github.com/katzenpost/mailproxy/config"
	"github.com/katzenpost/mailproxy/event"
	"github.com/katzenpost/minclient/block"
)

const (
	pkiName        = "default"
	kaetzenTimeout = 2 * time.Minute
)

// TimeoutError is returned on timeouts
type TimeoutError struct{}

func (t TimeoutError) Error() string {
	return "Timeout"
}

// Client is katzenpost object
type Client struct {
	address      string
	proxy        *mailproxy.Proxy
	eventSink    chan event.Event
	recvCh       chan bool
	connectionCh chan bool
	kaetzchenCh  chan kaetzchenKeyRequest
}

type kaetzchenKeyRequest struct {
	recipient string
	msgID     []byte
	ch        chan error
}

// New creates a katzenpost client
func New(cfg *Config) (*Client, error) {
	eventSink := make(chan event.Event)
	dataDir, err := cfg.getDataDir()
	if err != nil {
		return &Client{}, err
	}

	proxyCfg := config.Config{
		Proxy: &config.Proxy{
			NoLaunchListeners: true,
			DataDir:           dataDir,
			EventSink:         eventSink,
		},
		Logging: cfg.getLogging(),
		UpstreamProxy: &config.UpstreamProxy{
			Type: "none",
		},

		NonvotingAuthority: map[string]*config.NonvotingAuthority{
			pkiName: cfg.getAuthority(),
		},
		Account:    []*config.Account{cfg.getAccount()},
		Recipients: map[string]*ecdh.PublicKey{},
	}
	err = proxyCfg.FixupAndValidate()
	if err != nil {
		return &Client{}, err
	}

	recvCh := make(chan bool, 10)
	connectionCh := make(chan bool, 10)
	kaetzchenCh := make(chan kaetzchenKeyRequest, 10)
	proxy, err := mailproxy.New(&proxyCfg)
	c := &Client{cfg.getAddress(), proxy, eventSink, recvCh, connectionCh, kaetzchenCh}
	go c.eventHandler()
	return c, err
}

// WaitToConnect wait's to be connected
func (c Client) WaitToConnect() error {
	isConnected := <-c.connectionCh
	if !isConnected {
		return errors.New("Not connected")
	}
	return nil
}

// Shutdown the client
func (c Client) Shutdown() {
	c.proxy.Shutdown()
}

// Send a message into katzenpost
func (c Client) Send(recipient, msg string) ([]byte, error) {
	return c.proxy.SendMessage(c.address, recipient, []byte(msg))
}

// Message received from katzenpost
type Message struct {
	Sender    string
	Payload   string
	SenderKey string
	MessageID []byte
}

// GetMessage from katzenpost
func (c Client) GetMessage(timeoutMs int64) (*Message, error) {
	if timeoutMs == 0 {
		<-c.recvCh
		return c.getMsg()
	}

	select {
	case <-c.recvCh:
		return c.getMsg()
	case <-time.After(time.Millisecond * time.Duration(timeoutMs)):
		return nil, nil
	}
}

func (c Client) getMsg() (*Message, error) {
	msg, err := c.proxy.ReceivePop(c.address)
	return &Message{msg.SenderID, string(msg.Payload), hex.EncodeToString(msg.SenderKey.Bytes()), msg.MessageID}, err
}

// GetKey for recipient and add it to the local key storage
func (c Client) GetKey(recipient string) error {
	msgID, err := c.proxy.QueryKeyFromProvider(c.address, recipient)
	if err != nil {
		return err
	}
	ch := make(chan error)
	c.kaetzchenCh <- kaetzchenKeyRequest{recipient, msgID, ch}

	select {
	case err = <-ch:
		return err
	case <-time.After(kaetzenTimeout):
		return TimeoutError{}
	}
}

// HasKey returns if the key storage have a key for recipient
func (c Client) HasKey(recipient string) bool {
	key, err := c.proxy.GetRecipient(recipient)
	return err == nil && key != nil
}

type requestIndex [block.MessageIDLength]byte

func (c Client) eventHandler() {
	keyRequests := make(map[requestIndex]kaetzchenKeyRequest)

	for {
		select {
		case ev := <-c.eventSink:
			switch event := ev.(type) {
			case *event.MessageReceivedEvent:
				c.recvCh <- true
			case *event.ConnectionStatusEvent:
				c.connectionCh <- event.IsConnected
			case *event.KaetzchenReplyEvent:
				var index requestIndex
				copy(index[:], event.MessageID[:block.MessageIDLength])
				request, ok := keyRequests[index]
				if !ok {
					continue
				}

				if event.Err != nil {
					request.ch <- event.Err
					continue
				}

				_, key, err := c.proxy.ParseKeyQueryResponse(event.Payload)
				if err != nil {
					request.ch <- err
				} else {
					request.ch <- c.proxy.SetRecipient(request.recipient, key)
				}
				delete(keyRequests, index)
			default:
				continue
			}

		case request := <-c.kaetzchenCh:
			var index requestIndex
			copy(index[:], request.msgID[:block.MessageIDLength])
			keyRequests[index] = request
		}
	}
}
