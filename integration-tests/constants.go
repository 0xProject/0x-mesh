package integrationtests

import "time"

const (
	ethereumRPCURL  = "http://localhost:8545"
	ethereumChainID = 1337

	graphQLServerAddr                          = "localhost:60501"
	graphQLServerURL                           = "http://localhost:60501/graphql"
	browserIntegrationTestDataDir              = "./data/standalone-0"
	standaloneBlockPollingInterval             = "200ms"
	standaloneEthereumRPCMaxRequestsPer24HrUtc = "550000"

	// Various config options/information for the bootstrap node. The private key
	// for the bootstrap node is checked in to version control so we know it's
	// peer ID ahead of time.
	bootstrapAddr    = "/ip4/127.0.0.1/tcp/60500/ws"
	bootstrapList    = "/ip4/127.0.0.1/tcp/60500/ws/ipfs/16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7"
	bootstrapDataDir = "./data/bootstrap-0"

	// serverStartWaitTime is the amount of time to wait after seeing the "starting GraphQL server"
	// log message before attempting to connect to the server.
	serverStartWaitTime = 100 * time.Millisecond

	// blockProcessingWaitTime is the amount of time to wait for blockwatcher and orderwatcher to process block
	// events. Creating a valid order involves transferring sufficient funds to the maker, and setting their allowance for
	// the maker asset. These transactions must be mined and Mesh's BlockWatcher poller must process these blocks
	// in order for the order validation run at order submission to occur at a block number equal or higher then
	// the one where these state changes were included. With the BlockWatcher poller configured to run every 200ms,
	// we wait 500ms to give it ample time to run before submitting the above order to the Mesh node.
	blockProcessingWaitTime = 500 * time.Millisecond
)
