package keys

import (
	"crypto/rand"
	"io/ioutil"
	"os"
	"path/filepath"

	p2pcrypto "github.com/libp2p/go-libp2p-crypto"
)

func GetPrivateKeyFromPath(path string) (p2pcrypto.PrivKey, error) {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decodedKey, err := p2pcrypto.ConfigDecodeKey(string(keyBytes))
	if err != nil {
		return nil, err
	}
	priv, err := p2pcrypto.UnmarshalPrivateKey(decodedKey)
	if err != nil {
		return nil, err
	}
	return priv, nil
}

func GenerateAndSavePrivateKey(path string) (p2pcrypto.PrivKey, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	privKey, _, err := p2pcrypto.GenerateSecp256k1Key(rand.Reader)
	if err != nil {
		return nil, err
	}
	keyBytes, err := p2pcrypto.MarshalPrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	encodedKey := p2pcrypto.ConfigEncodeKey(keyBytes)
	if err := ioutil.WriteFile(path, []byte(encodedKey), os.ModePerm); err != nil {
		return nil, err
	}
	return privKey, nil
}
