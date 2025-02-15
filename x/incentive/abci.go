package incentive

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mokitanetwork/aether/x/incentive/keeper"
)

// BeginBlocker runs at the start of every block
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	params := k.GetParams(ctx)

	for _, rp := range params.USDXMintingRewardPeriods {
		k.AccumulateUSDXMintingRewards(ctx, rp)
	}
	for _, rp := range params.HardSupplyRewardPeriods {
		k.AccumulateHardSupplyRewards(ctx, rp)
	}
	for _, rp := range params.HardBorrowRewardPeriods {
		k.AccumulateHardBorrowRewards(ctx, rp)
	}
	for _, rp := range params.DelegatorRewardPeriods {
		k.AccumulateDelegatorRewards(ctx, rp)
	}
	for _, rp := range params.SwapRewardPeriods {
		k.AccumulateSwapRewards(ctx, rp)
	}
	for _, rp := range params.SavingsRewardPeriods {
		k.AccumulateSavingsRewards(ctx, rp)
	}
	for _, rp := range params.EarnRewardPeriods {
		if err := k.AccumulateEarnRewards(ctx, rp); err != nil {
			panic(fmt.Sprintf("failed to accumulate earn rewards: %s", err))
		}
	}

	// New generic RewardPeriods
	for _, mrp := range params.RewardPeriods {
		for _, rp := range mrp.RewardPeriods {
			if err := k.AccumulateRewards(ctx, mrp.ClaimType, rp); err != nil {
				panic(fmt.Errorf("failed to accumulate rewards for claim type %s: %w", mrp.ClaimType, err))
			}
		}
	}
}
