package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/mokitanetwork/aether/x/pricefeed/types"
)

// GetParams returns the params from the store
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	var p types.Params
	k.paramSubspace.GetParamSet(ctx, &p)
	return p
}

// SetParams sets params on the store
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetMarkets returns the markets from params
func (k Keeper) GetMarkets(ctx sdk.Context) types.Markets {
	return k.GetParams(ctx).Markets
}

// GetOracles returns the oracles in the pricefeed store
func (k Keeper) GetOracles(ctx sdk.Context, marketID string) ([]sdk.AccAddress, error) {
	for _, m := range k.GetMarkets(ctx) {
		if marketID == m.MarketID {
			return m.Oracles, nil
		}
	}
	return nil, sdkerrors.Wrap(types.ErrInvalidMarket, marketID)
}

// GetOracle returns the oracle from the store or an error if not found
func (k Keeper) GetOracle(ctx sdk.Context, marketID string, address sdk.AccAddress) (sdk.AccAddress, error) {
	oracles, err := k.GetOracles(ctx, marketID)
	if err != nil {
		// Error already wrapped
		return nil, err
	}
	for _, addr := range oracles {
		if addr.Equals(address) {
			return addr, nil
		}
	}
	return nil, sdkerrors.Wrap(types.ErrInvalidOracle, address.String())
}

// GetMarket returns the market if it is in the pricefeed system
func (k Keeper) GetMarket(ctx sdk.Context, marketID string) (types.Market, bool) {
	markets := k.GetMarkets(ctx)

	for i := range markets {
		if markets[i].MarketID == marketID {
			return markets[i], true
		}
	}
	return types.Market{}, false
}

// GetAuthorizedAddresses returns a list of addresses that have special authorization within this module, eg the oracles of all markets.
func (k Keeper) GetAuthorizedAddresses(ctx sdk.Context) []sdk.AccAddress {
	var oracles []sdk.AccAddress
	uniqueOracles := map[string]bool{}

	for _, m := range k.GetMarkets(ctx) {
		for _, o := range m.Oracles {
			// de-dup list of oracles
			if _, found := uniqueOracles[o.String()]; !found {
				oracles = append(oracles, o)
			}
			uniqueOracles[o.String()] = true
		}
	}
	return oracles
}
