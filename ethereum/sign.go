package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// Signer defines the methods needed to act as a elliptic curve signer
type Signer interface {
	EthSign(message []byte, signerAddress common.Address) (*ECSignature, error)
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

// EthSign signs a message via the `eth_sign` Ethereum JSON-RPC call
func (e *EthRPCSigner) EthSign(message []byte, signerAddress common.Address) (*ECSignature, error) {
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
