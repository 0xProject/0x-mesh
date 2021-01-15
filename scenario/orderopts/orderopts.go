package orderopts

import (
	"math/big"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Order             *zeroex.Order
	OrderV4           *zeroex.OrderV4
	SetupMakerState   bool
	SetupTakerAddress common.Address
}

type Option func(order *Config) error

// Apply applies the given options to the config, returning the first error
// encountered (if any).
func (cfg *Config) Apply(opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(cfg); err != nil {
			return err
		}
	}
	return nil
}

func MakerAddress(address common.Address) Option {
	return func(cfg *Config) error {
		cfg.Order.MakerAddress = address
		cfg.OrderV4.Maker = address
		return nil
	}
}

func MakerAssetData(assetData []byte) Option {
	return func(cfg *Config) error {
		cfg.Order.MakerAssetData = assetData
		// TODO: Could extract maker token and set on V4 order
		return nil
	}
}

func MakerAssetAmount(amount *big.Int) Option {
	return func(cfg *Config) error {
		cfg.Order.MakerAssetAmount = amount
		cfg.OrderV4.MakerAmount = amount
		return nil
	}
}

func TakerAssetData(assetData []byte) Option {
	return func(cfg *Config) error {
		cfg.Order.TakerAssetData = assetData
		// TODO: Could extract taker token and set on V4 order
		return nil
	}
}

func TakerAssetAmount(amount *big.Int) Option {
	return func(cfg *Config) error {
		cfg.Order.TakerAssetAmount = amount
		cfg.OrderV4.TakerAmount = amount
		return nil
	}
}

func ExpirationTimeSeconds(expirationTimeSeconds *big.Int) Option {
	return func(cfg *Config) error {
		cfg.Order.ExpirationTimeSeconds = expirationTimeSeconds
		cfg.OrderV4.Expiry = expirationTimeSeconds
		return nil
	}
}

func MakerFeeAssetData(assetData []byte) Option {
	return func(cfg *Config) error {
		cfg.Order.MakerFeeAssetData = assetData
		// V4 has no separate fee tokens
		return nil
	}
}

func MakerFee(amount *big.Int) Option {
	return func(cfg *Config) error {
		cfg.Order.MakerFee = amount
		// V4 has no maker fee
		return nil
	}
}

func SenderAddress(address common.Address) Option {
	return func(cfg *Config) error {
		cfg.Order.SenderAddress = address
		cfg.OrderV4.Sender = address
		return nil
	}
}

func TakerAddress(address common.Address) Option {
	return func(cfg *Config) error {
		cfg.Order.TakerAddress = address
		cfg.OrderV4.Taker = address
		return nil
	}
}

func FeeRecipientAddress(address common.Address) Option {
	return func(cfg *Config) error {
		cfg.Order.FeeRecipientAddress = address
		cfg.OrderV4.FeeRecipient = address
		return nil
	}
}

func SetupMakerState(b bool) Option {
	return func(cfg *Config) error {
		cfg.SetupMakerState = b
		return nil
	}
}

func SetupTakerAddress(takerAddress common.Address) Option {
	return func(cfg *Config) error {
		cfg.SetupTakerAddress = takerAddress
		return nil
	}
}
