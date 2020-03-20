package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/bep3/types"
)

// TODO: Consider gov proposal to remove asset - must clear out asset supply amounts?

// GetAssetSupplyInfo builds a new AssetSupplyInfo from data in the store
func (k Keeper) GetAssetSupplyInfo(ctx sdk.Context, denom string) (types.AssetSupplyInfo, sdk.Error) {
	asset, found := k.GetAssetByDenom(ctx, denom)
	if !found {
		return types.AssetSupplyInfo{}, types.ErrAssetNotSupported(k.codespace, denom)
	}

	inSwapSupply, foundInSwapSupply := k.GetInSwapSupply(ctx, []byte(denom))
	if !foundInSwapSupply {
		inSwapSupply = sdk.NewInt64Coin(denom, 0)
	}

	assetSupply, foundAssetSupply := k.GetAssetSupply(ctx, []byte(denom))
	if !foundAssetSupply {
		assetSupply = sdk.NewInt64Coin(denom, 0)
	}

	assetSupplyInfo := types.NewAssetSupplyInfo(denom,
		inSwapSupply.Amount.Int64(),
		assetSupply.Amount.Int64(),
		asset.Limit.Int64(),
	)
	return assetSupplyInfo, nil
}

// IncrementAssetSupply increments an asset's supply by the coin
func (k Keeper) IncrementAssetSupply(ctx sdk.Context, coin sdk.Coin) {
	currSupply, _ := k.GetAssetSupply(ctx, []byte(coin.Denom))
	k.SetAssetSupply(ctx, currSupply.Add(coin), []byte(coin.Denom))
}

// DecrementAssetSupply decrement an asset's supply by the coin
func (k Keeper) DecrementAssetSupply(ctx sdk.Context, coin sdk.Coin) sdk.Error {
	currSupply, _ := k.GetAssetSupply(ctx, []byte(coin.Denom))
	if !currSupply.Amount.Sub(coin.Amount).GTE(sdk.NewInt(0)) {
		return sdk.ErrInternal("cannot decrement asset's supply such that it results in a negative amount")
	}
	k.SetAssetSupply(ctx, currSupply.Sub(coin), []byte(coin.Denom))
	return nil
}

// IncrementInSwapSupply increments an asset's in swap supply by the coin
func (k Keeper) IncrementInSwapSupply(ctx sdk.Context, coin sdk.Coin) {
	currInSwapSupply, _ := k.GetInSwapSupply(ctx, []byte(coin.Denom))
	k.SetInSwapSupply(ctx, currInSwapSupply.Add(coin), []byte(coin.Denom))
}

// DecrementInSwapSupply decrement an asset's in swap supply by the coin
func (k Keeper) DecrementInSwapSupply(ctx sdk.Context, coin sdk.Coin) sdk.Error {
	currInSwapSupply, _ := k.GetInSwapSupply(ctx, []byte(coin.Denom))
	if !currInSwapSupply.Amount.Sub(coin.Amount).GTE(sdk.NewInt(0)) {
		return sdk.ErrInternal("cannot decrement asset's in swap supply such that it results in a negative amount")
	}
	k.SetInSwapSupply(ctx, currInSwapSupply.Sub(coin), []byte(coin.Denom))
	return nil
}

// ValidateCreateSwapAgainstSupplyLimit validates the proposed swap's amount against the asset's total supply limit
func (k Keeper) ValidateCreateSwapAgainstSupplyLimit(ctx sdk.Context, coin sdk.Coin) sdk.Error {
	currInSwapSupply, currAssetSupply := k.LoadAssetSupply(ctx, coin.Denom)
	currTotalSupply := currInSwapSupply.Add(currAssetSupply)
	asset, _ := k.GetAssetByDenom(ctx, coin.Denom)
	if currTotalSupply.Add(coin).Amount.GT(asset.Limit) {
		return types.ErrAboveTotalAssetSupplyLimit(
			k.codespace, coin.Denom, asset.Limit, currAssetSupply.Amount, currInSwapSupply.Amount,
		)
	}
	return nil
}

// LoadAssetSupply loads an asset's in swap supply and its current supply.
// If it's the first swap of this asset type, initialize both in swap supply
// and asset supply to 0.
func (k Keeper) LoadAssetSupply(ctx sdk.Context, denom string) (sdk.Coin, sdk.Coin) {
	currAssetSupply, found := k.GetAssetSupply(ctx, []byte(denom))
	if !found {
		initialSupply := sdk.NewInt64Coin(denom, 0)
		k.SetAssetSupply(ctx, initialSupply, []byte(denom))
		k.SetInSwapSupply(ctx, initialSupply, []byte(denom))
		return initialSupply, initialSupply
	}
	currInSwapSupply, _ := k.GetInSwapSupply(ctx, []byte(denom))
	return currInSwapSupply, currAssetSupply
}

// ValidateLiveAsset checks if an asset is both supported and active
func (k Keeper) ValidateLiveAsset(ctx sdk.Context, coin sdk.Coin) sdk.Error {
	asset, found := k.GetAssetByDenom(ctx, coin.Denom)
	if !found {
		return types.ErrAssetNotSupported(k.codespace, coin.Denom)
	}
	if !asset.Active {
		return types.ErrAssetNotActive(k.codespace, asset.Denom)
	}
	return nil
}
