package client

import (
	"github.com/cosmos/admin-module/x/adminmodule/client/cli"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// Param change proposal handler.
var ParamChangeProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd)

// Software upgrade proposal handler.
var SoftwareUpgradeProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpgradeProposal)

// Cancel software upgrade proposal handler.
var CancelUpgradeProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitCancelUpgradeProposal)

// IBC Client upgrade proposal handler.
var IBCClientUpgradeProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitIbcClientUpgradeProposal)

// IBC Client update proposal handler.
var IBCClientUpdateProposalHandler = govclient.NewProposalHandler(cli.NewCmdSubmitUpdateClientProposal)
