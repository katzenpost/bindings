package client

import (
	"github.com/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/mailproxy"
	"github.com/katzenpost/mailproxy/config"
	"github.com/katzenpost/mailproxy/event"
)

const (
	pkiName = "default"
)

var identityKeyBytes = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

// Client is katzenpost object
type Client struct {
	address   string
	proxy     *mailproxy.Proxy
	eventSink chan event.Event
}

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
		Recipients: map[string]string{},
	}
	err = proxyCfg.FixupAndValidate()
	if err != nil {
		return Client{}, err
	}

	proxy, err := mailproxy.New(&proxyCfg)
	return Client{cfg.getAddress(), proxy, eventSink}, err
}

func (c Client) Shutdown() {
	c.proxy.Shutdown()
	c.proxy.Wait()
}

func (c Client) Send(recipient, msg string) error {
	var identityKey ecdh.PrivateKey
	identityKey.FromBytes(identityKeyBytes)
	c.proxy.SetRecipient(recipient, identityKey.PublicKey())
	return c.proxy.SendMessage(c.address, recipient, []byte(msg))
}
