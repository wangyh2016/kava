package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/kava-labs/kava/x/auction"
	cdptypes "github.com/kava-labs/kava/x/cdp/types"
)

// CdpKeeper expected interface for the cdp keeper (noalias)
type CdpKeeper interface {
	GetCDP(ctx sdk.Context, collateralDenom string, cdpID uint64) (cdptypes.CDP, bool)
	DeleteCDP(ctx sdk.Context, cdp cdptypes.CDP)
	GetCdpByOwnerAndDenom(ctx sdk.Context, owner sdk.AccAddress, denom string) (cdptypes.CDP, bool)
	RemoveCdpOwnerIndex(ctx sdk.Context, cdp cdptypes.CDP)
	GetDeposits(ctx sdk.Context, cdpID uint64) (deposits cdptypes.Deposits)
	DeleteDeposit(ctx sdk.Context, status cdptypes.DepositStatus, cdpID uint64, depositor sdk.AccAddress)
	GetAllLiquidatedDeposits(ctx sdk.Context)
}

// SupplyKeeper defines the expected supply keeper for module accounts (noalias)
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) sdk.Error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) sdk.Error
	SendCoinsFromModuleToModule(ctx sdk.Context, senderModule, recipientModule string, amt sdk.Coins) sdk.Error
	BurnCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) sdk.Error
}

// AuctionKeeper expected interface for the auction keeper (noalias)
type AuctionKeeper interface {
	StartForwardAuction(sdk.Context, sdk.AccAddress, sdk.Coin, sdk.Coin) (auction.ID, sdk.Error)
	StartReverseAuction(sdk.Context, sdk.AccAddress, sdk.Coin, sdk.Coin) (auction.ID, sdk.Error)
	StartForwardReverseAuction(ctx sdk.Context, sellerModule string, lot sdk.Coin, maxBid sdk.Coin, deposits Deposits) (auction.ID, sdk.Error)
}
