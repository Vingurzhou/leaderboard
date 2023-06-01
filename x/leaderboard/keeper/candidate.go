package keeper

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"leaderboard/x/leaderboard/types"
)

// TransmitCandidatePacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitCandidatePacket(
	ctx sdk.Context,
	packetData types.CandidatePacketData,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) (uint64, error) {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return 0, sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packetBytes, err := packetData.GetBytes()
	if err != nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "cannot marshal the packet: %w", err)
	}

	return k.channelKeeper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, packetBytes)
}

// OnRecvCandidatePacket processes packet reception
func (k Keeper) OnRecvCandidatePacket(ctx sdk.Context, packet channeltypes.Packet, data types.CandidatePacketData) (packetAck types.CandidatePacketAck, err error) {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return packetAck, err
	}

	// TODO: packet reception logic
	k.SetPlayerInfo(ctx, *data.PlayerInfo)

	// Update the board
	board, found := k.GetBoard(ctx)
	if !found {
		panic("Leaderboard not found")
	}
	listed := board.PlayerInfo
	replaced := false
	for i := range listed {
		if listed[i].Index == data.PlayerInfo.Index {
			listed[i] = *data.PlayerInfo
			replaced = true
			break
		}
	}
	if !replaced {
		listed = append(listed, *data.PlayerInfo)
	}
	k.UpdateBoard(ctx, listed)

	return packetAck, nil
}

// OnAcknowledgementCandidatePacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementCandidatePacket(ctx sdk.Context, packet channeltypes.Packet, data types.CandidatePacketData, ack channeltypes.Acknowledgement) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:

		// TODO: failed acknowledgement logic
		_ = dispatchedAck.Error

		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck types.CandidatePacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return errors.New("cannot unmarshal acknowledgment")
		}

		// TODO: successful acknowledgement logic

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return errors.New("invalid acknowledgment format")
	}
}

// OnTimeoutCandidatePacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutCandidatePacket(ctx sdk.Context, packet channeltypes.Packet, data types.CandidatePacketData) error {

	// TODO: packet timeout logic

	return nil
}
