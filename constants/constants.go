package constants

import (
	"github.com/ethereum/go-ethereum/common"
)

/**
 * General
 */

// GanacheEndpoint specifies the Ganache test Ethereum node JSON RPC endpoint used in tests
const GanacheEndpoint = "http://localhost:8545"

// GanacheExchangeAddress specifies the 0x Exchange contract address on the Ganache snapshot
var GanacheExchangeAddress = common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788")
