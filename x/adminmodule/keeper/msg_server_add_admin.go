package keeper

import (
	"context"

	"fmt"

	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) AddAdmin(goCtx context.Context, msg *types.MsgAddAdmin) (*types.MsgAddAdminResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte(types.AdminKey))

	storeCreator := store.Get(types.ToAdminKey(msg.Creator))
	if storeCreator == nil {
		return nil, fmt.Errorf("requester %s must be admin to add admins", msg.Creator)
	}

	store.Set(types.ToAdminKey(msg.Admin), []byte{})

	return &types.MsgAddAdminResponse{}, nil
}
