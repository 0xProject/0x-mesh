package integrationtests

import (
	"math/big"

	"github.com/0xProject/0x-mesh/constants"
)

const (
	ethereumRPCURL  = "http://localhost:8545"
	ethereumChainID = 1337

	rpcPort        = 60501
	protocolString = "/meshsub/1.0.0"

	// Various config options/information for the bootstrap node. The private key
	// for the bootstrap node is checked in to version control so we know it's
	// peer ID ahead of time.
	bootstrapAddr    = "/ip4/127.0.0.1/tcp/60500/ws"
	bootstrapPeerID  = "16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7"
	bootstrapList    = "/ip4/127.0.0.1/tcp/60500/ws/ipfs/16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7"
	bootstrapDataDir = "./data/bootstrap-0"

	// Various config options/information for the standalone node. Like the
	// bootstrap node, we know the private key/peer ID ahead of time.
	standalonePeerID      = "16Uiu2HAmM9j68mgGGSFkXsuzbGJA8ezVHtQ2H9y6mRJAPhx6xtj9"
	standaloneDataDir     = "./data/standalone-"
	standaloneRPCEndpoint = "ws://localhost:"
	standaloneRPCAddr     = "localhost:"
)

var makerAddress = constants.GanacheAccount1
var takerAddress = constants.GanacheAccount2
var eighteenDecimalsInBaseUnits = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
var wethAmount = new(big.Int).Mul(big.NewInt(50), eighteenDecimalsInBaseUnits)
var zrxAmount = new(big.Int).Mul(big.NewInt(100), eighteenDecimalsInBaseUnits)
