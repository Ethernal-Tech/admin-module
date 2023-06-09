package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/admin-module/x/adminmodule/types"
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govv1types "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	// this line is used by starport scaffolding # ibc/keeper/import
)

type (
	Keeper struct {
		cdc                       codec.Codec
		storeKey                  storetypes.StoreKey
		memKey                    storetypes.StoreKey
		rtr                       govv1types.Router
		IsProposalTypeWhitelisted func(govv1types.Content) bool
		// this line is used by starport scaffolding # ibc/keeper/attribute
	}
)

func NewKeeper(
	cdc codec.Codec,
	storeKey,
	memKey storetypes.StoreKey,
	rtr govv1types.Router,
	isProposalTypeWhitelisted func(govv1types.Content) bool,
	// this line is used by starport scaffolding # ibc/keeper/parameter
) *Keeper {
	return &Keeper{
		cdc:                       cdc,
		storeKey:                  storeKey,
		memKey:                    memKey,
		rtr:                       rtr,
		IsProposalTypeWhitelisted: isProposalTypeWhitelisted,
		// this line is used by starport scaffolding # ibc/keeper/return
	}
}

// Router returns the adminmodule Keeper's Router
func (k Keeper) Router() govv1types.Router {
	return k.rtr
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
