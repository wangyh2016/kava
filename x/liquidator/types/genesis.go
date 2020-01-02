package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GenesisState is the state that must be provided at genesis.
type GenesisState struct {
	Params             Params             `json:"params" yaml:"params"`
	CollateralDeposits CollateralDeposits `json:"collateral_deposits" yaml:"collateral_deposits"`
	AuctionCollateral  sdk.Coins          `json:"auction_collateral" yaml:"auction_collateral"`
}

// CollateralDeposit stores the auction ID and it's corresponding deposit
type CollateralDeposit struct {
	Deposit   sdk.Coin
	AuctionID uint64
}

// CollateralDeposits array of CollateralDeposit
type CollateralDeposits []CollateralDeposit

// DefaultGenesisState returns a default genesis state
// TODO pick better values
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params:             DefaultParams(),
		CollateralDeposits: CollateralDeposits{},
		AuctionCollateral:  sdk.Coins{},
	}
}

// ValidateGenesis performs basic validation of genesis gs returning an error for any failed validation criteria.
func ValidateGenesis(gs GenesisState) error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}
	return nil
}
