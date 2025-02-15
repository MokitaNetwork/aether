package types_test

import (
	"encoding/json"
	"testing"

	"github.com/mokitanetwork/aether/x/swap/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/yaml"
)

func TestGenesis_Default(t *testing.T) {
	defaultGenesis := types.DefaultGenesisState()

	require.NoError(t, defaultGenesis.Validate())

	defaultParams := types.DefaultParams()
	assert.Equal(t, defaultParams, defaultGenesis.Params)
}

func TestGenesis_Validate_SwapFee(t *testing.T) {
	type args struct {
		name      string
		swapFee   sdk.Dec
		expectErr bool
	}
	// More comprehensive swap fee tests are in prams_test.go
	testCases := []args{
		{
			"normal",
			sdk.MustNewDecFromStr("0.25"),
			false,
		},
		{
			"negative",
			sdk.MustNewDecFromStr("-0.5"),
			true,
		},
		{
			"greater than 1.0",
			sdk.MustNewDecFromStr("1.001"),
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			genesisState := types.GenesisState{
				Params: types.Params{
					AllowedPools: types.DefaultAllowedPools,
					SwapFee:      tc.swapFee,
				},
			}

			err := genesisState.Validate()
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGenesis_Validate_AllowedPools(t *testing.T) {
	type args struct {
		name      string
		pairs     types.AllowedPools
		expectErr bool
	}
	// More comprehensive pair validation tests are in pair_test.go, params_test.go
	testCases := []args{
		{
			"normal",
			types.DefaultAllowedPools,
			false,
		},
		{
			"invalid",
			types.AllowedPools{
				{
					TokenA: "same",
					TokenB: "same",
				},
			},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			genesisState := types.GenesisState{
				Params: types.Params{
					AllowedPools: tc.pairs,
					SwapFee:      types.DefaultSwapFee,
				},
			}

			err := genesisState.Validate()
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGenesis_JSONEncoding(t *testing.T) {
	raw := `{
    "params": {
			"allowed_pools": [
			  {
			    "token_a": "uaeth",
					"token_b": "usdx"
				},
			  {
			    "token_a": "hard",
					"token_b": "busd"
				}
			],
			"swap_fee": "0.003000000000000000"
		},
		"pool_records": [
		  {
				"pool_id": "uaeth:usdx",
			  "reserves_a": { "denom": "uaeth", "amount": "1000000" },
			  "reserves_b": { "denom": "usdx", "amount": "5000000" },
			  "total_shares": "3000000"
			},
		  {
			  "pool_id": "hard:usdx",
			  "reserves_a": { "denom": "uaeth", "amount": "1000000" },
			  "reserves_b": { "denom": "usdx", "amount": "2000000" },
			  "total_shares": "2000000"
			}
		],
		"share_records": [
		  {
		    "depositor": "aeth1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w",
		    "pool_id": "uaeth:usdx",
		    "shares_owned": "100000"
			},
		  {
		    "depositor": "aeth1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea",
		    "pool_id": "hard:usdx",
		    "shares_owned": "200000"
			}
		]
	}`

	var state types.GenesisState
	err := json.Unmarshal([]byte(raw), &state)
	require.NoError(t, err)

	assert.Equal(t, 2, len(state.Params.AllowedPools))
	assert.Equal(t, sdk.MustNewDecFromStr("0.003"), state.Params.SwapFee)
	assert.Equal(t, 2, len(state.PoolRecords))
	assert.Equal(t, 2, len(state.ShareRecords))
}

func TestGenesis_YAMLEncoding(t *testing.T) {
	expected := `params:
  allowed_pools:
  - token_a: uaeth
    token_b: usdx
  - token_a: hard
    token_b: busd
  swap_fee: "0.003000000000000000"
pool_records:
- pool_id: uaeth:usdx
  reserves_a:
    amount: "1000000"
    denom: uaeth
  reserves_b:
    amount: "5000000"
    denom: usdx
  total_shares: "3000000"
- pool_id: hard:usdx
  reserves_a:
    amount: "1000000"
    denom: hard
  reserves_b:
    amount: "2000000"
    denom: usdx
  total_shares: "1500000"
share_records:
- depositor: aeth1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w
  pool_id: uaeth:usdx
  shares_owned: "100000"
- depositor: aeth1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea
  pool_id: hard:usdx
  shares_owned: "200000"
`

	depositor_1, err := sdk.AccAddressFromBech32("aeth1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w")
	require.NoError(t, err)
	depositor_2, err := sdk.AccAddressFromBech32("aeth1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea")
	require.NoError(t, err)

	state := types.NewGenesisState(
		types.NewParams(
			types.NewAllowedPools(
				types.NewAllowedPool("uaeth", "usdx"),
				types.NewAllowedPool("hard", "busd"),
			),
			sdk.MustNewDecFromStr("0.003"),
		),
		types.PoolRecords{
			types.NewPoolRecord(sdk.NewCoins(uaeth(1e6), usdx(5e6)), i(3e6)),
			types.NewPoolRecord(sdk.NewCoins(hard(1e6), usdx(2e6)), i(15e5)),
		},
		types.ShareRecords{
			types.NewShareRecord(depositor_1, types.PoolID("uaeth", "usdx"), i(1e5)),
			types.NewShareRecord(depositor_2, types.PoolID("hard", "usdx"), i(2e5)),
		},
	)

	data, err := yaml.Marshal(state)
	require.NoError(t, err)

	assert.Equal(t, expected, string(data))
}

func TestGenesis_ValidatePoolRecords(t *testing.T) {
	invalidPoolRecord := types.NewPoolRecord(sdk.NewCoins(uaeth(1e6), usdx(5e6)), i(-1))

	state := types.NewGenesisState(
		types.DefaultParams(),
		types.PoolRecords{invalidPoolRecord},
		types.ShareRecords{},
	)

	assert.Error(t, state.Validate())
}

func TestGenesis_ValidateShareRecords(t *testing.T) {
	depositor, err := sdk.AccAddressFromBech32("aeth1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w")
	require.NoError(t, err)

	invalidShareRecord := types.NewShareRecord(depositor, "", i(-1))

	state := types.NewGenesisState(
		types.DefaultParams(),
		types.PoolRecords{},
		types.ShareRecords{invalidShareRecord},
	)

	assert.Error(t, state.Validate())
}

func TestGenesis_Validate_PoolShareIntegration(t *testing.T) {
	depositor_1, err := sdk.AccAddressFromBech32("aeth1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w")
	require.NoError(t, err)
	depositor_2, err := sdk.AccAddressFromBech32("aeth1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea")
	require.NoError(t, err)

	testCases := []struct {
		name         string
		poolRecords  types.PoolRecords
		shareRecords types.ShareRecords
		expectedErr  string
	}{
		{
			name: "single pool record, zero share records",
			poolRecords: types.PoolRecords{
				types.NewPoolRecord(sdk.NewCoins(uaeth(1e6), usdx(5e6)), i(3e6)),
			},
			shareRecords: types.ShareRecords{},
			expectedErr:  "total depositor shares 0 not equal to pool 'uaeth:usdx' total shares 3000000",
		},
		{
			name:        "zero pool records, one share record",
			poolRecords: types.PoolRecords{},
			shareRecords: types.ShareRecords{
				types.NewShareRecord(depositor_1, types.PoolID("uaeth", "usdx"), i(5e6)),
			},
			expectedErr: "total depositor shares 5000000 not equal to pool 'uaeth:usdx' total shares 0",
		},
		{
			name: "one pool record, one share record",
			poolRecords: types.PoolRecords{
				types.NewPoolRecord(sdk.NewCoins(uaeth(1e6), usdx(5e6)), i(3e6)),
			},
			shareRecords: types.ShareRecords{
				types.NewShareRecord(depositor_1, "uaeth:usdx", i(15e5)),
			},
			expectedErr: "total depositor shares 1500000 not equal to pool 'uaeth:usdx' total shares 3000000",
		},
		{
			name: "more than one pool records, more than one share record",
			poolRecords: types.PoolRecords{
				types.NewPoolRecord(sdk.NewCoins(uaeth(1e6), usdx(5e6)), i(3e6)),
				types.NewPoolRecord(sdk.NewCoins(hard(1e6), usdx(2e6)), i(2e6)),
			},
			shareRecords: types.ShareRecords{
				types.NewShareRecord(depositor_1, types.PoolID("uaeth", "usdx"), i(15e5)),
				types.NewShareRecord(depositor_2, types.PoolID("uaeth", "usdx"), i(15e5)),
				types.NewShareRecord(depositor_1, types.PoolID("hard", "usdx"), i(1e6)),
			},
			expectedErr: "total depositor shares 1000000 not equal to pool 'hard:usdx' total shares 2000000",
		},
		{
			name: "valid case with many pool records and share records",
			poolRecords: types.PoolRecords{
				types.NewPoolRecord(sdk.NewCoins(uaeth(1e6), usdx(5e6)), i(3e6)),
				types.NewPoolRecord(sdk.NewCoins(hard(1e6), usdx(2e6)), i(2e6)),
				types.NewPoolRecord(sdk.NewCoins(hard(7e6), uaeth(10e6)), i(8e6)),
			},
			shareRecords: types.ShareRecords{
				types.NewShareRecord(depositor_1, types.PoolID("uaeth", "usdx"), i(15e5)),
				types.NewShareRecord(depositor_2, types.PoolID("uaeth", "usdx"), i(15e5)),
				types.NewShareRecord(depositor_1, types.PoolID("hard", "usdx"), i(2e6)),
				types.NewShareRecord(depositor_1, types.PoolID("hard", "uaeth"), i(3e6)),
				types.NewShareRecord(depositor_2, types.PoolID("hard", "uaeth"), i(5e6)),
			},
			expectedErr: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			state := types.NewGenesisState(types.DefaultParams(), tc.poolRecords, tc.shareRecords)
			err := state.Validate()

			if tc.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.expectedErr)
			}
		})
	}
}
