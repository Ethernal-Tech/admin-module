package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

// NewCmdSubmitUpgradeProposal implements a command handler for submitting an upgrade IBC client proposal transaction.
func NewCmdSubmitIbcClientUpgradeProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ibc-upgrade [name] [height] [path/to/upgraded_client_state.json] (--title [title] | --description [description]) [flags]",
		Args:  cobra.ExactArgs(3),
		Short: "Submit an IBC upgrade proposal",
		Long: "Submit an IBC client breaking upgrade proposal along with an initial deposit.\n" +
			"The client state specified is the upgraded client state representing the upgraded chain\n" +
			`Example Upgraded Client State JSON: 
{
	"@type":"/ibc.lightclients.tendermint.v1.ClientState",
 	"chain_id":"testchain1",
	"unbonding_period":"1814400s",
	"latest_height":{"revision_number":"0","revision_height":"2"},
	"proof_specs":[{"leaf_spec":{"hash":"SHA256","prehash_key":"NO_HASH","prehash_value":"SHA256","length":"VAR_PROTO","prefix":"AA=="},"inner_spec":{"child_order":[0,1],"child_size":33,"min_prefix_length":4,"max_prefix_length":12,"empty_child":null,"hash":"SHA256"},"max_depth":0,"min_depth":0},{"leaf_spec":{"hash":"SHA256","prehash_key":"NO_HASH","prehash_value":"SHA256","length":"VAR_PROTO","prefix":"AA=="},"inner_spec":{"child_order":[0,1],"child_size":32,"min_prefix_length":1,"max_prefix_length":1,"empty_child":null,"hash":"SHA256"},"max_depth":0,"min_depth":0}],
	"upgrade_path":["upgrade","upgradedIBCState"],
}
			`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			cdc := codec.NewProtoCodec(clientCtx.InterfaceRegistry)

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}

			name := args[0]

			height, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return err
			}

			plan := upgradetypes.Plan{
				Name:   name,
				Height: height,
			}

			clientState, err := tryUnmarshallClientState(cdc, args[2])
			if err != nil {
				return err
			}

			if err != nil {
				return err
			}

			content, err := clienttypes.NewUpgradeProposal(title, description, plan, clientState)
			if err != nil {
				return err
			}

			from := clientCtx.GetFromAddress()

			msg, err := types.NewMsgSubmitProposal(content, from)
			if err != nil {
				return err
			}

			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(FlagTitle, "", "title of proposal")
	cmd.Flags().String(FlagDescription, "", "description of proposal")

	return cmd
}

func tryUnmarshallClientState(cdc *codec.ProtoCodec, clientStateOrFileName string) (exported.ClientState, error) {
	var clientState exported.ClientState

	if err := cdc.UnmarshalInterfaceJSON([]byte(clientStateOrFileName), &clientState); err != nil {
		contents, err := os.ReadFile(clientStateOrFileName)
		if err != nil {
			return nil, fmt.Errorf("neither JSON input nor path to .json file for client state were provided: %w", err)
		}

		if err = cdc.UnmarshalInterfaceJSON(contents, &clientState); err != nil {
			return nil, err
		}

	}
	return clientState, nil
}
