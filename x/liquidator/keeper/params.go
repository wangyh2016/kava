package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/liquidator/types"
)

// GetParams returns the params for liquidator module
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var params types.Params
	k.paramSubspace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets params for the liquidator module
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetCollateral returns the collateral param with corresponding denom
func (k Keeper) GetCollateral(ctx sdk.Context, denom string) (types.CollateralParam, bool) {
	params := k.GetParams(ctx)
	for _, cp := range params.CollateralParams {
		if cp.Denom == denom {
			return cp, true
		}
	}
	return types.CollateralParam{}, false
}
