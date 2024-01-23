package core

import (
	"context"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/types"
)

type Core struct {
	CommandBridge CommandBridge
}

func NewCore(configurator ...Configurator) (Core, error) {
	c := Core{}

	for _, config := range configurator {
		var err error
		c, err = config(c)
		if err != nil {
			return Core{}, err
		}
	}
	return c, nil
}

func (c *Core) SetSwitch(switchID byte, direction bool) {
	c.CommandBridge.Send(types.SetSwitch{SwitchID: switchID, Direction: direction}.ToBridgeMsg())
}

func (c *Core) EventHandler(ctx context.Context) {
	for {
		msg, err := c.CommandBridge.BlockedReceive(ctx)
		if err != nil {
			return
		}
		fmt.Println(msg)
	}
}
