package adapter

import (
	"encoding/json"
	"log/slog"

	"github.com/MisakaTAT/GTerm/backend/enums"
	"github.com/MisakaTAT/GTerm/backend/pkg/terminal"
	"github.com/MisakaTAT/GTerm/backend/types"
	"github.com/gorilla/websocket"
	"go.bug.st/serial"
)

type Serial struct {
	port serial.Port
	ws   *websocket.Conn
}

func NewSerial(ws *websocket.Conn) *Serial {
	return &Serial{
		ws: ws,
	}
}

func (s *Serial) Open(portName string) error {
	slog.Info("Opening serial port: %s", portName)
	mode := &serial.Mode{
		BaudRate: 9600,
		Parity:   serial.NoParity,
		DataBits: 8,
		StopBits: serial.OneStopBit,
	}
	slog.Debug("Serial config: baudRate=%d, parity=%v, dataBits=%d, stopBits=%v",
		mode.BaudRate, mode.Parity, mode.DataBits, mode.StopBits)

	port, err := serial.Open(portName, mode)
	if err != nil {
		slog.Error("Failed to open serial port: %v", err)
		return err
	}
	s.port = port
	slog.Info("Serial port opened successfully: %s", portName)
	return nil
}

func (s *Serial) close() {
	slog.Info("Closing serial connection")
	if s.port != nil {
		_ = s.port.Close()
	}
}

func (s *Serial) Input(quitSignal chan bool) {
	slog.Info("Starting WebSocket input monitoring")
	defer s.setQuit(quitSignal)
	for {
		select {
		case <-quitSignal:
			slog.Debug("Received quit signal, stopping input handler")
			return
		default:
			_, data, err := s.ws.ReadMessage()
			if err != nil {
				slog.Error("Failed to read WebSocket message: %v", err)
				return
			}
			msg := &terminal.Payload{}
			_ = json.Unmarshal(data, &msg)
			if msg.Type == enums.TerminalTypeCMD {
				slog.Debug("Sending command to serial port: %s", msg.Cmd)
				if _, err = s.port.Write([]byte(msg.Cmd)); err != nil {
					slog.Error("Failed to write to serial port: %v", err)
				}
			}
		}
	}
}

func (s *Serial) Output(quitSignal chan bool) {
	slog.Info("Starting serial port output reading")
	defer s.setQuit(quitSignal)
	if s.port == nil {
		slog.Error("Serial port not open")
		s.setQuit(quitSignal)
	}
	buff := make([]byte, 100)
	for {
		select {
		case <-quitSignal:
			slog.Debug("Received quit signal, stopping output handler")
			return
		default:
			for {
				n, err := s.port.Read(buff)
				if err != nil {
					slog.Error("Failed to read data from serial port: %v", err)
					return
				}
				if n == 0 || len(buff[:n]) == 0 {
					continue
				}
				slog.Debug("Read %d bytes of data from serial port", n)
				if err = s.ws.WriteJSON(&types.Message{
					Type:    enums.TerminalTypeData,
					Content: string(buff[:n]),
				}); err != nil {
					slog.Error("Failed to write WebSocket message: %v", err)
				}
			}
		}
	}
}

func (s *Serial) Wait(quitSignal chan bool) {
	defer s.setQuit(quitSignal)
	<-quitSignal
	s.close()
}

func (s *Serial) setQuit(ch chan bool) {
	ch <- true
}
