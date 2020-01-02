package liquidator

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis sets the genesis state in the keeper.
func InitGenesis(ctx sdk.Context, k Keeper, gs GenesisState) {

	k.SetParams(ctx, gs.Params)
	for _, ac := range gs.AuctionCollateral {
		k.IncrementTotalAuctionDeposits(ctx, ac)
	}
	for _, cd := range gs.CollateralDeposits {
		k.SetAuctionDeposit(ctx, cd.AuctionID, cd.Deposit)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	params := k.GetParams(ctx)
	auctionCollateralDeposits := k.GetCollateralDeposits(ctx)
	totalAuctionCollateral := k.GetTotalAuctionCollateral(ctx)

	return GenesisState{
		Params:             params,
		CollateralDeposits: auctionCollateralDeposits,
		AuctionCollateral:  totalAuctionCollateral,
	}
}
