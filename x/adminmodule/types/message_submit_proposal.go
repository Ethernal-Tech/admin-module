package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1b1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/gogoproto/proto"
	"gopkg.in/yaml.v2"
)

var (
	_ sdk.Msg                          = &MsgSubmitProposal{}
	_ cdctypes.UnpackInterfacesMessage = &MsgSubmitProposal{}
)

func NewMsgSubmitProposal(content govtypesv1b1.Content, proposer sdk.AccAddress) (*MsgSubmitProposal, error) {
	m := &MsgSubmitProposal{
		Proposer: proposer.String(),
	}

	err := m.SetContent(content)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MsgSubmitProposal) GetContent() govtypesv1b1.Content {
	content, ok := m.Content.GetCachedValue().(govtypesv1b1.Content)
	if !ok {
		return nil
	}
	return content
}

func (m *MsgSubmitProposal) SetContent(content govtypesv1b1.Content) error {
	msg, ok := content.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := types.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Content = any
	return nil
}

func (m *MsgSubmitProposal) Route() string {
	return RouterKey
}

func (m *MsgSubmitProposal) Type() string {
	return "SubmitProposal"
}

func (m *MsgSubmitProposal) GetSigners() []sdk.AccAddress {
	proposer, err := sdk.AccAddressFromBech32(m.Proposer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{proposer}
}

func (m *MsgSubmitProposal) GetSignBytes() []byte {
	bz := govcodec.ModuleCdc.MustMarshalJSON(m)
	return sdk.MustSortJSON(bz)
}

// String implements the Stringer interface
func (m *MsgSubmitProposal) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// ValidateBasic implements Msg
func (m *MsgSubmitProposal) ValidateBasic() error {
	if m.Proposer == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, m.Proposer)
	}

	content := m.GetContent()
	if content == nil {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "missing content")
	}
	if !govtypesv1b1.IsValidProposalType(content.ProposalType()) {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalType, content.ProposalType())
	}
	if err := content.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

func (m MsgSubmitProposal) UnpackInterfaces(unpacker cdctypes.AnyUnpacker) error {
	var content govtypesv1b1.Content
	return unpacker.UnpackAny(m.Content, &content)
}
