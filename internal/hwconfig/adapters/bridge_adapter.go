package adapters

import (
	"github.com/nanderv/traincontrol-prototype/internal/bridge"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain/codes"
	"github.com/nanderv/traincontrol-prototype/internal/hwconfig"
	"log/slog"
)

type MessageAdapter struct {
	core   *hwconfig.HwConfigurator
	sender hwconfig.BridgeSender[domain.Msg]
}

// Receive a message from a layout
func (ma *MessageAdapter) Receive(msg domain.Msg) error {
	if msg.Type != codes.HW {
		return nil
	}
	switch msg.Val[0] {
	case 1:
		slog.Info("NODE REGISTERED", "Mac", [3]byte{msg.Val[1], msg.Val[2], msg.Val[3]})
	}
	return nil
}

// Send a message towards a layout
func (ma *MessageAdapter) Send(msg hwconfig.Msger[domain.Msg]) error {
	return ma.sender.Send(msg.ToBridgeMsg())
}

func NewMessageAdapter(c *hwconfig.HwConfigurator, b *bridge.SerialBridge) *MessageAdapter {
	m := MessageAdapter{core: c, sender: b}
	c.AddCommandBridge(&m)
	b.AddReceiver(&m)
	return &m
}
