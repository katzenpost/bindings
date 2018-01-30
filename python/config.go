// config.go - mixnet client configuration
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
	"fmt"
	"os"
	"path"

	"github.com/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/core/crypto/eddsa"
	"github.com/katzenpost/mailproxy/config"
)

// Config has the client configuration
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
		StorageKey:  nil,
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
