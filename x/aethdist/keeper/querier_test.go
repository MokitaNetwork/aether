package keeper_test

// import (
// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	abci "github.com/tendermint/tendermint/abci/types"

// 	"github.com/mokitanetwork/aether/x/aethdist/keeper"
// 	"github.com/mokitanetwork/aether/x/aethdist/types"
// )

// func (suite *KeeperTestSuite) TestQuerierGetParams() {
// 	querier := keeper.NewQuerier(suite.keeper)
// 	bz, err := querier(suite.ctx, []string{types.QueryGetParams}, abci.RequestQuery{})
// 	suite.Require().NoError(err)
// 	suite.NotNil(bz)

// 	testParams := types.NewParams(true, testPeriods)
// 	var p types.Params
// 	suite.Nil(types.ModuleCdc.UnmarshalJSON(bz, &p))
// 	suite.Require().Equal(testParams, p)
// }

// func (suite *KeeperTestSuite) TestQuerierGetBalance() {
// 	sk := suite.supplyKeeper

// 	sk.MintCoins(suite.ctx, types.AetherDistMacc, sdk.NewCoins(sdk.NewCoin("uaeth", sdk.NewInt(100e6))))

// 	querier := keeper.NewQuerier(suite.keeper)
// 	bz, err := querier(suite.ctx, []string{types.QueryGetBalance}, abci.RequestQuery{})
// 	suite.Require().NoError(err)
// 	suite.Require().NotNil(bz)

// 	var coins sdk.Coins
// 	types.ModuleCdc.UnmarshalJSON(bz, &coins)
// 	suite.Require().Equal(sdk.NewInt(100e6), coins.AmountOf("uaeth"))
// }
