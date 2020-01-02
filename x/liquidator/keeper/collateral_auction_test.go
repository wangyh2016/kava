package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/kava-labs/kava/app"
	"github.com/kava-labs/kava/x/auction"
	"github.com/kava-labs/kava/x/liquidator/keeper"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtime "github.com/tendermint/tendermint/types/time"
)

type AuctionTestSuite struct {
	suite.Suite

	keeper keeper.Keeper
	addrs  []sdk.AccAddress
	app    app.TestApp
	ctx    sdk.Context
}

func (suite *AuctionTestSuite) SetupTest() {
	tApp := app.NewTestApp()
	ctx := tApp.NewContext(true, abci.Header{Height: 1, Time: tmtime.Now()})
	_, addrs := app.GeneratePrivKeyAddressPairs(1)
	coins := []sdk.Coins{cs(c("xrp", 5000000000))}
	authGS := app.NewAuthGenState(addrs, coins)
	tApp.InitializeFromGenesisStates(
		authGS,
		NewPricefeedGenStateMulti(),
		NewCDPGenStateMulti(),
		NewAuctionGenStateMulti(),
		NewLiquidatorGenStateMulti(),
	)
	cdpKeeper := tApp.GetCDPKeeper()
	err := cdpKeeper.AddCdp(ctx, addrs[0], cs(c("xrp", 5000000000)), cs(c("usdx", 625000000)))
	suite.NoError(err)
	testCDP, f := cdpKeeper.GetCDP(ctx, "xrp", 1)
	suite.True(f)
	cdpKeeper.SeizeCollateral(ctx, testCDP)
	keeper := tApp.GetLiquidatorKeeper()
	suite.app = tApp
	suite.ctx = ctx
	suite.keeper = keeper
}

func (suite *AuctionTestSuite) TestSendCollateralToAuction() {

	cdpKeeper := suite.app.GetCDPKeeper()
	deps := cdpKeeper.GetAllLiquidatedDeposits(suite.ctx)
	suite.keeper.SendCollateralToAuction(suite.ctx, deps)
	deps = cdpKeeper.GetAllLiquidatedDeposits(suite.ctx)
	suite.Equal(0, len(deps))
	auctionKeeper := suite.app.GetAuctionKeeper()
	ac, f := auctionKeeper.GetAuction(suite.ctx, 0)
	suite.True(f)
	fra, _ := ac.(auction.ForwardReverseAuction)
	suite.Equal(c("usdx", 625000000), fra.MaxBid)
	suite.Equal(c("xrp", 5000000000), fra.Lot)
}

func TestSeizeTestSuite(t *testing.T) {
	suite.Run(t, new(AuctionTestSuite))
}
