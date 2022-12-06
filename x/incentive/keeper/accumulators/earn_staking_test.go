package accumulators_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/mokitanetwork/aether/app"
	earntypes "github.com/mokitanetwork/aether/x/earn/types"
	"github.com/mokitanetwork/aether/x/incentive/testutil"
	"github.com/mokitanetwork/aether/x/incentive/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type EarnAccumulatorStakingRewardsTestSuite struct {
	testutil.IntegrationTester

	keeper    testutil.TestKeeper
	userAddrs []sdk.AccAddress
	valAddrs  []sdk.ValAddress
}

func TestEarnStakingRewardsIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(EarnAccumulatorStakingRewardsTestSuite))
}

func (suite *EarnAccumulatorStakingRewardsTestSuite) SetupTest() {
	suite.IntegrationTester.SetupTest()

	suite.keeper = testutil.TestKeeper{
		Keeper: suite.App.GetIncentiveKeeper(),
	}

	_, addrs := app.GeneratePrivKeyAddressPairs(5)
	suite.userAddrs = addrs[0:2]
	suite.valAddrs = []sdk.ValAddress{
		sdk.ValAddress(addrs[2]),
		sdk.ValAddress(addrs[3]),
	}

	// Setup app with test state
	authBuilder := app.NewAuthBankGenesisBuilder().
		WithSimpleAccount(addrs[0], cs(c("uaeth", 1e12))).
		WithSimpleAccount(addrs[1], cs(c("uaeth", 1e12))).
		WithSimpleAccount(addrs[2], cs(c("uaeth", 1e12))).
		WithSimpleAccount(addrs[3], cs(c("uaeth", 1e12)))

	incentiveBuilder := testutil.NewIncentiveGenesisBuilder().
		WithGenesisTime(suite.GenesisTime).
		WithSimpleRewardPeriod(types.CLAIM_TYPE_EARN, "baeth", cs())

	savingsBuilder := testutil.NewSavingsGenesisBuilder().
		WithSupportedDenoms("baeth")

	earnBuilder := testutil.NewEarnGenesisBuilder().
		WithAllowedVaults(earntypes.AllowedVault{
			Denom:             "baeth",
			Strategies:        earntypes.StrategyTypes{earntypes.STRATEGY_TYPE_SAVINGS},
			IsPrivateVault:    false,
			AllowedDepositors: nil,
		})

	stakingBuilder := testutil.NewStakingGenesisBuilder()

	mintBuilder := testutil.NewMintGenesisBuilder().
		WithInflationMax(sdk.OneDec()).
		WithInflationMin(sdk.OneDec()).
		WithMinter(sdk.OneDec(), sdk.ZeroDec()).
		WithMintDenom("uaeth")

	suite.StartChainWithBuilders(
		authBuilder,
		incentiveBuilder,
		savingsBuilder,
		earnBuilder,
		stakingBuilder,
		mintBuilder,
	)
}

func (suite *EarnAccumulatorStakingRewardsTestSuite) TestStakingRewardsDistributed() {
	// derivative 1: 8 total staked, 7 to earn, 1 not in earn
	// derivative 2: 2 total staked, 1 to earn, 1 not in earn
	userMintAmount0 := c("uaeth", 8e9)
	userMintAmount1 := c("uaeth", 2e9)

	userDepositAmount0 := i(7e9)
	userDepositAmount1 := i(1e9)

	// Create two validators
	derivative0, err := suite.MintLiquidAnyValAddr(suite.userAddrs[0], suite.valAddrs[0], userMintAmount0)
	suite.Require().NoError(err)

	derivative1, err := suite.MintLiquidAnyValAddr(suite.userAddrs[0], suite.valAddrs[1], userMintAmount1)
	suite.Require().NoError(err)

	err = suite.DeliverEarnMsgDeposit(suite.userAddrs[0], sdk.NewCoin(derivative0.Denom, userDepositAmount0), earntypes.STRATEGY_TYPE_SAVINGS)
	suite.NoError(err)
	err = suite.DeliverEarnMsgDeposit(suite.userAddrs[0], sdk.NewCoin(derivative1.Denom, userDepositAmount1), earntypes.STRATEGY_TYPE_SAVINGS)
	suite.NoError(err)

	// Get derivative denoms
	lq := suite.App.GetLiquidKeeper()
	vaultDenom1 := lq.GetLiquidStakingTokenDenom(suite.valAddrs[0])
	vaultDenom2 := lq.GetLiquidStakingTokenDenom(suite.valAddrs[1])

	previousAccrualTime := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.Ctx = suite.Ctx.WithBlockTime(previousAccrualTime)

	initialVault1RewardFactor := d("0.04")
	initialVault2RewardFactor := d("0.04")

	globalIndexes := types.MultiRewardIndexes{
		{
			CollateralType: vaultDenom1,
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "uaeth",
					RewardFactor:   initialVault1RewardFactor,
				},
			},
		},
		{
			CollateralType: vaultDenom2,
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "uaeth",
					RewardFactor:   initialVault2RewardFactor,
				},
			},
		},
	}

	suite.keeper.StoreGlobalIndexes(suite.Ctx, types.CLAIM_TYPE_EARN, globalIndexes)

	suite.keeper.Store.SetRewardAccrualTime(suite.Ctx, types.CLAIM_TYPE_EARN, vaultDenom1, suite.Ctx.BlockTime())
	suite.keeper.Store.SetRewardAccrualTime(suite.Ctx, types.CLAIM_TYPE_EARN, vaultDenom2, suite.Ctx.BlockTime())

	val := suite.GetAbciValidator(suite.valAddrs[0])

	// Mint tokens, distribute to validators, claim staking rewards
	// 1 hour later
	_, resBeginBlock := suite.NextBlockAfterWithReq(
		1*time.Hour,
		abci.RequestEndBlock{},
		abci.RequestBeginBlock{
			LastCommitInfo: abci.LastCommitInfo{
				Votes: []abci.VoteInfo{{
					Validator:       val,
					SignedLastBlock: true,
				}},
			},
		},
	)

	// check time and factors
	suite.StoredTimeEquals(types.CLAIM_TYPE_EARN, vaultDenom1, suite.Ctx.BlockTime())
	suite.StoredTimeEquals(types.CLAIM_TYPE_EARN, vaultDenom2, suite.Ctx.BlockTime())

	validatorRewards, _ := suite.GetBeginBlockClaimedStakingRewards(resBeginBlock)

	suite.Require().Contains(validatorRewards, suite.valAddrs[0].String(), "there should be claim events for validator 1")
	suite.Require().Contains(validatorRewards, suite.valAddrs[1].String(), "there should be claim events for validator 2")

	// Total staking rewards / total source shares (**deposited in earn** not total minted)
	// types.RewardIndexes.Quo() uses Dec.Quo() which uses bankers rounding.
	// So we need to use Dec.Quo() to also round vs Dec.QuoInt() which truncates
	expectedIndexes1 := validatorRewards[suite.valAddrs[0].String()].
		AmountOf("uaeth").
		ToDec().
		Quo(userDepositAmount0.ToDec())

	expectedIndexes2 := validatorRewards[suite.valAddrs[1].String()].
		AmountOf("uaeth").
		ToDec().
		Quo(userDepositAmount1.ToDec())

	// Only contains staking rewards
	suite.StoredIndexesEqual(types.CLAIM_TYPE_EARN, vaultDenom1, types.RewardIndexes{
		{
			CollateralType: "uaeth",
			RewardFactor:   initialVault1RewardFactor.Add(expectedIndexes1),
		},
	})

	suite.StoredIndexesEqual(types.CLAIM_TYPE_EARN, vaultDenom2, types.RewardIndexes{
		{
			CollateralType: "uaeth",
			RewardFactor:   initialVault2RewardFactor.Add(expectedIndexes2),
		},
	})
}
