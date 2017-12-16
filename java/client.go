// minclient.go - mixnet client
// Copyright (C) 2017  Yawning Angel.
// Copyright (C) 2017  Ruben Pollan.
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

// Package client provides a mixnet client library
package client

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
func NewClient(pkiAddress, pkiKey string, logConfig *LogConfig) (*Client, error) {
	var pubKey eddsa.PublicKey
	err := pubKey.FromString(pkiKey)
	if err != nil {
		return nil, err
	}

	logLevel := "NOTICE"
	if logConfig.Level != "" {
		logLevel = logConfig.Level
	}
	client := new(Client)
	client.log, err = log.New(logConfig.File, logLevel, !logConfig.Enabled)
	if err != nil {
		return nil, err
	}

	pkiCfg := npki.Config{
		LogBackend: client.log,
		Address:    pkiAddress,
		PublicKey:  &pubKey,
	}
	client.pki, err = npki.New(&pkiCfg)
	return client, err
}
