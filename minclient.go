package minclient

import (
	"github.com/katzenpost/core/crypto/eddsa"
	cpki"github.com/katzenpost/core/pki"
	"github.com/katzenpost/core/log"
	npki "github.com/katzenpost/authority/nonvoting/client"
)

// TODO: let's make this configurable
var clientLog *log.Backend
var pkiClient cpki.Client

// SetUpPKI configures the pki to be used
func SetUpPKI(address string, key string) error {
	var pubKey eddsa.PublicKey
	err := pubKey.FromString(key)
	if err != nil {
		return err
	}

	clientLog, err = log.New("/tmp/katzenpost.log", "DEBUG", false)
	if err != nil {
		return err
	}

	pkiCfg := npki.Config{
		LogBackend: clientLog,
		Address: address,
		PublicKey: &pubKey,
	}
	pkiClient, err = npki.New(&pkiCfg)
	return err
}
