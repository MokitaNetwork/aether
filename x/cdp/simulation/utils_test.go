package simulation_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	sdk "github.com/cosmos/cosmos-sdk/types"

// 	"github.com/mokitanetwork/aether/x/cdp/simulation"
// )

// func TestShiftDec(t *testing.T) {
// 	tests := []struct {
// 		value    sdk.Dec
// 		shift    sdk.Int
// 		expected sdk.Dec
// 	}{
// 		{sdk.MustNewDecFromStr("5.5"), sdk.NewInt(1), sdk.MustNewDecFromStr("55")},
// 		{sdk.MustNewDecFromStr("5.5"), sdk.NewInt(-1), sdk.MustNewDecFromStr("0.55")},
// 		{sdk.MustNewDecFromStr("5.5"), sdk.NewInt(2), sdk.MustNewDecFromStr("550")},
// 		{sdk.MustNewDecFromStr("5.5"), sdk.NewInt(-2), sdk.MustNewDecFromStr("0.055")},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.value.String(), func(t *testing.T) {
// 			require.Equal(t, tt.expected, simulation.ShiftDec(tt.value, tt.shift))
// 		})
// 	}
// }
