package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cdptypes "github.com/kava-labs/kava/x/cdp/types"
)

// PartialDeposit stores deposits that are pending being sent to auction
type PartialDeposit struct {
	cdptypes.Deposit
	DebtAmount sdk.Int `json:"debt_amount" yaml:"debt_amount"`
}

// PartialDeposits array of PartialDeposit
type PartialDeposits []PartialDeposit

// SumDeposits returns the sum of all deposits in the input collection
func SumDeposits(ds []cdptypes.Deposit) sdk.Int {
	// Sum returns the total amount of deposits
	total := sdk.ZeroInt()
	for _, d := range ds {
		total = total.Add(d.Amount[0].Amount)
	}
	return total
}

// SumPartialDeposits returns the sum of all deposits in the input collection
func SumPartialDeposits(pd PartialDeposits) sdk.Int {
	// Sum returns the total amount of deposits
	total := sdk.ZeroInt()
	for _, d := range pd {
		total = total.Add(d.Amount[0].Amount)
	}
	return total
}

// SumCoins returns the sum of all amounts for the input coins
func SumCoins(cs sdk.Coins) sdk.Int {
	total := sdk.ZeroInt()
	for _, c := range cs {
		total = total.Add(c.Amount)
	}
	return total
}

// SumDebt returns the sum of all debt amounts for the input partial deposits
func SumDebt(pds PartialDeposits) sdk.Int {
	total := sdk.ZeroInt()
	for _, pd := range pds {
		total = total.Add(pd.DebtAmount)
	}
	return total
}

// CalculateRatableDebtShare calculates the ratable distribution of the debt based on the partial deposit amount.
func CalculateRatableDebtShare(depositAmount, totalDepositAmount, debtAmount sdk.Int) sdk.Int {
	return getShareWithFloatRounded(sdk.NewUintFromBigInt(depositAmount.BigInt()), sdk.NewUintFromBigInt(totalDepositAmount.BigInt()), sdk.NewUintFromBigInt((debtAmount.BigInt())))
}
