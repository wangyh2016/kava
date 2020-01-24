# End Block

At the end of each block, auctions that have reached `EndTime` are closed. The logic to close auctions is as follows:

```go
// EndBlocker runs at the end of every block.
func EndBlocker(ctx sdk.Context, k Keeper) {
	err := k.CloseExpiredAuctions(ctx)
	if err != nil {
		panic(err)
	}
}

// CloseExpiredAuctions finds all auctions that are past (or at) their ending times and closes them, paying out to the highest bidder.
func (k Keeper) CloseExpiredAuctions(ctx sdk.Context) sdk.Error {
	var expiredAuctions []uint64
	k.IterateAuctionsByTime(ctx, ctx.BlockTime(), func(id uint64) bool {
		expiredAuctions = append(expiredAuctions, id)
		return false
	})
	for _, id := range expiredAuctions {
		if err := k.CloseAuction(ctx, id); err != nil {
			return err
		}
	}
	return nil
}
```
