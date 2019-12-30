package types

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetShareWithFloatRounded calculates the ratable distribution of an allocation based on the part contributed to the total.
func getShareWithFloatRounded(part, total, allocation sdk.Uint) sdk.Int {

	return sdk.NewInt(int64(sdk.NewUint(uint64(math.Round((float64(allocation.Uint64()) / (float64(total.Uint64()) / float64(part.Uint64())))))).Uint64()))
}

// GetSharesPercentage calculates the ratable percentage of shares that should be allocated based on the part contributed to the whole
func CalculateSharesPercentage(part, total sdk.Int) sdk.Dec {
	allocation := sdk.NewUint(100000000000000)
	shares := getShareWithFloatRounded(sdk.NewUintFromBigInt(part.BigInt()), sdk.NewUintFromBigInt(total.BigInt()), allocation)
	return sdk.NewDecWithPrec(shares.Int64(), 14)
}
