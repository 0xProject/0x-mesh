package integrationtests

const (
	ethereumRPCURL  = "http://localhost:8545"
	ethereumChainID = 1337
	wsRPCPort       = 60501
	httpRPCPort     = 60701

	browserIntegrationTestDataDir              = "./data/standalone-0"
	standaloneWSRPCEndpointPrefix              = "ws://localhost:"
	standaloneHTTPRPCEndpointPrefix            = "http://localhost:"
	standaloneRPCAddrPrefix                    = "localhost:"
	standaloneBlockPollingInterval             = "200ms"
	standaloneEthereumRPCMaxRequestsPer24HrUtc = "550000"

	// Various config options/information for the bootstrap node. The private key
	// for the bootstrap node is checked in to version control so we know it's
	// peer ID ahead of time.
	bootstrapAddr    = "/ip4/127.0.0.1/tcp/60500/ws"
	bootstrapPeerID  = "16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7"
	bootstrapList    = "/ip4/127.0.0.1/tcp/60500/ws/ipfs/16Uiu2HAmGd949LwaV4KNvK2WDSiMVy7xEmW983VH75CMmefmMpP7"
	bootstrapDataDir = "./data/bootstrap-0"
)
