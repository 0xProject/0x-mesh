package ethereum

import (
	"crypto/ecdsa"
	"errors"
	"fmt"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

// Signer defines the methods needed to act as a elliptic curve signer
type Signer interface {
	Sign(message []byte, signerAddress common.Address) (*ECSignature, error)
}

// ECSignature contains the parameters of an elliptic curve signature
type ECSignature struct {
	V byte
	R common.Hash
	S common.Hash
}

// EthRPCSigner is a signer that uses a call to Ethereum JSON-RPC method `eth_call`
// to produce a signature
type EthRPCSigner struct {
	rpcClient *rpc.Client
}

// NewEthRPCSigner instantiates a new EthRPCSigner
func NewEthRPCSigner(rpcClient *rpc.Client) Signer {
	return &EthRPCSigner{
		rpcClient: rpcClient,
	}
}

// Sign signs a message via the `eth_sign` Ethereum JSON-RPC call
func (e *EthRPCSigner) Sign(message []byte, signerAddress common.Address) (*ECSignature, error) {
	var signatureHex string
	if err := e.rpcClient.Call(&signatureHex, "eth_sign", signerAddress.Hex(), common.Bytes2Hex(message)); err != nil {
		return nil, err
	}
	// `eth_sign` returns the signature in the [R || S || V] format where V is 0 or 1.
	signatureBytes := common.Hex2Bytes(signatureHex[2:])
	vParam := signatureBytes[64]
	if vParam == byte(0) {
		vParam = byte(27)
	} else if vParam == byte(1) {
		vParam = byte(28)
	}

	ecSignature := &ECSignature{
		V: vParam,
		R: common.BytesToHash(signatureBytes[0:32]),
		S: common.BytesToHash(signatureBytes[32:64]),
	}
	return ecSignature, nil
}

// LocalSigner is a signer that produces an `eth_sign`-compatible signature locally using
// a private key
type LocalSigner struct {
	privateKey *ecdsa.PrivateKey
}

// NewLocalSigner instantiates a new LocalSigner
func NewLocalSigner(privateKey *ecdsa.PrivateKey) Signer {
	return &LocalSigner{
		privateKey: privateKey,
	}
}

// GetSignerAddress returns the signerAddress corresponding to LocalSigner's private key
func (l *LocalSigner) GetSignerAddress() common.Address {
	return crypto.PubkeyToAddress(l.privateKey.PublicKey)
}

// Sign mimicks the signing of `eth_sign` locally its supplied private key
func (l *LocalSigner) Sign(message []byte, signerAddress common.Address) (*ECSignature, error) {
	expectedSignerAddress := l.GetSignerAddress()
	if signerAddress != expectedSignerAddress {
		return nil, fmt.Errorf("Cannot sign with signerAddress %s since LocalSigner contains private key for %s", signerAddress, expectedSignerAddress)
	}

	messageWithPrefix, _ := accounts.TextAndHash(message)

	// The produced signature is in the [R || S || V] format where V is 0 or 1.
	signatureBytes, err := crypto.Sign(messageWithPrefix, l.privateKey)
	if err != nil {
		return nil, err
	}

	vParam := signatureBytes[64]
	if vParam == byte(0) {
		vParam = byte(27)
	} else if vParam == byte(1) {
		vParam = byte(28)
	}

	ecSignature := &ECSignature{
		V: vParam,
		R: common.BytesToHash(signatureBytes[0:32]),
		S: common.BytesToHash(signatureBytes[32:64]),
	}

	return ecSignature, nil
}

// TestSigner generates `eth_sign` signatures for test accounts available on the test
// Ethereum node Ganache
type TestSigner struct{}

// NewTestSigner instantiates a new LocalSigner
func NewTestSigner() Signer {
	return &TestSigner{}
}

// Sign generates an `eth_sign` equivalent signature using an public/private key
// pair hard-coded in the constants package.
func (t *TestSigner) Sign(message []byte, signerAddress common.Address) (*ECSignature, error) {
	pkBytes, ok := constants.GanacheAccountToPrivateKey[signerAddress]
	if !ok {
		return nil, errors.New("Unrecognized Ganache account supplied to ECSignForTests")
	}
	privateKey, err := crypto.ToECDSA(pkBytes)
	if err != nil {
		return nil, err
	}

	localSigner := NewLocalSigner(privateKey)
	return localSigner.Sign(message, signerAddress)
}
