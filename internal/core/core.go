package core

import (
	"errors"
	"fmt"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	layout2 "github.com/nanderv/traincontrol-prototype/internal/core/domain/layout"
)

type Core struct {
	layoutBridges        []MessageSender
	notifyChangeChannels []*chan layout2.Layout
	layout               layout2.Layout
}

func (c *Core) AddNewReturnChannel() *chan layout2.Layout {
	ch := make(chan layout2.Layout)
	c.notifyChangeChannels = append(c.notifyChangeChannels, &ch)
	return &ch
}

func (c *Core) AddCommandBridge(cc MessageSender) {
	c.layoutBridges = append(c.layoutBridges, cc)
	return
}

func NewCore(configurator ...Configurator) (*Core, error) {
	c := Core{}
	c.layout.TrackSwitches = make([]layout2.TrackSwitch, 0)
	c.notifyChangeChannels = make([]*chan layout2.Layout, 0)
	for _, config := range configurator {
		var err error
		err = config(&c)
		if err != nil {
			return &Core{}, err
		}
	}
	return &c, nil
}

func (c *Core) SetSwitchAction(switchID byte, direction bool) error {
	var found bool
	for _, sw := range c.layout.TrackSwitches {
		if sw.Number == switchID {
			found = true
		}
	}
	if !found {
		return fmt.Errorf("switch with id %v not found", switchID)
	}

	return c.sendToBridges(NewSetSwitch(switchID, direction).ToBridgeMsg())
}

func (c *Core) sendToBridges(msg domain.Msg) error {
	var errOut error
	for _, b := range c.layoutBridges {
		err := b.Send(msg)
		if err != nil {
			if errOut == nil {
				errOut = err
			} else {
				errOut = errors.Join(errOut, err)
			}
		}
	}
	return errOut
}
func (c *Core) SetSwitchEvent(msg SetSwitchResult) {
	for i, sw := range c.layout.TrackSwitches {
		if sw.Number == msg.SetSwitch.switchID {
			c.layout.TrackSwitches[i].Direction = msg.SetSwitch.direction
		}
	}
	c.notify()
	fmt.Println(msg.String())
}

func (c *Core) notify() {
	fmt.Println("NOTIFY")
	for _, ch := range c.notifyChangeChannels {
		fmt.Println(ch)
		select {
		case *ch <- c.layout:
			fmt.Println(c.layout)
			return
		}
	}
}
