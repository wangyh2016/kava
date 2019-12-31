package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/x/auction"
	cdptypes "github.com/kava-labs/kava/x/cdp/types"
	"github.com/kava-labs/kava/x/liquidator/types"
)

// SendCollateralToAuction iterate over liquidated deposits and sends them to the auction module
// deposits are broken down into lots so that the maximum amount sent to a single auction is always less than
// or equal to the parameter 'AuctionSize' for that collateral type. Deposits less than a lot are combined to
// form lots until only a partial lot remains, which is then sent to auction.
func (k Keeper) SendCollateralToAuction(ctx sdk.Context, deposits []cdptypes.Deposit) {
	cp, _ := k.GetCollateral(ctx, deposits[0].Amount[0].Denom)
	partialAuctionDeposits := types.PartialDeposits{}

	// while there is a positive amount of collateral left to auction
	for types.SumDeposits(deposits).GT(sdk.ZeroInt()) {
		// while there is at least one lot worth of collateral to auction
		for types.SumDeposits(deposits).GT(cp.AuctionSize) {
			for i, dep := range deposits {
				// ensure the cdp corresponding to that deposit still exists
				partialAuctionPending := false
				cdpFromDeposit, found := k.cdpKeeper.GetCDP(ctx, dep.Amount[0].Denom, dep.CdpID)
				if !found {
					panic(fmt.Sprintf("cdp not found for liquidated deposit: %s", dep))
				}
				// check if there are other deposits on the cdp
				otherDeposits := k.cdpKeeper.GetDeposits(ctx, dep.CdpID)
				// if there is 1 deposit, then the collateral of that deposit is covering 100% of the debt
				debtCoveredByDeposit := types.SumCoins(cdpFromDeposit.Principal.Add(cdpFromDeposit.AccumulatedFees))
				// if there are multiple deposits, figure out how much debt is covered by this particular deposit
				if len(otherDeposits) > 1 {
					debtCoveredByDeposit = types.CalculateRatableDebtShare(dep.Amount[0].Amount, types.SumDeposits(otherDeposits), debtCoveredByDeposit)
				}
				// if this deposit contains more than a lot, start auctions with it
				// until there is less than a lot remaining
				k.CreateAuctionsFromDeposit(ctx, &dep, &debtCoveredByDeposit, cp.AuctionSize, cdpFromDeposit.Principal[0].Denom)
				// the deposit at this point contains zero collateral or collateral less than one lot
				if !dep.Amount[0].Amount.IsZero() {
					// if adding this deposit to the partial deposit pool is still less than a lot,
					// append it to the lot
					if types.SumPartialDeposits(partialAuctionDeposits).Add(dep.Amount[0].Amount).LT(cp.AuctionSize) {
						partialAuctionDeposits = append(partialAuctionDeposits, types.PartialDeposit{dep, debtCoveredByDeposit})
						dep.Amount[0].Amount = sdk.ZeroInt()
						// flag that prevents the deposit and cdp from being removed from the store
						// until the partial deposits have been sent to auction
						partialAuctionPending = true
					} else {
						// if adding the deposit to the partial deposit pool creates
						// a lot, create a new auction from partial deposits
						k.CreateAuctionFromPartialDeposits(ctx, &dep, &debtCoveredByDeposit, &partialAuctionDeposits, cp.AuctionSize, dep.Amount[0].Denom, cdpFromDeposit.Principal[0].Denom)
						// go through the partial deposit pool (except the current deposit we just appended)
						// and delete the deposits from the store and delete the cdp if appropriate
						for i := 0; i < len(partialAuctionDeposits)-1; i++ {
							k.RemovePartialDeposit(ctx, partialAuctionDeposits[i])
						}
						// reset the partial deposits pool
						partialAuctionDeposits = types.PartialDeposits{}
					}
				}
				if dep.Amount[0].Amount.IsZero() {
					// remove the deposit from the slice if it is empty
					deposits = append(deposits[:i], deposits[i+1:]...)

					if !partialAuctionPending {
						// delete the deposit and cdp from the store if there is no partial auctions pending
						k.cdpKeeper.DeleteDeposit(ctx, cdptypes.StatusLiquidated, cdpFromDeposit.ID, dep.Depositor)
						// remove the deposit from the slice

						// remove the cdp from the store if there are no other deposits
						if len(otherDeposits) == 1 {
							k.cdpKeeper.DeleteCDP(ctx, cdpFromDeposit)
							k.cdpKeeper.RemoveCdpOwnerIndex(ctx, cdpFromDeposit)
						}
					}
				}
			}
		}
		// if there is a partial lot remaining, send it to auction
		if len(partialAuctionDeposits) > 0 {
			cdpFromDeposit, found := k.cdpKeeper.GetCdpByOwnerAndDenom(ctx, partialAuctionDeposits[0].Depositor, partialAuctionDeposits[0].Amount[0].Denom)
			if !found {
				panic(fmt.Sprintf("cdp not found for liquidated deposit: %s", partialAuctionDeposits[0]))
			}
			k.CreateAuctionFromPartialDeposits(ctx, &cdptypes.Deposit{}, &sdk.Int{}, &partialAuctionDeposits, types.SumPartialDeposits(partialAuctionDeposits), partialAuctionDeposits[0].Amount[0].Denom, cdpFromDeposit.Principal[0].Denom)
			// Delete the deposits from the store and delete the cdps from the store if appropriate
			for _, pd := range partialAuctionDeposits {
				k.RemovePartialDeposit(ctx, pd)
			}
		}
	}
}

// CreateAuctionsFromDeposit creates auctions from the input deposit until there is less than auctionSize left on the deposit
func (k Keeper) CreateAuctionsFromDeposit(ctx sdk.Context, dep *cdptypes.Deposit, debt *sdk.Int, auctionSize sdk.Int, principalDenom string) {
	for dep.Amount[0].Amount.GTE(auctionSize) {
		// figure out how much debt is covered by one lots worth of collateral
		depositDebtAmount := (auctionSize.Quo(dep.Amount[0].Amount)).Mul(*debt)
		// subtract one lot's worth of debt from the total debt covered by this deposit
		*debt = debt.Sub(depositDebtAmount)
		// start an auction for one lot, attempting to raise depositDebtAmount
		auctionID, err := k.auctionKeeper.StartForwardReverseAuction(
			ctx, types.ModuleName, sdk.NewCoin(dep.Amount[0].Denom, auctionSize),
			sdk.NewCoin(principalDenom, depositDebtAmount), []sdk.AccAddress{dep.Depositor},
			[]sdk.Int{depositDebtAmount})
		if err != nil {
			panic(err)
		}
		// store the ongoing auction with how much collateral is potentially being sold
		k.SetAuctionDeposit(ctx, auctionID, sdk.NewCoin(dep.Amount[0].Denom, auctionSize))
		// increment the total amount of collateral at auctions for this collateral type
		k.IncrementTotalAuctionDeposits(ctx, sdk.NewCoin(dep.Amount[0].Denom, auctionSize))
		dep.Amount[0].Amount = dep.Amount[0].Amount.Sub(auctionSize)
	}
}

func (k Keeper) CreateAuctionFromPartialDeposits(ctx sdk.Context, dep *cdptypes.Deposit, debt *sdk.Int, partialDeps *types.PartialDeposits, auctionSize sdk.Int, collateralDenom, principalDenom string) {
	if dep != (&cdptypes.Deposit{}) {
		collateralToAdd := auctionSize.Sub(types.SumPartialDeposits(*partialDeps))
		debtToAdd := collateralToAdd.Quo(dep.Amount[0].Amount).Mul(*debt)
		*debt = debt.Sub(debtToAdd)
		depositToAdd := cdptypes.NewDeposit(dep.CdpID, dep.Depositor, sdk.NewCoins(sdk.NewCoin(collateralDenom, collateralToAdd)))
		*partialDeps = append(*partialDeps, types.PartialDeposit{depositToAdd, debtToAdd})
		dep.Amount[0].Amount = dep.Amount[0].Amount.Sub(collateralToAdd)
	}
	depositors := []sdk.AccAddress{}
	deposits := []sdk.Int{}
	for _, dep := range *partialDeps {
		depositors = append(depositors, dep.Depositor)
		deposits = append(deposits, dep.DebtAmount)
	}
	auctionID, err := k.auctionKeeper.StartForwardReverseAuction(
		ctx, types.ModuleName, sdk.NewCoin(collateralDenom, auctionSize),
		sdk.NewCoin(principalDenom, types.SumDebt(*partialDeps)), depositors, deposits)
	if err != nil {
		panic(err)
	}
	// store the auction collateral and total collateral
	k.SetAuctionDeposit(ctx, auctionID, sdk.NewCoin(collateralDenom, auctionSize))
	k.IncrementTotalAuctionDeposits(ctx, sdk.NewCoin(dep.Amount[0].Denom, auctionSize))
}

func (k Keeper) RemovePartialDeposit(ctx sdk.Context, partialDep types.PartialDeposit) {
	cdpFromDeposit, found := k.cdpKeeper.GetCDP(ctx, partialDep.Amount[0].Denom, partialDep.CdpID)
	if !found {
		panic(fmt.Sprintf("cdp not found for liquidated deposit: %s", partialDep))
	}
	otherDeposits := k.cdpKeeper.GetDeposits(ctx, partialDep.CdpID)
	if len(otherDeposits) == 1 {
		k.cdpKeeper.DeleteCDP(ctx, cdpFromDeposit)
		k.cdpKeeper.RemoveCdpOwnerIndex(ctx, cdpFromDeposit)
	}
	k.cdpKeeper.DeleteDeposit(ctx, cdptypes.StatusLiquidated, partialDep.CdpID, partialDep.Depositor)
}

// SetAuctionDeposit stores the amount of collateral at auction for a particular auction
func (k Keeper) SetAuctionDeposit(ctx sdk.Context, auctionID auction.ID, auctionDeposit sdk.Coin) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.CollateralAuctionKeyPrefix)
	store.Set(types.GetAuctionIDBytes(uint64(auctionID)), k.cdc.MustMarshalBinaryLengthPrefixed(auctionDeposit))
}

// IncrementTotalAuctionDeposits increments the total amount of collateral at auction for a particular collateral type
func (k Keeper) IncrementTotalAuctionDeposits(ctx sdk.Context, auctionDeposit sdk.Coin) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.TotalCollateralPrefix)
	bz := store.Get([]byte(auctionDeposit.Denom))
	if bz == nil {
		store.Set([]byte(auctionDeposit.Denom), k.cdc.MustMarshalBinaryLengthPrefixed(auctionDeposit))
		return
	}
	var previousTotal sdk.Coin
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &previousTotal)
	previousTotal = previousTotal.Add(auctionDeposit)
	store.Set([]byte(auctionDeposit.Denom), k.cdc.MustMarshalBinaryLengthPrefixed(auctionDeposit))
	return
}
