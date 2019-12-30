package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
)

// Parameter keys
var (
	KeyDebtAuctionSize  = []byte("DebtAuctionSize")
	KeyCollateralParams = []byte("CollateralParams")
)

// Params store params for the liquidator module
type Params struct {
	DebtAuctionSize sdk.Int
	//SurplusAuctionSize sdk.Int
	CollateralParams CollateralParams
}

// NewParams returns a new params object for the liquidator module
func NewParams(debtAuctionSize sdk.Int, collateralParams CollateralParams) Params {
	return Params{
		DebtAuctionSize:  debtAuctionSize,
		CollateralParams: collateralParams,
	}
}

// String implements fmt.Stringer
func (p Params) String() string {
	out := fmt.Sprintf(`Params:
		Debt Auction Size: %s
		Collateral Params: `,
		p.DebtAuctionSize,
	)
	for _, cp := range p.CollateralParams {
		out += fmt.Sprintf(`
		%s`, cp.String())
	}
	return out
}

// CollateralParam params storing information about each collateral for the liquidator module
type CollateralParam struct {
	Denom              string  // Coin name of collateral type
	AuctionSize        sdk.Int // Max amount of collateral to sell off in any one auction. Known as lump in Maker.
	LiquidationPenalty sdk.Dec
}

// String implements stringer interface
func (cp CollateralParam) String() string {
	return fmt.Sprintf(`
  Denom:        %s
	AuctionSize: %s
	Liquidation Penalty: %s`, cp.Denom, cp.AuctionSize, cp.LiquidationPenalty)
}

// CollateralParams array of CollateralParam
type CollateralParams []CollateralParam

// String implements fmt.Stringer
func (cps CollateralParams) String() string {
	out := "Collateral Params\n"
	for _, cp := range cps {
		out += fmt.Sprintf("%s\n", cp)
	}
	return out
}

// ParamKeyTable for the liquidator module
func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of liquidator module's parameters.
// nolint
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		subspace.NewParamSetPair(KeyDebtAuctionSize, &p.DebtAuctionSize),
		subspace.NewParamSetPair(KeyCollateralParams, &p.CollateralParams),
	}
}

// DefaultParams for the liquidator module
func DefaultParams() Params {
	return Params{
		DebtAuctionSize:  sdk.NewInt(1000),
		CollateralParams: CollateralParams{},
	}
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if p.DebtAuctionSize.IsNegative() {
		return fmt.Errorf("debt auction size should be positive, is %s", p.DebtAuctionSize)
	}
	denomDupMap := make(map[string]int)
	for _, cp := range p.CollateralParams {
		_, found := denomDupMap[cp.Denom]
		if found {
			return fmt.Errorf("duplicate denom: %s", cp.Denom)
		}
		denomDupMap[cp.Denom] = 1
		if cp.AuctionSize.IsNegative() {
			return fmt.Errorf(
				"auction size for each collateral should be positive, is %s for %s", cp.AuctionSize, cp.Denom,
			)
		}
		if cp.LiquidationPenalty.LT(sdk.ZeroDec()) || cp.LiquidationPenalty.GT(sdk.OneDec()) {
			return fmt.Errorf(
				"liquidation penalty should be between 0 and 1, is %s", cp.LiquidationPenalty)
		}
	}
	return nil
}
