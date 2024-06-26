package hardware

import (
	domain2 "github.com/nanderv/traincontrol-prototype/internal/hardware/domain"
	"time"
)

type TrackService struct {
	layoutBridge         Sender
	notifyChangeChannels []*chan domain2.HardwareState
	Layout               domain2.HardwareState
}

func (svc *TrackService) AddNewReturnChannel() *chan domain2.HardwareState {
	ch := make(chan domain2.HardwareState)
	svc.notifyChangeChannels = append(svc.notifyChangeChannels, &ch)
	return &ch
}

func (svc *TrackService) SetLayoutSender(cc Sender) {
	svc.layoutBridge = cc
	return
}

func NewTrackService(lay domain2.HardwareState) (*TrackService, error) {
	c := TrackService{}
	c.Layout = lay
	go c.notifyEveryOnce()
	return &c, nil
}

func (svc *TrackService) SetSwitchDirection(switchID string, direction bool) error {
	sw, err := svc.Layout.GetSwitch(switchID)

	if err != nil {
		return err
	}

	return svc.layoutBridge.SetSwitchDirection(sw, direction)
}

func (svc *TrackService) UpdateSwitchState(sw *domain2.TrackSwitch, direction bool) error {
	sw.UpdateDirection(direction)
	svc.notify()
	return nil
}

func (svc *TrackService) notify() {
	for _, ch := range svc.notifyChangeChannels {
		*ch <- svc.Layout
	}
}
func (svc *TrackService) notifyEveryOnce() {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			svc.notify()
		}
	}
}
