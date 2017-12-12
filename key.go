package minclient

import (
	"encoding/hex"

	"github.com/katzenpost/core/crypto/ecdh"
	"github.com/katzenpost/core/crypto/rand"
)

// Key keeps the key public and private data
type Key struct {
	Private string
	Public  string
	priv    *ecdh.PrivateKey
}

// GenKey creates a new ecdh key
func GenKey() (Key, error) {
	key, err := ecdh.NewKeypair(rand.Reader)
	if err != nil {
		return Key{}, err
	}
	return buildKey(key), nil
}

// StringToKey builds a Key from a string
func StringToKey(keyStr string) (Key, error) {
	var key ecdh.PrivateKey

	keyBytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return Key{}, err
	}

	err = key.FromBytes(keyBytes)
	if err != nil {
		return Key{}, err
	}

	return buildKey(&key), nil
}

func buildKey(key *ecdh.PrivateKey) Key {
	return Key{
		Private: hex.EncodeToString(key.Bytes()),
		Public:  key.PublicKey().String(),
		priv:    key,
	}
}
