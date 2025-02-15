package types_test

import (
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/mokitanetwork/aether/x/committee/testutil"
	"github.com/mokitanetwork/aether/x/committee/types"
)

func TestBaseCommittee(t *testing.T) {
	addresses := []sdk.AccAddress{
		sdk.AccAddress(crypto.AddressHash([]byte("AetherTest1"))),
		sdk.AccAddress(crypto.AddressHash([]byte("AetherTest2"))),
		sdk.AccAddress(crypto.AddressHash([]byte("AetherTest3"))),
	}

	testCases := []struct {
		name            string
		createCommittee func() (*types.MemberCommittee, error)
		expectPass      bool
	}{
		{
			name: "normal",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This base committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: true,
		},
		{
			name: "description length too long",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					fmt.Sprintln("This base committee has a long description.",
						"This base committee has a long description. This base committee has a long description.",
						"This base committee has a long description. This base committee has a long description.",
						"This base committee has a long description. This base committee has a long description.",
						"This base committee has a long description. This base committee has a long description.",
						"This base committee has a long description. This base committee has a long description.",
						"This base committee has a long description. This base committee has a long description.",
						"This base committee has a long description. This base committee has a long description."),
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: false,
		},
		{
			name: "no members",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This base committee is for testing.",
					[]sdk.AccAddress{},
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: false,
		},
		{
			name: "duplicate member",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This base committee is for testing.",
					[]sdk.AccAddress{addresses[2], addresses[2]},
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: false,
		},
		{
			name: "nil permissions",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This base committee is for testing.",
					addresses[:3],
					[]types.Permission{nil},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: false,
		},
		{
			name: "negative proposal duration",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This base committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*-7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: false,
		},
		{
			name: "vote threshold is nil",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This base committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					sdk.Dec{},
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: false,
		},
		{
			name: "vote threshold is 0",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This base committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: false,
		},
		{
			name: "vote threshold above 1",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This base committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("1.001"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: false,
		},
		{
			name: "invalid tally option",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This base committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_UNSPECIFIED,
				)
			},
			expectPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			committee, err := tc.createCommittee()
			if err != nil {
				require.False(t, tc.expectPass)
			} else {
				err = committee.BaseCommittee.Validate()
				if tc.expectPass {
					require.NoError(t, err)
				} else {
					require.Error(t, err)
				}
			}
		})
	}
}

func TestMemberCommittee(t *testing.T) {
	addresses := []sdk.AccAddress{
		sdk.AccAddress(crypto.AddressHash([]byte("AetherTest1"))),
		sdk.AccAddress(crypto.AddressHash([]byte("AetherTest2"))),
		sdk.AccAddress(crypto.AddressHash([]byte("AetherTest3"))),
	}

	testCases := []struct {
		name            string
		createCommittee func() (*types.MemberCommittee, error)
		expectPass      bool
	}{
		{
			name: "normal",
			createCommittee: func() (*types.MemberCommittee, error) {
				return types.NewMemberCommittee(
					1,
					"This member committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
				)
			},
			expectPass: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			committee, err := tc.createCommittee()
			require.NoError(t, err)
			require.Equal(t, types.MemberCommitteeType, committee.GetType())

			err = committee.Validate()
			if tc.expectPass {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

// TestTokenCommittee tests unique TokenCommittee functionality
func TestTokenCommittee(t *testing.T) {
	addresses := []sdk.AccAddress{
		sdk.AccAddress(crypto.AddressHash([]byte("AetherTest1"))),
		sdk.AccAddress(crypto.AddressHash([]byte("AetherTest2"))),
		sdk.AccAddress(crypto.AddressHash([]byte("AetherTest3"))),
	}

	testCases := []struct {
		name            string
		createCommittee func() (*types.TokenCommittee, error)
		expectPass      bool
	}{
		{
			name: "normal",
			createCommittee: func() (*types.TokenCommittee, error) {
				return types.NewTokenCommittee(
					1,
					"This token committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
					testutil.D("0.4"),
					"hard",
				)
			},
			expectPass: true,
		},
		{
			name: "nil quorum",
			createCommittee: func() (*types.TokenCommittee, error) {
				return types.NewTokenCommittee(
					1,
					"This token committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
					sdk.Dec{},
					"hard",
				)
			},
			expectPass: false,
		},
		{
			name: "negative quorum",
			createCommittee: func() (*types.TokenCommittee, error) {
				return types.NewTokenCommittee(
					1,
					"This token committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
					testutil.D("-0.1"),
					"hard",
				)
			},
			expectPass: false,
		},
		{
			name: "quroum greater than 1",
			createCommittee: func() (*types.TokenCommittee, error) {
				return types.NewTokenCommittee(
					1,
					"This token committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
					testutil.D("1.001"),
					"hard",
				)
			},
			expectPass: false,
		},
		{
			name: "bond denom as tally denom",
			createCommittee: func() (*types.TokenCommittee, error) {
				return types.NewTokenCommittee(
					1,
					"This token committee is for testing.",
					addresses[:3],
					[]types.Permission{&types.GodPermission{}},
					testutil.D("0.667"),
					time.Hour*24*7,
					types.TALLY_OPTION_FIRST_PAST_THE_POST,
					testutil.D("0.4"),
					types.BondDenom,
				)
			},
			expectPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			committee, err := tc.createCommittee()
			assert.NoError(t, err)
			assert.Equal(t, types.TokenCommitteeType, committee.GetType())

			err = committee.Validate()
			if tc.expectPass {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestProposalGetContent(t *testing.T) {
	mockTitle := "A Title"
	mockDescription := "A Description"
	proposal, err := types.NewProposal(
		govtypes.NewTextProposal(mockTitle, mockDescription),
		1, 1, time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC),
	)
	assert.NoError(t, err)
	content := proposal.GetContent()
	assert.NotNil(t, content)
	assert.Equal(t, mockTitle, content.GetTitle())
	assert.Equal(t, mockDescription, content.GetDescription())
}
