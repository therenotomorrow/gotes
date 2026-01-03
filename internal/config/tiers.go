package config

import (
	"github.com/therenotomorrow/ex"
)

const (
	ErrInvalidTier ex.Error = "invalid tier"
)

type Tier string

const (
	TierDev  Tier = "dev"
	TierRC   Tier = "rc"
	TierProd Tier = "prod"
	TierTest Tier = "test"
)

func (t Tier) Validate() error {
	switch t {
	case TierDev, TierRC, TierProd, TierTest:
		return nil
	default:
		return ErrInvalidTier
	}
}
