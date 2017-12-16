// key.go - mixnet user key
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

package client

import (
	"github.com/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/core/crypto/rand"
)

// Key keeps the key public and private data
type Key struct {
	priv *ecdh.PrivateKey
}

// GenKey creates a new ecdh key
func GenKey() (*Key, error) {
	mKey := new(Key)
	var err error
	mKey.priv, err = ecdh.NewKeypair(rand.Reader)
	if err != nil {
		return mKey, err
	}
	return mKey, err
}
