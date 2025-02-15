package swap_test

import (
	"testing"

	"github.com/mokitanetwork/aether/app"
	"github.com/mokitanetwork/aether/x/swap"
	"github.com/mokitanetwork/aether/x/swap/testutil"
	"github.com/mokitanetwork/aether/x/swap/types"
	"github.com/stretchr/testify/suite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type genesisTestSuite struct {
	testutil.Suite
}

func (suite *genesisTestSuite) Test_InitGenesis_ValidationPanic() {
	invalidState := types.NewGenesisState(
		types.Params{
			SwapFee: sdk.NewDec(-1),
		},
		types.PoolRecords{},
		types.ShareRecords{},
	)

	suite.Panics(func() {
		swap.InitGenesis(suite.Ctx, suite.Keeper, invalidState)
	}, "expected init genesis to panic with invalid state")
}

func (suite *genesisTestSuite) Test_InitAndExportGenesis() {
	depositor_1, err := sdk.AccAddressFromBech32("aeth1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w")
	suite.Require().NoError(err)
	depositor_2, err := sdk.AccAddressFromBech32("aeth1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea")
	suite.Require().NoError(err)

	// slices are sorted by key as stored in the data store, so init and export can be compared with equal
	state := types.NewGenesisState(
		types.Params{
			AllowedPools: types.AllowedPools{types.NewAllowedPool("uaeth", "usdx")},
			SwapFee:      sdk.MustNewDecFromStr("0.00255"),
		},
		types.PoolRecords{
			types.NewPoolRecord(sdk.NewCoins(sdk.NewCoin("hard", sdk.NewInt(1e6)), sdk.NewCoin("usdx", sdk.NewInt(2e6))), sdk.NewInt(1e6)),
			types.NewPoolRecord(sdk.NewCoins(sdk.NewCoin("uaeth", sdk.NewInt(1e6)), sdk.NewCoin("usdx", sdk.NewInt(5e6))), sdk.NewInt(3e6)),
		},
		types.ShareRecords{
			types.NewShareRecord(depositor_2, types.PoolID("hard", "usdx"), sdk.NewInt(1e6)),
			types.NewShareRecord(depositor_1, types.PoolID("uaeth", "usdx"), sdk.NewInt(3e6)),
		},
	)

	swap.InitGenesis(suite.Ctx, suite.Keeper, state)
	suite.Equal(state.Params, suite.Keeper.GetParams(suite.Ctx))

	poolRecord1, _ := suite.Keeper.GetPool(suite.Ctx, types.PoolID("hard", "usdx"))
	suite.Equal(state.PoolRecords[0], poolRecord1)
	poolRecord2, _ := suite.Keeper.GetPool(suite.Ctx, types.PoolID("uaeth", "usdx"))
	suite.Equal(state.PoolRecords[1], poolRecord2)

	shareRecord1, _ := suite.Keeper.GetDepositorShares(suite.Ctx, depositor_2, types.PoolID("hard", "usdx"))
	suite.Equal(state.ShareRecords[0], shareRecord1)
	shareRecord2, _ := suite.Keeper.GetDepositorShares(suite.Ctx, depositor_1, types.PoolID("uaeth", "usdx"))
	suite.Equal(state.ShareRecords[1], shareRecord2)

	exportedState := swap.ExportGenesis(suite.Ctx, suite.Keeper)
	suite.Equal(state, exportedState)
}

func (suite *genesisTestSuite) Test_Marshall() {
	depositor_1, err := sdk.AccAddressFromBech32("aeth1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w")
	suite.Require().NoError(err)
	depositor_2, err := sdk.AccAddressFromBech32("aeth1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea")
	suite.Require().NoError(err)

	// slices are sorted by key as stored in the data store, so init and export can be compared with equal
	state := types.NewGenesisState(
		types.Params{
			AllowedPools: types.AllowedPools{types.NewAllowedPool("uaeth", "usdx")},
			SwapFee:      sdk.MustNewDecFromStr("0.00255"),
		},
		types.PoolRecords{
			types.NewPoolRecord(sdk.NewCoins(sdk.NewCoin("hard", sdk.NewInt(1e6)), sdk.NewCoin("usdx", sdk.NewInt(2e6))), sdk.NewInt(1e6)),
			types.NewPoolRecord(sdk.NewCoins(sdk.NewCoin("uaeth", sdk.NewInt(1e6)), sdk.NewCoin("usdx", sdk.NewInt(5e6))), sdk.NewInt(3e6)),
		},
		types.ShareRecords{
			types.NewShareRecord(depositor_2, types.PoolID("hard", "usdx"), sdk.NewInt(1e6)),
			types.NewShareRecord(depositor_1, types.PoolID("uaeth", "usdx"), sdk.NewInt(3e6)),
		},
	)

	encodingCfg := app.MakeEncodingConfig()
	cdc := encodingCfg.Marshaler

	bz, err := cdc.Marshal(&state)
	suite.Require().NoError(err, "expected genesis state to marshal without error")

	var decodedState types.GenesisState
	err = cdc.Unmarshal(bz, &decodedState)
	suite.Require().NoError(err, "expected genesis state to unmarshal without error")

	suite.Equal(state, decodedState)
}

func (suite *genesisTestSuite) Test_LegacyJSONConversion() {
	depositor_1, err := sdk.AccAddressFromBech32("aeth1mq9qxlhze029lm0frzw2xr6hem8c3k9ts54w0w")
	suite.Require().NoError(err)
	depositor_2, err := sdk.AccAddressFromBech32("aeth1esagqd83rhqdtpy5sxhklaxgn58k2m3s3mnpea")
	suite.Require().NoError(err)

	// slices are sorted by key as stored in the data store, so init and export can be compared with equal
	state := types.NewGenesisState(
		types.Params{
			AllowedPools: types.AllowedPools{types.NewAllowedPool("uaeth", "usdx")},
			SwapFee:      sdk.MustNewDecFromStr("0.00255"),
		},
		types.PoolRecords{
			types.NewPoolRecord(sdk.NewCoins(sdk.NewCoin("hard", sdk.NewInt(1e6)), sdk.NewCoin("usdx", sdk.NewInt(2e6))), sdk.NewInt(1e6)),
			types.NewPoolRecord(sdk.NewCoins(sdk.NewCoin("uaeth", sdk.NewInt(1e6)), sdk.NewCoin("usdx", sdk.NewInt(5e6))), sdk.NewInt(3e6)),
		},
		types.ShareRecords{
			types.NewShareRecord(depositor_2, types.PoolID("hard", "usdx"), sdk.NewInt(1e6)),
			types.NewShareRecord(depositor_1, types.PoolID("uaeth", "usdx"), sdk.NewInt(3e6)),
		},
	)

	encodingCfg := app.MakeEncodingConfig()
	cdc := encodingCfg.Marshaler
	legacyCdc := encodingCfg.Amino

	protoJson, err := cdc.MarshalJSON(&state)
	suite.Require().NoError(err, "expected genesis state to marshal amino json without error")

	aminoJson, err := legacyCdc.MarshalJSON(&state)
	suite.Require().NoError(err, "expected genesis state to marshal amino json without error")

	suite.JSONEq(string(protoJson), string(aminoJson), "expected json outputs to be equal")

	var importedState types.GenesisState
	err = cdc.UnmarshalJSON(aminoJson, &importedState)
	suite.Require().NoError(err, "expected amino json to unmarshall to proto without error")

	suite.Equal(state, importedState, "expected genesis state to be equal")
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(genesisTestSuite))
}
