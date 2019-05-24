package ethereum

import (
	"crypto/ecdsa"
	"errors"

	"github.com/0xProject/0x-mesh/constants"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

// ECSignature contains the parameters of an elliptic curve signature
type ECSignature struct {
	V byte
	R common.Hash
	S common.Hash
}

// EthSign signs a message via the `eth_sign` Ethereum JSON-RPC call
func EthSign(message []byte, signerAddress common.Address, rpcClient *rpc.Client) (*ECSignature, error) {
	var signatureHex string
	if err := rpcClient.Call(&signatureHex, "eth_sign", signerAddress.Hex(), common.Bytes2Hex(message)); err != nil {
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

// LocalEthSign mimicks the signing of `eth_sign` locally using the supplied private key
func LocalEthSign(message []byte, privateKey *ecdsa.PrivateKey) (*ECSignature, error) {

	messageWithPrefix, _ := accounts.TextAndHash(message)

	// The produced signature is in the [R || S || V] format where V is 0 or 1.
	signatureBytes, err := crypto.Sign(messageWithPrefix, privateKey)
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

// EthSignForTests generates an `eth_sign` equivalent signature using an public/private key
// pair hard-coded in the constants package.
func EthSignForTests(message []byte, signerAddress common.Address) (*ECSignature, error) {
	pkBytes, ok := constants.GanacheAccountToPrivateKey[signerAddress]
	if !ok {
		return nil, errors.New("Unrecognized Ganache account supplied to ECSignForTests")
	}
	privateKey, err := crypto.ToECDSA(pkBytes)
	if err != nil {
		return nil, err
	}

	return LocalEthSign(message, privateKey)
}
