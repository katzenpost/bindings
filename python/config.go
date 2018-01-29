package client

import (
	"fmt"
	"os"
	"path"

	"github.com/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/core/crypto/eddsa"
	"github.com/katzenpost/mailproxy/config"
)

type Config struct {
	PkiAddress string
	PkiKey     string
	User       string
	Provider   string
	LinkKey    Key
	Log        LogConfig
	DataDir    string
}

// LogConfig keeps the configuration of the loger
type LogConfig struct {
	File    string
	Level   string
	Enabled bool
}

func (c Config) getAuthority() *config.NonvotingAuthority {
	var pkiPublicKey eddsa.PublicKey
	pkiPublicKey.FromString(c.PkiKey)
	return &config.NonvotingAuthority{
		Address:   c.PkiAddress,
		PublicKey: &pkiPublicKey,
	}
}

func (c Config) getAccount() *config.Account {
	var identityKey ecdh.PrivateKey
	identityKey.FromBytes(identityKeyBytes)
	return &config.Account{
		User:        c.User,
		Provider:    c.Provider,
		Authority:   pkiName,
		LinkKey:     c.LinkKey.priv,
		IdentityKey: &identityKey,
	}
}

func (c Config) getDataDir() (string, error) {
	if c.DataDir != "" {
		return c.DataDir, nil
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return path.Join(workingDir, "data"), nil
}

func (c Config) getLogging() *config.Logging {
	if c.Log.Level != "" {
		return &config.Logging{
			File:    c.Log.File,
			Level:   c.Log.Level,
			Disable: !c.Log.Enabled,
		}
	}
	return nil
}

func (c Config) getAddress() string {
	return fmt.Sprintf("%s@%s", c.User, c.Provider)
}
