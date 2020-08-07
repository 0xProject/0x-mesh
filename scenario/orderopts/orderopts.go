package orderopts

import (
	"math/big"

	"github.com/0xProject/0x-mesh/zeroex"
	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	Order             *zeroex.Order
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
		return nil
	}
}

func MakerAssetData(assetData []byte) Option {
	return func(cfg *Config) error {
		cfg.Order.MakerAssetData = assetData
		return nil
	}
}

func MakerAssetAmount(amount *big.Int) Option {
	return func(cfg *Config) error {
		cfg.Order.MakerAssetAmount = amount
		return nil
	}
}

func TakerAssetData(assetData []byte) Option {
	return func(cfg *Config) error {
		cfg.Order.TakerAssetData = assetData
		return nil
	}
}

func TakerAssetAmount(amount *big.Int) Option {
	return func(cfg *Config) error {
		cfg.Order.TakerAssetAmount = amount
		return nil
	}
}

func ExpirationTimeSeconds(expirationTimeSeconds *big.Int) Option {
	return func(cfg *Config) error {
		cfg.Order.ExpirationTimeSeconds = expirationTimeSeconds
		return nil
	}
}

func MakerFeeAssetData(assetData []byte) Option {
	return func(cfg *Config) error {
		cfg.Order.MakerFeeAssetData = assetData
		return nil
	}
}

func MakerFee(amount *big.Int) Option {
	return func(cfg *Config) error {
		cfg.Order.MakerFee = amount
		return nil
	}
}

func SenderAddress(address common.Address) Option {
	return func(cfg *Config) error {
		cfg.Order.SenderAddress = address
		return nil
	}
}

func TakerAddress(address common.Address) Option {
	return func(cfg *Config) error {
		cfg.Order.TakerAddress = address
		return nil
	}
}

func FeeRecipientAddress(address common.Address) Option {
	return func(cfg *Config) error {
		cfg.Order.FeeRecipientAddress = address
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
