package cmd

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/config"
	"github.com/cosmos/cosmos-sdk/client/debug"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"

	"github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	"github.com/spf13/cobra"
	tmcli "github.com/tendermint/tendermint/libs/cli"
	ethermintclient "github.com/tharsis/ethermint/client"
	"github.com/tharsis/ethermint/crypto/hd"
	ethermintserver "github.com/tharsis/ethermint/server"
	servercfg "github.com/tharsis/ethermint/server/config"

	"github.com/mokitanetwork/aether/app"
	"github.com/mokitanetwork/aether/app/params"
	aethclient "github.com/mokitanetwork/aether/client"
	"github.com/mokitanetwork/aether/migrate"
)

// EnvPrefix is the prefix environment variables must have to configure the app.
const EnvPrefix = "AETH"

// NewRootCmd creates a new root command for the aeth blockchain.
func NewRootCmd() *cobra.Command {
	app.SetSDKConfig().Seal()

	encodingConfig := app.MakeEncodingConfig()

	initClientCtx := client.Context{}.
		WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(types.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastBlock).
		WithHomeDir(app.DefaultNodeHome).
		WithKeyringOptions(hd.EthSecp256k1Option()).
		WithViper(EnvPrefix)

	rootCmd := &cobra.Command{
		Use:   "aeth",
		Short: "Daemon and CLI for the Aether blockchain.",
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())

			initClientCtx, err := client.ReadPersistentCommandFlags(initClientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			initClientCtx, err = config.ReadFromClientConfig(initClientCtx)
			if err != nil {
				return err
			}

			if err = client.SetCmdClientContextHandler(initClientCtx, cmd); err != nil {
				return err
			}

			customAppTemplate, customAppConfig := servercfg.AppConfig("uaeth")

			return server.InterceptConfigsPreRunHandler(cmd, customAppTemplate, customAppConfig)
		},
	}

	addSubCmds(rootCmd, encodingConfig, app.DefaultNodeHome)

	return rootCmd
}

// addSubCmds registers all the sub commands used by aeth.
func addSubCmds(rootCmd *cobra.Command, encodingConfig params.EncodingConfig, defaultNodeHome string) {
	rootCmd.AddCommand(
		ethermintclient.ValidateChainID(
			genutilcli.InitCmd(app.ModuleBasics, defaultNodeHome),
		),
		genutilcli.CollectGenTxsCmd(banktypes.GenesisBalancesIterator{}, defaultNodeHome),
		migrate.MigrateGenesisCmd(),
		migrate.AssertInvariantsCmd(encodingConfig),
		genutilcli.GenTxCmd(app.ModuleBasics, encodingConfig.TxConfig, banktypes.GenesisBalancesIterator{}, defaultNodeHome),
		genutilcli.ValidateGenesisCmd(app.ModuleBasics),
		AddGenesisAccountCmd(defaultNodeHome),
		tmcli.NewCompletionCmd(rootCmd, true), // TODO add other shells, drop tmcli dependency, unhide?
		// testnetCmd(app.ModuleBasics, banktypes.GenesisBalancesIterator{}), // TODO add
		debug.Cmd(),
		config.Cmd(),
	)

	ac := appCreator{
		encodingConfig: encodingConfig,
	}

	// ethermintserver adds additional flags to start the JSON-RPC server for evm support
	ethermintserver.AddCommands(rootCmd, defaultNodeHome, ac.newApp, ac.appExport, ac.addStartCmdFlags)

	// add keybase, auxiliary RPC, query, and tx child commands
	rootCmd.AddCommand(
		StatusCommand(),
		newQueryCmd(),
		newTxCmd(),
		aethclient.KeyCommands(app.DefaultNodeHome),
	)
}
