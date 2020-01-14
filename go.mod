module github.com/0xProject/0x-mesh

go 1.12

replace (
	github.com/ethereum/go-ethereum => github.com/0xProject/go-ethereum v1.8.8-0.20191104224527-9d5c202240be
	github.com/libp2p/go-ws-transport => github.com/libp2p/go-ws-transport v0.0.0-20191008032742-3098bba549e8
	github.com/syndtr/goleveldb => github.com/0xProject/goleveldb v1.0.1-0.20191115232649-6a187a47701c
)

require (
	github.com/albrow/stringset v2.1.0+incompatible
	github.com/allegro/bigcache v0.0.0-20190618191010-69ea0af04088 // indirect
	github.com/aristanetworks/goarista v0.0.0-20190712234253-ed1100a1c015 // indirect
	github.com/benbjohnson/clock v0.0.0-20161215174838-7dc76406b6d3
	github.com/cespare/cp v1.1.1 // indirect
	github.com/chromedp/cdproto v0.0.0-20190827000638-b5ac1e37ce90
	github.com/chromedp/chromedp v0.4.0
	github.com/deckarep/golang-set v1.7.1 // indirect
	github.com/edsrzf/mmap-go v1.0.0 // indirect
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/ethereum/go-ethereum v0.0.0-00010101000000-000000000000
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/gballet/go-libpcsclite v0.0.0-20190528105824-2fd9b619dd3c // indirect
	github.com/gibson042/canonicaljson-go v1.0.3
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/uuid v1.1.1
	github.com/hashicorp/golang-lru v0.5.3
	github.com/ipfs/go-datastore v0.1.1
	github.com/ipfs/go-ds-leveldb v0.0.2
	github.com/jpillora/backoff v0.0.0-20170918002102-8eab2debe79d
	github.com/karalabe/usb v0.0.0-20191104083709-911d15fe12a9 // indirect
	github.com/karlseguin/ccache v2.0.3+incompatible
	github.com/karlseguin/expect v1.0.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/lib/pq v1.2.0
	github.com/libp2p/go-libp2p v0.3.1
	github.com/libp2p/go-libp2p-autonat-svc v0.1.0
	github.com/libp2p/go-libp2p-circuit v0.1.2
	github.com/libp2p/go-libp2p-connmgr v0.1.1
	github.com/libp2p/go-libp2p-core v0.2.4
	github.com/libp2p/go-libp2p-crypto v0.1.0
	github.com/libp2p/go-libp2p-discovery v0.2.0
	github.com/libp2p/go-libp2p-kad-dht v0.2.1
	github.com/libp2p/go-libp2p-peer v0.2.0
	github.com/libp2p/go-libp2p-peerstore v0.1.3
	github.com/libp2p/go-libp2p-pubsub v0.2.3
	github.com/libp2p/go-libp2p-swarm v0.2.2
	github.com/libp2p/go-maddr-filter v0.0.5
	github.com/libp2p/go-ws-transport v0.1.0
	github.com/mattn/go-colorable v0.1.2 // indirect
	github.com/mattn/go-runewidth v0.0.4 // indirect
	github.com/multiformats/go-multiaddr v0.1.2
	github.com/multiformats/go-multiaddr-dns v0.2.0
	github.com/ocdogan/rbt v0.0.0-20160425054511-de6e2b48be33
	github.com/olekukonko/tablewriter v0.0.1 // indirect
	github.com/opaolini/go-ds-sql v0.0.0-20191105113501-15a4c8e2fea5
	github.com/opentracing/opentracing-go v1.1.0 // indirect
	github.com/pborman/uuid v0.0.0-20180906182336-adf5a7427709 // indirect
	github.com/plaid/go-envvar v1.1.0
	github.com/prometheus/tsdb v0.10.0 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/status-im/keycard-go v0.0.0-20190424133014-d95853db0f48 // indirect
	github.com/steakknife/bloomfilter v0.0.0-20180906043351-99ee86d9200f // indirect
	github.com/steakknife/hamming v0.0.0-20180906055317-003c143a81c2 // indirect
	github.com/stretchr/testify v1.3.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/tyler-smith/go-bip39 v1.0.2 // indirect
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208 // indirect
	github.com/wsxiaoys/terminal v0.0.0-20160513160801-0940f3fc43a0 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190809123943-df4f5c81cb3b // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/net v0.0.0-20190724013045-ca1201d0de80 // indirect
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/xerrors v0.0.0-20191011141410-1b5146add898 // indirect
	gopkg.in/karlseguin/expect.v1 v1.0.1 // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20190709231704-1e4459ed25ff // indirect
)
