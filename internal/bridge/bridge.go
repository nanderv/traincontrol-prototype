package bridge

import (
	"github.com/dsyx/serialport-go"
	"github.com/nanderv/traincontrol-prototype/internal/bridge/domain"
	"log/slog"
	"sync"
)

// The SerialBridge is responsible for translating commands towards things the railway can understand
type SerialBridge struct {
	returners []MessageReceiver
	port      *serialport.SerialPort

	sync.Mutex
}

func NewSerialBridge() *SerialBridge {
	port, err := serialport.Open("/dev/ttyACM0", serialport.DefaultConfig())
	if err != nil {
		slog.Info("couldn't open serial conn", err)
		port, err = serialport.Open("/dev/ttyACM1", serialport.DefaultConfig())
		if err != nil {
			slog.Info("couldn't open serial conn", err)
			port, err = serialport.Open("/dev/ttyACM2", serialport.DefaultConfig())
			if err != nil {
				slog.Error("couldn't open serial conn", err)
				return nil
			}
		}
	}

	return &SerialBridge{port: port}
}

func (f *SerialBridge) AddReceiver(r MessageReceiver) {
	f.returners = append(f.returners, r)
}

func (f *SerialBridge) Send(m domain.Msg) error {
	slog.Info("OUTBOUND", "message", m)
	f.Lock()
	defer f.Unlock()

	encoded := m.Encode()
	_, err := f.port.Write(encoded[:])
	if err != nil {
		return err
	}

	return nil
}

func (f *SerialBridge) Handle() {
	var buffer = make([]byte, 0)
	for {
		buffer = append(buffer, f.readMessageBytes()...)

		for len(buffer) > domain.RawSize {
			buffer = f.handleMessageFromBuffer(buffer)
		}
	}
}

func (f *SerialBridge) readMessageBytes() []byte {
	bytes := make([]byte, 16)
	n, err := f.port.Read(bytes)
	if err != nil {
		slog.Error("could not read", err)
	}

	bytes = bytes[:n]
	return bytes
}

func (f *SerialBridge) handleMessageFromBuffer(byteBuffer []byte) []byte {
	messageBytesCorrect, msg, numBytesRead := getRawMessage(byteBuffer)

	if messageBytesCorrect {
		f.handleReceivedMessage(msg)
	}

	byteBuffer = byteBuffer[numBytesRead:]

	return byteBuffer
}

func getRawMessage(byteBuffer []byte) (bool, domain.RawMsg, int) {
	counter := 0
	var msg = domain.RawMsg{}
	correct := true
	for i, v := range byteBuffer {
		counter = i + 1
		msg[i] = v
		if !domain.ValidChar(v) {
			correct = false
			break
		}
		if i >= domain.RawSize-1 {
			break
		}
	}
	return correct, msg, counter
}

func (f *SerialBridge) handleReceivedMessage(msg domain.RawMsg) {
	mm, err := msg.Decode()
	if err != nil {
		slog.Error("incorrect message", err)
		return
	}
	slog.Info("message received and sent on", "msg", msg)
	for _, r := range f.returners {
		err = r.Receive(mm)

		if err != nil {
			slog.Error("incorrect message", err)
		}
	}
}
