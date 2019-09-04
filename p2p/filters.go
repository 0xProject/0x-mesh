package p2p

import (
	libp2p "github.com/libp2p/go-libp2p"
	filter "github.com/libp2p/go-maddr-filter"
)

func Filters(filters *filter.Filters) libp2p.Option {
	return func(cfg *libp2p.Config) error {
		cfg.Filters = filters
		return nil
	}
}
