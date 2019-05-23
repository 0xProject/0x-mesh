package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// ECSignature contains the parameters of an elliptic curve signature
type ECSignature struct {
	V byte
	R []byte
	S []byte
}

// ECSign signs a message via the `eth_sign` Ethereum JSON-RPC call
func ECSign(message []byte, signerAddress common.Address, rpcClient *rpc.Client) (*ECSignature, error) {
	var signatureHex string
	if err := rpcClient.Call(&signatureHex, "eth_sign", signerAddress.Hex(), common.Bytes2Hex(message)); err != nil {
		return nil, err
	}
	// `eth_sign` returns the signature as r+s+v and the `v` parameter as 0/1 instead of 27/28
	signatureBytes := common.Hex2Bytes(signatureHex[2:])
	vParam := signatureBytes[64]
	if vParam == byte(0) {
		vParam = byte(27)
	} else if vParam == byte(1) {
		vParam = byte(28)
	}
	rParam := signatureBytes[0:32]
	sParam := signatureBytes[32:64]

	ecSignature := &ECSignature{
		V: vParam,
		R: rParam,
		S: sParam,
	}
	return ecSignature, nil
}
