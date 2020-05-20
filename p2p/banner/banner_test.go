package banner

import (
	"net"
	"testing"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPNetFromMaddr(t *testing.T) {
	testCases := []struct {
		maddr    ma.Multiaddr
		expected net.IPNet
	}{
		{
			maddr: newMaddr(t, "/ip4/159.65.4.82/tcp/60558"),
			expected: net.IPNet{
				IP:   net.IP{0x9f, 0x41, 0x4, 0x52},
				Mask: ipv4AllMask,
			},
		},
		{
			maddr: newMaddr(t, "/ip4/159.65.4.82/tcp/60558/ipfs/16Uiu2HAm9brLYhoM1wCTRtGRR7ZqXhk8kfEt6a2rSFSZpeV8eB7L/p2p-circuit"),
			expected: net.IPNet{
				IP:   net.IP{0x9f, 0x41, 0x4, 0x52},
				Mask: ipv4AllMask,
			},
		},
		{
			maddr: newMaddr(t, "/ip6/fe80:cd00:0000:0cde:1257:0000:211e:729c/tcp/60558"),
			expected: net.IPNet{
				IP:   net.IP{0xfe, 0x80, 0xcd, 0x0, 0x0, 0x0, 0xc, 0xde, 0x12, 0x57, 0x0, 0x0, 0x21, 0x1e, 0x72, 0x9c},
				Mask: ipv6AllMask,
			},
		},
		{
			maddr: newMaddr(t, "/ip6/fe80:cd00:0000:0cde:1257:0000:211e:729c/tcp/60558/ipfs/16Uiu2HAm9brLYhoM1wCTRtGRR7ZqXhk8kfEt6a2rSFSZpeV8eB7L/p2p-circuit"),
			expected: net.IPNet{
				IP:   net.IP{0xfe, 0x80, 0xcd, 0x0, 0x0, 0x0, 0xc, 0xde, 0x12, 0x57, 0x0, 0x0, 0x21, 0x1e, 0x72, 0x9c},
				Mask: ipv6AllMask,
			},
		},
	}

	for i, tc := range testCases {
		actual, err := ipNetFromMaddr(tc.maddr)
		require.NoError(t, err, "test case %d (%s)", i, tc.maddr.String())
		assert.Equal(t, tc.expected, actual, "test case %d (%s)", i, tc.maddr.String())
	}
}

func newMaddr(t *testing.T, s string) ma.Multiaddr {
	maddr, err := ma.NewMultiaddr(s)
	require.NoError(t, err)
	return maddr
}
