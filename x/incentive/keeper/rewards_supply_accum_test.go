package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/mokitanetwork/aether/x/incentive/types"
)

type AccumulateSupplyRewardsTests struct {
	unitTester
}

func (suite *AccumulateSupplyRewardsTests) storedTimeEquals(denom string, expected time.Time) {
	storedTime, found := suite.keeper.GetPreviousHardSupplyRewardAccrualTime(suite.ctx, denom)
	suite.True(found)
	suite.Equal(expected, storedTime)
}

func (suite *AccumulateSupplyRewardsTests) storedIndexesEqual(denom string, expected types.RewardIndexes) {
	storedIndexes, found := suite.keeper.GetHardSupplyRewardIndexes(suite.ctx, denom)
	suite.Equal(found, expected != nil)

	if found {
		suite.Equal(expected, storedIndexes)
	} else {
		suite.Empty(storedIndexes)
	}
}

func TestAccumulateSupplyRewards(t *testing.T) {
	suite.Run(t, new(AccumulateSupplyRewardsTests))
}

func (suite *AccumulateSupplyRewardsTests) TestStateUpdatedWhenBlockTimeHasIncreased() {
	denom := "bnb"

	hardKeeper := newFakeHardKeeper().addTotalSupply(c(denom, 1e6), d("1"))
	suite.keeper = suite.NewKeeper(&fakeParamSubspace{}, nil, nil, hardKeeper, nil, nil, nil, nil, nil, nil)

	suite.storeGlobalSupplyIndexes(types.MultiRewardIndexes{
		{
			CollateralType: denom,
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "hard",
					RewardFactor:   d("0.02"),
				},
				{
					CollateralType: "uaeth",
					RewardFactor:   d("0.04"),
				},
			},
		},
	})
	previousAccrualTime := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.keeper.SetPreviousHardSupplyRewardAccrualTime(suite.ctx, denom, previousAccrualTime)

	newAccrualTime := previousAccrualTime.Add(1 * time.Hour)
	suite.ctx = suite.ctx.WithBlockTime(newAccrualTime)

	period := types.NewMultiRewardPeriod(
		true,
		denom,
		time.Unix(0, 0), // ensure the test is within start and end times
		distantFuture,
		cs(c("hard", 2000), c("uaeth", 1000)), // same denoms as in global indexes
	)

	suite.keeper.AccumulateHardSupplyRewards(suite.ctx, period)

	// check time and factors

	suite.storedTimeEquals(denom, newAccrualTime)
	suite.storedIndexesEqual(denom, types.RewardIndexes{
		{
			CollateralType: "hard",
			RewardFactor:   d("7.22"),
		},
		{
			CollateralType: "uaeth",
			RewardFactor:   d("3.64"),
		},
	})
}

func (suite *AccumulateSupplyRewardsTests) TestStateUnchangedWhenBlockTimeHasNotIncreased() {
	denom := "bnb"

	hardKeeper := newFakeHardKeeper().addTotalSupply(c(denom, 1e6), d("1"))
	suite.keeper = suite.NewKeeper(&fakeParamSubspace{}, nil, nil, hardKeeper, nil, nil, nil, nil, nil, nil)

	previousIndexes := types.MultiRewardIndexes{
		{
			CollateralType: denom,
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "hard",
					RewardFactor:   d("0.02"),
				},
				{
					CollateralType: "uaeth",
					RewardFactor:   d("0.04"),
				},
			},
		},
	}
	suite.storeGlobalSupplyIndexes(previousIndexes)
	previousAccrualTime := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.keeper.SetPreviousHardSupplyRewardAccrualTime(suite.ctx, denom, previousAccrualTime)

	suite.ctx = suite.ctx.WithBlockTime(previousAccrualTime)

	period := types.NewMultiRewardPeriod(
		true,
		denom,
		time.Unix(0, 0), // ensure the test is within start and end times
		distantFuture,
		cs(c("hard", 2000), c("uaeth", 1000)), // same denoms as in global indexes
	)

	suite.keeper.AccumulateHardSupplyRewards(suite.ctx, period)

	// check time and factors

	suite.storedTimeEquals(denom, previousAccrualTime)
	expected, f := previousIndexes.Get(denom)
	suite.True(f)
	suite.storedIndexesEqual(denom, expected)
}

func (suite *AccumulateSupplyRewardsTests) TestNoAccumulationWhenSourceSharesAreZero() {
	denom := "bnb"

	hardKeeper := newFakeHardKeeper() // zero total supplys
	suite.keeper = suite.NewKeeper(&fakeParamSubspace{}, nil, nil, hardKeeper, nil, nil, nil, nil, nil, nil)

	previousIndexes := types.MultiRewardIndexes{
		{
			CollateralType: denom,
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "hard",
					RewardFactor:   d("0.02"),
				},
				{
					CollateralType: "uaeth",
					RewardFactor:   d("0.04"),
				},
			},
		},
	}
	suite.storeGlobalSupplyIndexes(previousIndexes)
	previousAccrualTime := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.keeper.SetPreviousHardSupplyRewardAccrualTime(suite.ctx, denom, previousAccrualTime)

	firstAccrualTime := previousAccrualTime.Add(7 * time.Second)
	suite.ctx = suite.ctx.WithBlockTime(firstAccrualTime)

	period := types.NewMultiRewardPeriod(
		true,
		denom,
		time.Unix(0, 0), // ensure the test is within start and end times
		distantFuture,
		cs(c("hard", 2000), c("uaeth", 1000)), // same denoms as in global indexes
	)

	suite.keeper.AccumulateHardSupplyRewards(suite.ctx, period)

	// check time and factors

	suite.storedTimeEquals(denom, firstAccrualTime)
	expected, f := previousIndexes.Get(denom)
	suite.True(f)
	suite.storedIndexesEqual(denom, expected)
}

func (suite *AccumulateSupplyRewardsTests) TestStateAddedWhenStateDoesNotExist() {
	denom := "bnb"

	hardKeeper := newFakeHardKeeper().addTotalSupply(c(denom, 1e6), d("1"))
	suite.keeper = suite.NewKeeper(&fakeParamSubspace{}, nil, nil, hardKeeper, nil, nil, nil, nil, nil, nil)

	period := types.NewMultiRewardPeriod(
		true,
		denom,
		time.Unix(0, 0), // ensure the test is within start and end times
		distantFuture,
		cs(c("hard", 2000), c("uaeth", 1000)),
	)

	firstAccrualTime := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(firstAccrualTime)

	suite.keeper.AccumulateHardSupplyRewards(suite.ctx, period)

	// After the first accumulation only the current block time should be stored.
	// The indexes will be empty as no time has passed since the previous block because it didn't exist.
	suite.storedTimeEquals(denom, firstAccrualTime)
	suite.storedIndexesEqual(denom, nil)

	secondAccrualTime := firstAccrualTime.Add(10 * time.Second)
	suite.ctx = suite.ctx.WithBlockTime(secondAccrualTime)

	suite.keeper.AccumulateHardSupplyRewards(suite.ctx, period)

	// After the second accumulation both current block time and indexes should be stored.
	suite.storedTimeEquals(denom, secondAccrualTime)
	suite.storedIndexesEqual(denom, types.RewardIndexes{
		{
			CollateralType: "hard",
			RewardFactor:   d("0.02"),
		},
		{
			CollateralType: "uaeth",
			RewardFactor:   d("0.01"),
		},
	})
}

func (suite *AccumulateSupplyRewardsTests) TestNoPanicWhenStateDoesNotExist() {
	denom := "bnb"

	hardKeeper := newFakeHardKeeper()
	suite.keeper = suite.NewKeeper(&fakeParamSubspace{}, nil, nil, hardKeeper, nil, nil, nil, nil, nil, nil)

	period := types.NewMultiRewardPeriod(
		true,
		denom,
		time.Unix(0, 0), // ensure the test is within start and end times
		distantFuture,
		cs(),
	)

	accrualTime := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.ctx = suite.ctx.WithBlockTime(accrualTime)

	// Accumulate with no source shares and no rewards per second will result in no increment to the indexes.
	// No increment and no previous indexes stored, results in an updated of nil. Setting this in the state panics.
	// Check there is no panic.
	suite.NotPanics(func() {
		suite.keeper.AccumulateHardSupplyRewards(suite.ctx, period)
	})

	suite.storedTimeEquals(denom, accrualTime)
	suite.storedIndexesEqual(denom, nil)
}

func (suite *AccumulateSupplyRewardsTests) TestNoAccumulationWhenBeforeStartTime() {
	denom := "bnb"

	hardKeeper := newFakeHardKeeper().addTotalSupply(c(denom, 1e6), d("1"))
	suite.keeper = suite.NewKeeper(&fakeParamSubspace{}, nil, nil, hardKeeper, nil, nil, nil, nil, nil, nil)

	previousIndexes := types.MultiRewardIndexes{
		{
			CollateralType: denom,
			RewardIndexes: types.RewardIndexes{
				{
					CollateralType: "hard",
					RewardFactor:   d("0.02"),
				},
				{
					CollateralType: "uaeth",
					RewardFactor:   d("0.04"),
				},
			},
		},
	}
	suite.storeGlobalSupplyIndexes(previousIndexes)
	previousAccrualTime := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.keeper.SetPreviousHardSupplyRewardAccrualTime(suite.ctx, denom, previousAccrualTime)

	firstAccrualTime := previousAccrualTime.Add(10 * time.Second)

	period := types.NewMultiRewardPeriod(
		true,
		denom,
		firstAccrualTime.Add(time.Nanosecond), // start time after accrual time
		distantFuture,
		cs(c("hard", 2000), c("uaeth", 1000)),
	)

	suite.ctx = suite.ctx.WithBlockTime(firstAccrualTime)

	suite.keeper.AccumulateHardSupplyRewards(suite.ctx, period)

	// The accrual time should be updated, but the indexes unchanged
	suite.storedTimeEquals(denom, firstAccrualTime)
	expectedIndexes, f := previousIndexes.Get(denom)
	suite.True(f)
	suite.storedIndexesEqual(denom, expectedIndexes)
}

func (suite *AccumulateSupplyRewardsTests) TestPanicWhenCurrentTimeLessThanPrevious() {
	denom := "bnb"

	hardKeeper := newFakeHardKeeper().addTotalSupply(c(denom, 1e6), d("1"))
	suite.keeper = suite.NewKeeper(&fakeParamSubspace{}, nil, nil, hardKeeper, nil, nil, nil, nil, nil, nil)

	previousAccrualTime := time.Date(1998, 1, 1, 0, 0, 0, 0, time.UTC)
	suite.keeper.SetPreviousHardSupplyRewardAccrualTime(suite.ctx, denom, previousAccrualTime)

	firstAccrualTime := time.Time{}

	period := types.NewMultiRewardPeriod(
		true,
		denom,
		time.Time{}, // start time after accrual time
		distantFuture,
		cs(c("hard", 2000), c("uaeth", 1000)),
	)

	suite.ctx = suite.ctx.WithBlockTime(firstAccrualTime)

	suite.Panics(func() {
		suite.keeper.AccumulateHardSupplyRewards(suite.ctx, period)
	})
}
