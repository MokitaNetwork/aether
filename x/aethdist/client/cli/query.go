package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/mokitanetwork/aether/x/aethdist/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group aethdist queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		queryParamsCmd(),
		queryBalanceCmd(),
	)

	return cmd
}

func queryParamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "get the aethdist module parameters",
		Long:  "Get the current global aethdist module parameters.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(cliCtx)
			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return cliCtx.PrintProto(&res.Params)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func queryBalanceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance",
		Short: "get the aethdist module balance",
		Long:  "Get the current aethdist module account balance.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(cliCtx)
			res, err := queryClient.Balance(context.Background(), &types.QueryBalanceRequest{})
			if err != nil {
				return err
			}
			return cliCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
