package keeper_test

import (
	"errors"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/stretchr/testify/require"
)

var TestProposal = govv1beta1types.NewTextProposal("Test", "description")

type invalidProposalRoute struct{ govv1beta1types.TextProposal }

func (invalidProposalRoute) ProposalRoute() string { return "nonexistingroute" }

func TestGetSetProposal(t *testing.T) {
	_, ctx, keeper := setupMsgServer(t)

	// Init genesis ProposalID
	keeper.SetProposalID(sdk.UnwrapSDKContext(ctx), 1)

	tp := TestProposal
	proposal, err := keeper.SubmitProposal(sdk.UnwrapSDKContext(ctx), tp)
	require.NoError(t, err)
	proposalID := proposal.ProposalId
	keeper.SetProposal(sdk.UnwrapSDKContext(ctx), proposal)

	gotProposal, ok := keeper.GetProposal(sdk.UnwrapSDKContext(ctx), proposalID)
	require.True(t, ok)
	require.True(t, proposal.Equal(gotProposal))
}

func TestSubmitProposal(t *testing.T) {
	_, ctx, keeper := setupMsgServer(t)

	// Init genesis ProposalID
	keeper.SetProposalID(sdk.UnwrapSDKContext(ctx), 1)

	testCases := []struct {
		content     govv1beta1types.Content
		expectedErr error
	}{
		{&govv1beta1types.TextProposal{Title: "title", Description: "description"}, nil},
		// Keeper does not check the validity of title and description, no error
		{&govv1beta1types.TextProposal{Title: "", Description: "description"}, nil},
		{&govv1beta1types.TextProposal{Title: strings.Repeat("1234567890", 100), Description: "description"}, nil},
		{&govv1beta1types.TextProposal{Title: "title", Description: ""}, nil},
		{&govv1beta1types.TextProposal{Title: "title", Description: strings.Repeat("1234567890", 1000)}, nil},
		// error only when invalid route
		{&invalidProposalRoute{}, govtypes.ErrNoProposalHandlerExists},
	}

	for i, tc := range testCases {
		_, err := keeper.SubmitProposal(sdk.UnwrapSDKContext(ctx), tc.content)
		require.True(t, errors.Is(tc.expectedErr, err), "tc #%d; got: %v, expected: %v", i, err, tc.expectedErr)
	}
}
