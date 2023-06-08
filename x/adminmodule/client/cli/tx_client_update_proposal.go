package cli

import (
	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	"github.com/spf13/cobra"
)

func NewCmdSubmitUpdateClientProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-client [subject-client-id] [substitute-client-id] (--title [title]  | --description [description]) [flags]",
		Args:  cobra.ExactArgs(2),
		Short: "Submit an update IBC client proposal",
		Long: "Submit an update IBC client proposal along with an initial deposit.\n" +
			"Please specify a subject client identifier you want to update..\n" +
			"Please specify the substitute client the subject client will be updated to.",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			title, err := cmd.Flags().GetString(FlagTitle)
			if err != nil {
				return err
			}

			description, err := cmd.Flags().GetString(FlagDescription)
			if err != nil {
				return err
			}

			subjectClientID := args[0]
			substituteClientID := args[1]

			content := clienttypes.NewClientUpdateProposal(title, description, subjectClientID, substituteClientID)

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
