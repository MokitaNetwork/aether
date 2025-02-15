package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mokitanetwork/aether/x/aethdist/types"
)

// GetParams returns the params from the store
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets params on the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}
