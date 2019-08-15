package keys

import (
	"crypto/rand"
	"path/filepath"

	p2pcrypto "github.com/libp2p/go-libp2p-core/crypto"
)

func GetPrivateKeyFromPath(path string) (p2pcrypto.PrivKey, error) {
	keyBytes, err := readFile(path)
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
	if err := mkdirAll(dir); err != nil {
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
	if err := writeFile(path, []byte(encodedKey)); err != nil {
		return nil, err
	}
	return privKey, nil
}
