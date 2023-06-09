package keeper

import (
	"fmt"

	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

// SubmitProposal create new proposal given a content
func (k Keeper) SubmitProposal(ctx sdk.Context, content govv1types.Content) (govv1types.Proposal, error) {
	if !k.rtr.HasRoute(content.ProposalRoute()) {
		return govv1types.Proposal{}, sdkerrors.Wrap(govtypes.ErrNoProposalHandlerExists, content.ProposalRoute())
	}

	cacheCtx, _ := ctx.CacheContext()
	handler := k.rtr.GetRoute(content.ProposalRoute())
	if err := handler(cacheCtx, content); err != nil {
		return govv1types.Proposal{}, sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, err.Error())
	}

	proposalID, err := k.GetProposalID(ctx)
	if err != nil {
		return govv1types.Proposal{}, err
	}

	headerTime := ctx.BlockHeader().Time

	// submitTime and depositEndTime would not be used
	proposal, err := govv1types.NewProposal(content, proposalID, headerTime, headerTime)
	if err != nil {
		return govv1types.Proposal{}, err
	}

	k.SetProposal(ctx, proposal)
	k.InsertActiveProposalQueue(ctx, proposalID)
	k.SetProposalID(ctx, proposalID+1)

	return proposal, nil
}

// GetProposalID gets the highest proposal ID
func (k Keeper) GetProposalID(ctx sdk.Context) (proposalID uint64, err error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ProposalIDKey)
	if bz == nil {
		return 0, sdkerrors.Wrap(types.ErrInvalidGenesis, "initial proposal ID hasn't been set")
	}

	proposalID = types.GetProposalIDFromBytes(bz)
	return proposalID, nil
}

// SetProposalID sets the new proposal ID to the store
func (k Keeper) SetProposalID(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposalIDKey, types.GetProposalIDBytes(proposalID))
}

// SetProposal set a proposal to store
func (k Keeper) SetProposal(ctx sdk.Context, proposal govv1types.Proposal) {
	store := ctx.KVStore(k.storeKey)

	bz := k.MustMarshalProposal(proposal)

	store.Set(types.ProposalKey(proposal.ProposalId), bz)
}

// GetProposal get proposal from store by ProposalID
func (k Keeper) GetProposal(ctx sdk.Context, proposalID uint64) (govv1types.Proposal, bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.ProposalKey(proposalID))
	if bz == nil {
		return govv1types.Proposal{}, false
	}

	var proposal govv1types.Proposal
	k.MustUnmarshalProposal(bz, &proposal)

	return proposal, true
}

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue
func (k Keeper) InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ActiveProposalQueueKey(proposalID), types.GetProposalIDBytes(proposalID))
}

// RemoveFromActiveProposalQueue removes a proposalID from the Active Proposal Queue
func (k Keeper) RemoveFromActiveProposalQueue(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ActiveProposalQueueKey(proposalID))
}

// IterateActiveProposalsQueue iterates over the proposals in the active proposal queue
// and performs a callback function
func (k Keeper) IterateActiveProposalsQueue(ctx sdk.Context, cb func(proposal govv1types.Proposal) (stop bool)) {
	iterator := k.ActiveProposalQueueIterator(ctx)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		proposalID := types.GetProposalIDFromBytes(iterator.Value())
		proposal, found := k.GetProposal(ctx, proposalID)
		if !found {
			panic(fmt.Sprintf("proposal %d does not exist", proposalID))
		}

		if cb(proposal) {
			break
		}
	}
}

// ActiveProposalQueueIterator returns an sdk.Iterator for all the proposals in the Active Queue
func (k Keeper) ActiveProposalQueueIterator(ctx sdk.Context) sdk.Iterator {
	prefixStore := prefix.NewStore(ctx.KVStore(k.storeKey), types.ActiveProposalQueuePrefix)
	return prefixStore.Iterator(nil, nil)
}

func (k Keeper) MarshalProposal(proposal govv1types.Proposal) ([]byte, error) {
	bz, err := k.cdc.Marshal(&proposal)
	if err != nil {
		return nil, err
	}
	return bz, nil
}

func (k Keeper) UnmarshalProposal(bz []byte, proposal *govv1types.Proposal) error {
	err := k.cdc.Unmarshal(bz, proposal)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) MustMarshalProposal(proposal govv1types.Proposal) []byte {
	bz, err := k.MarshalProposal(proposal)
	if err != nil {
		panic(err)
	}
	return bz
}

func (k Keeper) MustUnmarshalProposal(bz []byte, proposal *govv1types.Proposal) {
	err := k.UnmarshalProposal(bz, proposal)
	if err != nil {
		panic(err)
	}
}
