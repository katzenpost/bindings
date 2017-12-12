package minclient

import (
	npki "github.com/katzenpost/authority/nonvoting/client"
	"github.com/katzenpost/core/crypto/eddsa"
	"github.com/katzenpost/core/log"
	cpki "github.com/katzenpost/core/pki"
)

// Client is katzenpost object
type Client struct {
	log *log.Backend
	pki cpki.Client
}

// LogConfig keeps the configuration of the loger
type LogConfig struct {
	File    string
	Level   string
	Enabled bool
}

// NewClient configures the pki to be used
func NewClient(pkiAddress, pkiKey string, logConfig LogConfig) (Client, error) {
	var client Client

	var pubKey eddsa.PublicKey
	err := pubKey.FromString(pkiKey)
	if err != nil {
		return client, err
	}

	logLevel := "NOTICE"
	if logConfig.Level != "" {
		logLevel = logConfig.Level
	}
	client.log, err = log.New(logConfig.File, logLevel, !logConfig.Enabled)
	if err != nil {
		return client, err
	}

	pkiCfg := npki.Config{
		LogBackend: client.log,
		Address:    pkiAddress,
		PublicKey:  &pubKey,
	}
	client.pki, err = npki.New(&pkiCfg)
	return client, err
}
