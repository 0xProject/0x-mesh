package configs

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
)

/**
 * General
 */

// GanacheExchangeAddress specifies the 0x Exchange contract address on the Ganache snapshot
var GanacheExchangeAddress = common.HexToAddress("0x48bacb9266a570d521063ef5dd96e61686dbe788")

/**
 * OrderWatcher configs
 */

// MinCleanupInterval specified the minimum amount of time between orderbook cleanup intervals. These
// cleanups are meant to catch any stale orders that somehow were not caught by the event watcher
// process.
var MinCleanupInterval = 1 * time.Hour
