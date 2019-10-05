package orderwatch

import (
	"github.com/0xProject/0x-mesh/zeroex/orderwatch/decoder"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// GetRelevantTopics returns the OrderWatcher-relevant topics that should be used when filtering
// the logs retrieved for Ethereum blocks
func GetRelevantTopics() []common.Hash {
	topics := []common.Hash{}
	for _, signature := range decoder.EVENT_SIGNATURES {
		topic := common.BytesToHash(crypto.Keccak256([]byte(signature)))
		topics = append(topics, topic)
	}

	return topics
}
