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
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/core/pki"
	"github.com/katzenpost/mailproxy"
	"github.com/katzenpost/mailproxy/config"
	"github.com/katzenpost/mailproxy/event"
)

const (
	pkiName = "default"
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
}

// New creates a katzenpost client
func New(cfg Config) (Client, error) {
	eventSink := make(chan event.Event)
	dataDir, err := cfg.getDataDir()
	if err != nil {
		return Client{}, err
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
		return Client{}, err
	}

	recvCh := make(chan bool, 10)
	connectionCh := make(chan bool, 10)
	proxy, err := mailproxy.New(&proxyCfg)
	c := Client{cfg.getAddress(), proxy, eventSink, recvCh, connectionCh}
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

// ListProviders returns the provider list
func (c Client) ListProviders() ([]string, error) {
	providers, err := c.proxy.ListProviders(pkiName)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(providers))
	for i, provider := range providers {
		names[i] = provider.Name
	}
	return names, nil
}

// Shutdown the client
func (c Client) Shutdown() {
	c.proxy.Shutdown()
}

// Send a message into katzenpost
func (c Client) Send(recipient, msg string) error {
	err := c.fetchKey(recipient)
	if err != nil {
		return err
	}
	return c.proxy.SendMessage(c.address, recipient, []byte(msg))
}

func (c Client) fetchKey(address string) error {
	parts := strings.Split(address, "@")
	if len(parts) != 2 {
		return errors.New("Not valid address address: " + address)
	}
	user := strings.ToLower(parts[0])
	providerName := parts[1]

	providers, err := c.proxy.ListProviders(pkiName)
	if err != nil {
		return err
	}
	providerAddress := ""
	for _, provider := range providers {
		if provider.Name == providerName {
			addr := provider.Addresses[pki.TransportTCPv4][0]
			providerAddress = strings.Split(addr, ":")[0]
			break
		}
	}
	if providerAddress == "" {
		return errors.New("Recipient provider doesn't exist in the authority document: " + providerName)
	}

	resp, err := http.PostForm("http://"+providerAddress+":7900/getidkey", url.Values{"user": {user}})
	if err != nil {
		return errors.New("Can't fetch key for address: " + err.Error())
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var response struct {
		Getidkey string
	}
	err = decoder.Decode(&response)
	if err != nil {
		return errors.New("There was a problem reading the key fetch response: " + err.Error())
	}

	var key ecdh.PublicKey
	key.FromString(response.Getidkey)
	c.proxy.SetRecipient(address, &key)
	return nil
}

// Message received from katzenpost
type Message struct {
	Sender  string
	Payload string
}

// GetMessage from katzenpost
func (c Client) GetMessage(timeout int64) (Message, error) {
	if timeout == 0 {
		<-c.recvCh
		return c.getMsg()
	}

	select {
	case <-c.recvCh:
		return c.getMsg()
	case <-time.After(time.Millisecond * time.Duration(timeout)):
		return Message{}, TimeoutError{}
	}
}

func (c Client) getMsg() (Message, error) {
	msg, err := c.proxy.ReceivePop(c.address)
	return Message{msg.SenderID, string(msg.Payload)}, err
}

func (c Client) eventHandler() {
	for {
		ev := <-c.eventSink
		switch ev.(type) {
		case *event.MessageReceivedEvent:
			c.recvCh <- true
		case *event.ConnectionStatusEvent:
			conEv := ev.(*event.ConnectionStatusEvent)
			c.connectionCh <- conEv.IsConnected
		default:
			continue
		}
	}
}
