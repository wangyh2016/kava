package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	cdptypes "github.com/kava-labs/kava/x/cdp/types"
)

type Deposit struct {
	Depositor sdk.AccAddress
	Weight    sdk.Dec
}
type Deposits []Deposit

type PartialDeposit struct {
	cdptypes.Deposit
	DebtAmount sdk.Int `json:"debt_amount" yaml:"debt_amount"`
}

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

func SumCoins(cs sdk.Coins) sdk.Int {
	total := sdk.ZeroInt()
	for _, c := range cs {
		total = total.Add(c.Amount)
	}
	return total
}

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

func ConvertDepositsToWeights(pds PartialDeposits) Deposits {
	// total amount of collateral being sent to auction
	totalAuctionDeposits := SumPartialDeposits(pds)
	weights := Deposits{}
	totalWeight := sdk.ZeroDec()
	for i, pd := range pds {
		collateralPercentage := CalculateSharesPercentage(pd.Amount[0].Amount, totalAuctionDeposits)
		totalWeight = totalWeight.Add(collateralPercentage)
		weights[i] = Deposit{pd.Depositor, collateralPercentage}
	}
	if !totalWeight.Equal(sdk.OneDec()) {
		error := sdk.OneDec().Sub(totalWeight)
		weights[0].Weight = weights[0].Weight.Add(error)
	}
	return weights
}
