package keeper_test

import (
	"github.com/mokitanetwork/aether/x/savings/types"
)

func (suite *KeeperTestSuite) TestGetSetParams() {
	params := suite.keeper.GetParams(suite.ctx)
	suite.Require().Equal(
		types.Params{SupportedDenoms: []string(nil)},
		params,
	)

	newParams := types.NewParams([]string{"btc", "test"})
	suite.keeper.SetParams(suite.ctx, newParams)

	fetchedParams := suite.keeper.GetParams(suite.ctx)
	suite.Require().Equal(newParams, fetchedParams)
}
